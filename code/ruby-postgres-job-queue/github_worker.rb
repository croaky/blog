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
