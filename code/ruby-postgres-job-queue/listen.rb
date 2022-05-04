# begindoc: all
require "pg"

require_relative "job_one"
require_relative "job_two"

begin
  conn = PG.connect(ENV.fetch("DATABASE_URL"))

  conn.exec "LISTEN job_queued"
  puts "Waiting on job_queued channel..."

  loop do
    conn.wait_for_notify do |event, pid, job_id|
      job = conn.exec_params(<<~SQL, [job_id]).first
        SELECT
          id,
          name,
          data
        FROM
          job_queue
        WHERE
          id = $1
      SQL
      if !job
        next
      end

      t = Process.clock_gettime(Process::CLOCK_MONOTONIC)

      status = case job["name"]
      when "JobOne"
        JobOne.call(job["data"])
      when "JobTwo"
        JobTwo.call(job["data"])
      else
        "err: invalid job #{job["name"]}"
      end

      conn.exec_params(<<~SQL, [status, job["id"]])
        UPDATE
          job_queue
        SET
          status = $1,
          worked_at = now()
        WHERE
          id = $2
      SQL

      elapsed = (Process.clock_gettime(Process::CLOCK_MONOTONIC) - t).round(2)
      puts "#{elapsed}s job #{job["id"]}: #{status}"
    end
  end
ensure
  conn&.exec "UNLISTEN job_queued"
  conn&.close
end
# enddoc: all
