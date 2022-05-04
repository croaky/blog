# Postgres-Backed Job Queue in Ruby

A few lines of Ruby with the [pg](https://github.com/ged/ruby-pg) driver
is a simple alternative to a job queuing library.
Job queues are defined as database tables
and workers are defined in one Ruby file.

Depending on the queue requirements,
either polling or Postgres' `LISTEN/NOTIFY` may be appropriate.

```
queue_poll: bundle exec ruby queue_poll.rb
queue_listen: bundle exec ruby queue_listen.rb
```

To run one worker on Heroku:

```
heroku ps:scale queue_poll=1
```

Or:

```
heroku ps:scale queue_listen=1
```

Queues can contain heterogeneous job types.
But, if you want to avoid backup in one queue affecting jobs of another type,
create N numbers of (1 job queue table + 1 job worker).

## Poll

With a `job_queue` table...

```sql
CREATE TABLE job_queue (
  id SERIAL,
  created_at timestamp DEFAULT now() NOT NULL,
  status text DEFAULT 'pending'::text NOT NULL,
  name text NOT NULL,
  data jsonb NOT NULL,
  worked_at timestamp
);
```

...the job worker could look like:

```embed
code/ruby-postgres-job-queue/poll.rb all
```

## Listen

With a `job_queue` table, `NOTIFY` function, and `TRIGGER`...

```sql
CREATE TABLE job_queue (
  id SERIAL,
  created_at timestamp DEFAULT now() NOT NULL,
  status text DEFAULT 'pending'::text NOT NULL,
  name text NOT NULL,
  data jsonb NOT NULL,
  worked_at timestamp
);

CREATE FUNCTION notify_job_queued() RETURNS TRIGGER AS $$
BEGIN
  PERFORM
    pg_notify('job_queued', cast(NEW.id AS varchar));
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER on_job_queue
  AFTER INSERT ON job_queue
  FOR EACH ROW
  EXECUTE PROCEDURE notify_job_queued();
```

...the job worker could look like:

```embed
code/ruby-postgres-job-queue/listen.rb all
```

## Enqueue

In either case, enqueue a job by `INSERT`ing into the queue table.

```ruby
require "pg"
require "json"

conn = PG.connect(ENV.fetch("DATABASE_URL"))

conn.exec_params(<<~SQL, [{company_id: 1}.to_json])
  INSERT INTO job_queue (name, data)
  VALUES ('JobOne', $1)
SQL

job = conn.exec(<<~SQL).first
  SELECT data
  FROM job_queue
SQL

puts JSON.parse(job["data"]).dig("company_id") # 1
```
