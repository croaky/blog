# ruby / clock

I run recurring jobs in a Clock process instead of cron.
The Clock inserts jobs into the [job queue system](/ruby/job-queues),
which handles the actual work.

## Why not cron

Cron syntax is terse to the point of inscrutability.
Debugging failures requires digging through system logs.
Your application code is separated from its schedule.

## Clock process

Run a single Clock process:

```ruby
# schedule/clock.rb
require_relative "../lib/db"
require_relative "../lib/jobs/insert"
require_relative "../lib/calendar"

module Schedule
  JOBS = [
    # Daily at midnight
    {
      queue: "discord",
      name: "Discord::Ingest",
      at?: proc { |t| t.min == 30 && t.hour == 0 }
    },
    # Daily cleanup at 8:30am
    {
      queue: "no_throttle",
      name: "Queues::CleanUp",
      at?: proc { |t| t.min == 30 && t.hour == 8 }
    },
    # Weekly on Tuesdays at 3am
    {
      queue: "github",
      name: "Github::IngestStars",
      at?: proc { |t| t.min == 0 && t.hour == 3 && t.tuesday? }
    },
    # Daily at noon, skip holidays
    {
      queue: "mail",
      name: "Mail::Reminders",
      at?: proc { |t| 
        t.min == 0 && t.hour == 12 && 
        !Calendar.holiday?(t)
      }
    }
  ]

  class Clock
    attr_reader :db

    def initialize(db)
      @db = db
    end

    def tick(seconds:)
      i = Jobs::Insert.new(db)

      loop do
        Schedule::JOBS.each do |job|
          if job.fetch(:at?).call(Time.now.utc)
            puts "insert #{job.fetch(:name)}"
            i.call(
              queue: job.fetch(:queue),
              name: job.fetch(:name)
            )
          end
        end

        sleep(seconds)
      end
    end
  end
end

# Validate job definitions
Schedule::JOBS.each do |job|
  if !job[:at?] || !job[:name] || !job[:queue]
    raise "job missing required keys: #{job.inspect}"
  end
end

if $0 == __FILE__
  $stdout.sync = true
  puts "clock running with #{Schedule::JOBS.size} jobs"
  Schedule::Clock.new(DB.new).tick(seconds: 60)
end
```

## Schedule definitions

Each job is a hash with:

- `queue`: which queue to insert into
- `name`: the job class name
- `at?`: a proc that returns true when the job should run

The proc receives a Time object in UTC.
Use Ruby's time methods for readable schedules:

```ruby
# Sub-hourly
proc { |t| t.min % 15 == 0 }          # every 15 minutes
proc { |t| t.min == 30 }               # once per hour at :30

# Daily
proc { |t| t.min == 0 && t.hour == 3 } # daily at 3am

# Weekly
proc { |t| t.min == 0 && t.hour == 3 && t.tuesday? }

# Monthly
proc { |t| t.min == 0 && t.hour == 3 && t.day == 1 }   # first of month
proc { |t| t.to_date == Date.new(t.year, t.month, -1) } # last of month

# Conditional
proc { |t| t.min == 0 && t.hour == 12 && !Calendar.holiday?(t) }
```

## Integration

This provides:

- A single source of truth for scheduled work
- Job deduplication through database constraints
- Visibility into pending work
- Retry logic from the queue workers
- Rate limiting per API

The Clock runs as a single process.
If it crashes, your process supervisor restarts it.
Jobs that should have run during downtime
will run on the next matching interval.

## Benefits

- Schedule defined in Ruby alongside the job code
- No cron syntax to remember
- Easy to test schedule logic
- Integrates with existing job queue
