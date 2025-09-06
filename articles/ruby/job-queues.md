# ruby / job queues

I use a simple Ruby and Postgres job queuing system:

- Each queue runs 1 job at a time.
- Jobs are worked First In, First Out.
- Jobs are any object with an interface `Job.new(db).call`.
  with optional args `Job.new(db).call(foo: 1, bar: "baz")`.
- The only dependencies are Ruby, Postgres,
  and a custom [DB](/ruby/db) wrapper around the `pg` driver.

## Modest needs

In my app, I have ~20 queues.
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
  callsite text,
  created_at timestamp DEFAULT now() NOT NULL,
  started_at timestamp,
  finished_at timestamp
);
```

Run a Ruby process like:

```bash
bundle exec ruby queues/poll.rb
```

Edit a `queues/poll.rb` file like:

```ruby
require_relative "../lib/db"
require_relative "discord_worker"
require_relative "github_worker"
require_relative "postmark_worker"
require_relative "slack_worker"

$stdout.sync = true

module Queues
  WORKERS = [
    Queues::DiscordWorker,
    Queues::GithubWorker,
    Queues::PostmarkWorker,
    Queues::SlackWorker
  ].freeze
end

# Ensure all workers implement the interface.
Queues::WORKERS.each(&:validate!)

# Ensure queues are only worked on by one worker.
dup_queues = Queues::WORKERS.map(&:queue).tally.select { |_, v| v > 1 }.keys
if dup_queues.any?
  raise "duplicate queues: #{dup_queues.join(", ")}"
end

children = Queues::WORKERS.map do |worker|
  fork do
    worker.new(DB.new).poll
  rescue SignalException
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

Create a base worker class:

```ruby
# queues/poll_worker.rb
module Queues
  class PollWorker
    class << self
      attr_accessor :queue, :jobs
    end

    def self.validate!
      if queue.to_s.strip.empty?
        raise NotImplementedError, "#{name} does not specify a queue"
      end

      if @jobs.nil? || @jobs.empty?
        raise NotImplementedError, "#{name} does not define any jobs"
      end
    end

    attr_reader :db

    def initialize(db)
      @db = db
    end

    def poll
      puts "queue=#{queue} poll=#{poll_interval}s"

      loop do
        sleep poll_interval

        pending_jobs.each do |job|
          result = db.exec(<<~SQL, [job["id"]]).first
            UPDATE jobs
            SET started_at = now(), status = 'started'
            WHERE id = $1
            RETURNING EXTRACT(EPOCH FROM now() - created_at) AS latency
          SQL

          latency = result["latency"].round(2).to_f
          status = work(job_name: job["name"], job_args: job["args"])
        rescue => err
          status = "err: #{err}"
          Sentry.capture_exception(err)
        ensure
          if job && job["id"]
            result = db.exec(<<~SQL, [status, job["id"]]).first
              UPDATE jobs
              SET finished_at = now(), status = $1
              WHERE id = $2
              RETURNING EXTRACT(EPOCH FROM now() - started_at) AS elapsed
            SQL

            elapsed = result["elapsed"].round(2).to_f
            puts %(queue=#{queue} job=#{job["name"]} id=#{job["id"]} status="#{status}" latency=#{latency}s duration=#{elapsed}s)

            min_job_time = 1.0 / max_jobs_per_second
            sleep [min_job_time - elapsed, 0].max
          end
        end
      end
    end

    private def pending_jobs
      db.exec(<<~SQL, [queue])
        SELECT id, name, args
        FROM jobs
        WHERE queue = $1
          AND started_at IS NULL
          AND status = 'pending'
        ORDER BY created_at ASC
      SQL
    end

    private def poll_interval
      10
    end

    private def max_jobs_per_second
      Float::INFINITY
    end

    private def queue
      self.class.queue
    end

    private def work(job_name:, job_args:)
      worker = self.class.jobs.find { |job| job.name == job_name }

      if !worker
        msg = "unknown job `#{job_name}` for queue `#{queue}`"
        Sentry.capture_message(msg, extra: {job_args: job_args})
        return "err: #{msg}"
      end

      if worker.instance_method(:call).arity == 0
        worker.new(db).call
      else
        worker.new(db).call(**job_args.transform_keys(&:to_sym))
      end
    end
  end
end
```

Implement a specific worker:

```ruby
# queues/github_worker.rb
require_relative "poll_worker"
require_relative "../lib/github/job_one"
require_relative "../lib/github/job_two"

module Queues
  class GithubWorker < PollWorker
    @queue = "github"
    @jobs = [Github::JobOne, Github::JobTwo]

    private def max_jobs_per_second
      10 # GitHub API rate limit
    end
  end
end
```

## Enqueuing jobs

Create a helper for inserting jobs:

```ruby
# lib/jobs/insert.rb
module Jobs
  class Insert
    attr_reader :db

    def initialize(db)
      @db = db
    end

    def call(queue:, name:, args: {}, args_params: [])
      loc = caller_locations(1, 1).first
      callsite = "#{loc.path}:#{loc.lineno}"
      param_offset = args_params.size

      case args
      when Hash
        params = args_params + [queue, name, callsite, args]
        sql = "SELECT $#{params.size} AS args"
      when Array
        params = args_params + [queue, name, callsite, args.to_json]
        sql = "SELECT json_array_elements($#{params.size}) AS args"
      when String
        params = args_params + [queue, name, callsite]
        sql = args
      else
        raise ArgumentError, "args must be an array, hash or sql string."
      end

      db.exec(<<~SQL, params)
        WITH data AS (#{sql})
        INSERT INTO jobs (queue, name, callsite, args)
        SELECT
          $#{param_offset + 1},
          $#{param_offset + 2},
          $#{param_offset + 3},
          data.args::jsonb
        FROM data
        ON CONFLICT DO NOTHING
        RETURNING id
      SQL
    end
  end
end
```

Use it to enqueue jobs:

```ruby
require_relative "lib/db"
require_relative "lib/jobs/insert"

i = Jobs::Insert.new(DB.pool)

# Single job
i.call(
  queue: "github",
  name: "JobOne",
  args: {
    company_id: 42
  }
)

# Multiple jobs from SQL
i.call(
  queue: "github",
  name: "JobOne",
  args: <<~SQL
    SELECT jsonb_build_object('company_id', id) AS args
    FROM companies
    WHERE status = 'active'
  SQL
)
```

## Scheduling jobs

A Clock process inserts jobs on a schedule.
See [ruby / clock](/ruby/clock).
