# Postgres-backed job queues with Ruby

The following describes a simple Ruby and Postgres job queuing system
with these attributes

- Each queue runs 1 job at a time.
- Jobs are worked First In, First Out.
- Jobs are any object with an interface `Job.new(db).call`.
  with optional args `Job.new(db).call(foo: 1, bar: "baz")`.
- The only dependencies are Ruby, Postgres, and
  [the pg gem](https://github.com/ged/ruby-pg).

I have been running a system like this in production for a few years.

## Short history of Ruby job queuing systems

There have been many open source job queuing systems in the Ruby community.

[Sidekiq](https://sidekiq.org/) and [Resque](https://github.com/resque/resque)
are popular but store their jobs in Redis, an additional database to operate.
Other systems have been backed by a relational database
such as [QueueClassic](https://github.com/QueueClassic/queue_classic),
[Delayed Job](https://github.com/collectiveidea/delayed_job) (from Shopify),
[Que](https://github.com/que-rb/que), and
[GoodJob](https://github.com/bensheldon/good_job).

The relational database architecture may be gaining in popularity.
`FOR UPDATE SKIP LOCKED`, which avoids blocking and waiting on locks when polling jobs,
was added to Postgres in 2016 and to MySQL in 2018.
In 2023, Basecamp released
[Solid Queue](https://github.com/basecamp/solid_queue),
an abstraction built around `FOR UPDATE SKIP LOCKED`
that will be [the default ActiveJob backend in Rails 8](https://github.com/rails/rails/issues/50442).

While these are all great projects, I haven't been using them
as I have more modest needs.

## Modest needs

In my application, I have ~20 queues.
~80% of these invoke third-party APIs that have rate limits
such as GitHub, Discord, Slack, and Postmark.
I don't need these jobs to be high-throughput or highly parallel;
processing one at a time is fine.

## How

Create a `jobs` table in Postgres:

```sql
CREATE TABLE jobs (
  id SERIAL,
  queue text NOT NULL,
  name text NOT NULL,
  args jsonb DEFAULT '{}' NOT NULL,
  status text DEFAULT 'pending'::text NOT NULL,
  created_at timestamp DEFAULT now() NOT NULL,
  started_at timestamp,
  finished_at timestamp
);
```

Run a Ruby process like:

```bash
bundle exec ruby queues.rb
```

Edit a `queues.rb` file like:

```embed
code/ruby/queues.rb
```

Edit a `lib/github/worker.rb` file like:

```embed
code/ruby/lib/github/worker.rb
```

Enqueue a job by `INSERT`ing into the jobs table:

```ruby
require "json"
require "pg"

db = PG.connect(ENV.fetch("DATABASE_URL"))

db.exec_params(<<~SQL, [{company_id: 1}.to_json])
  INSERT INTO jobs (queue, name, args)
  VALUES ('github', 'JobOne', $1)
SQL

job = conn.exec(<<~SQL).first
  SELECT
    args
  FROM
    jobs
  ORDER BY
    created_at DESC
  LIMIT 1
SQL

puts JSON.parse(job["args"]).dig("company_id") # 1
```
