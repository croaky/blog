# Heroku Postgres read replica with Rails

As a production Ruby on Rails app on Heroku scales,
you may want to move some database read load to a replica
in order to keep some web processes and their database connections
free for other read and write throughput.

As of Rails 6, the framework offers this via
[built-in facilities](https://guides.rubyonrails.org/active_record_multiple_databases.html).
This article demonstrates one possible configuration of these facilities.

In this configuration, automatic connection switching is not enabled.
Instead, the replica is used at targeted call sites in the Rails app.

For example, in `ThingsController#index`:

```ruby
def index
  ApplicationRecord.read_only do
    @things = Thing.order(created_at: :desc)
  end
end
```

This targeted approach is well-suited for offloading
admin actions or queries that don't need to be consistent up to the second.

Both this targeted approach
and the automatic connection switching approach
are not well-suited for user flows where
a user makes a write action
and needs to immediately see the read result
because of the race condition created by
a Heroku Postgres database under load
replicating while the read query comes in.

This is `app/models/application_record.rb`:

```ruby
class ApplicationRecord < ActiveRecord::Base
  self.abstract_class = true

  FOLLOWER_COLOR = ENV["DATABASE_FOLLOWER_COLOR"].to_s.upcase
  FOLLOWER_URL = ENV["HEROKU_POSTGRESQL_#{FOLLOWER_COLOR}_URL"]

  if FOLLOWER_URL.present?
    connects_to database: {writing: :primary, reading: :follower}
  end

  def self.read_only
    if FOLLOWER_URL.present?
      ActiveRecord::Base.connected_to(role: :reading) do
        # All code in this block will run against the follower database.
        # If a write is attempted, an ActiveRecord::ReadOnlyError will raise.
        yield
      end
    else
      # All code in this block will run against the only configured database.
      # Database roles `primary` and `follower` are not available.
      yield
    end
  end
end
```

Heroku Postgres assigns a "color" identifier to each database
that is guaranteed to be unique within the Heroku app.
[`ActiveRecord::ConnectionHandling#connects_to`](https://api.rubyonrails.org/classes/ActiveRecord/ConnectionHandling.html#method-i-connects_to)
is only invoked if the follower database URL is present.

This is important because ActiveRecord will spawn a database pool for
each database "role" of `writing` and `reading` in this case.

If the configuration instead had a fallback for `follower` to
the `primary` Heroku database at `DATABASE_URL`,
we would double our database connections,
which max out at 500 at the highest Heroku Postgres plans.

The `primary` and `follower` databases are therefore
only defined in `config/database.yml` based on the same logic,
depending on the presence of the follower database URL:

```yaml
<% pool = ENV.fetch("PUMA_THREADS", 5) %>
<% follower_color = ENV["DATABASE_FOLLOWER_COLOR"].to_s.upcase %>
<% follower_url = ENV["HEROKU_POSTGRESQL_#{follower_color}_URL"] %>

default: &default
  adapter: postgresql
  encoding: unicode
  pool: <%= pool %>

development:
  primary:
    <<: *default
    database: multi_db_development
  follower:
    <<: *default
    database: multi_db_development
    replica: true

test:
  <<: *default
  database: multi_db_test

production:
  <% if follower_url.present? %>
  primary:
    <<: *default
  follower:
    <<: *default
    replica: true
    url: <%= follower_url %>
  <% else %>
  <<: *default
  <% end %>
```

For simplicity,
the database pool is set to use `PUMA_THREADS` and
all configurable environment variables are set to 5 by default.

```
DATABASE_FOLLOWER_COLOR=SILVER
PUMA_WORKERS=5
PUMA_THREADS=5
```

See `config/puma.rb`:

```ruby
workers ENV.fetch("PUMA_WORKERS", 5).to_i
threads_count = ENV.fetch("PUMA_THREADS", 5).to_i
threads threads_count, threads_count

preload_app!

rackup DefaultRackup
port ENV.fetch("PORT", 3000)
environment ENV.fetch("RAILS_ENV", "development")

on_worker_boot do
  ActiveRecord::Base.establish_connection
end
```

If your app is considering a read replica, you may also be using
[Heroku's Performance dynos](https://devcenter.heroku.com/articles/optimizing-dyno-usage).
With `PUMA_WORKERS=5`, at 400MB of memory for each Rails `web` process,
you'd be at 80% capacity using Performance-M dynos (2.5GB memory each).

The cheapest plan with 500 connections is a
[`standard-3` Heroku Postgres database](https://devcenter.heroku.com/articles/heroku-postgres-plans),
offering 15GB RAM and 512GB storage.
A calculation like the following can help set an upper limit
on the number of dynos for
[Heroku Autoscaling](https://devcenter.heroku.com/articles/scaling#autoscaling).

```
500 maximum database connections =
(up to 25 connections for Heroku internal use) +
(19 web dynos * 5 Puma workers * 5 Puma threads)
```

Test the number of connections by opening a Postgres prompt
using `heroku pg:psql` and running:

```sql
SELECT
  count(*)
FROM
  pg_stat_activity
WHERE
  pid <> pg_backend_pid()
  AND usename = current_user;
```

In another shell, send traffic to the app with a tool like Apache Bench:

```bash
ab -n 100 -c 8 https://tranquil-brushlands-56319.herokuapp.com/
```

As the traffic increases, you'll see the connections max out.
