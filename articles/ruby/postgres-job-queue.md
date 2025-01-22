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

```ruby
require "pg"
require_relative "lib/discord/worker"
require_relative "lib/github/worker"
require_relative "lib/postmark/worker"
require_relative "lib/slack/worker"

$stdout.sync = true

workers = [
  Discord::Worker,
  Github::Worker,
  Postmark::Worker,
  Slack::Worker
].freeze

# Ensure all workers implement the interface.
workers.each(&:validate!)

# Ensure queues are only worked on by one worker.
dup_queues = workers.map(&:queue).tally.select { |_, v| v > 1 }.keys
if dup_queues.any?
  raise "duplicate queues: #{dup_queues.join(", ")}"
end

children = workers.map do |worker|
  # Fork a thread for each worker.
  fork do
    # Initialize worker with its own db connection.
    db = PG.connect(ENV.fetch("DATABASE_URL"))
    worker.new(db).poll
  rescue SignalException
    # Prevent child processes from being interrupted.
    # Leave signal handling to the parent process.
  end
end

begin
  children.each { |pid| Process.wait(pid) }
rescue SignalException => sig
  if Signal.list.values_at("HUP", "INT", "KILL", "QUIT", "TERM").include?(sig.signo)
    children.each { |pid| Process.kill("KILL", pid) }
  end
end
```

Edit a `lib/github/worker.rb` file like:

```ruby
require "json"
require "pg"
require_relative "job_one"
require_relative "job_two"

module Github
  class Worker
    attr_reader :db, :queue, :jobs, :poll_interval, :max_jobs_per_second

    def initialize(db)
      @db = db
      @queue = queue
      @jobs = [JobOne, JobTwo]
      @poll_interval = 10

      # https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/rate-limits-for-github-apps
      @max_jobs_per_second = 10
    end

    def poll
      puts "queue=#{queue} poll=#{poll_interval}s"

      loop do
        sleep poll_interval

        pending_jobs.each do |job|
          db.exec_params(<<~SQL, [job["id"]])
            UPDATE
              jobs
            SET
              started_at = now(),
              status = 'started'
            WHERE
              id = $1
          SQL

          worker = jobs.find { |job| job.name == job["name"] }
          status =
            if !worker
              "err: Unknown job `#{name}` for queue `#{queue}`"
            elsif worker.instance_method(:call).arity == 0
              worker.new(db).call
            else
              worker.new(db).call(**job["args"].transform_keys(&:to_sym))
            end
        rescue => err
          status = "err: #{err}"
        ensure
          if job && job["id"]
            elapsed = db.exec_params(<<~SQL, [status, job["id"]]).first["elapsed"]
              UPDATE
                jobs
              SET
                finished_at = now(),
                status = 'ok'
              WHERE
                id = 1
              RETURNING
                round(extract(EPOCH FROM (finished_at - started_at)), 2) AS elapsed
            SQL

            puts %(queue=#{queue} job=#{job["name"]} id=#{job["id"]} status="#{status}" duration=#{elapsed}s)

            min_job_time = 1.0 / max_jobs_per_second
            sleep [min_job_time - elapsed, 0].max
          end
        end
      end
    end

    private def pending_jobs
      db.exec_params(<<~SQL, [queue])
        SELECT
          id,
          name,
          args
        FROM
          jobs
        WHERE
          queue = $1
          AND started_at IS NULL
          AND status = 'pending'
        ORDER BY
          created_at ASC
      SQL
    end
  end
end
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
