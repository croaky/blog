# begindoc: all
require "pg"
require "json"

require_relative "job_one"
require_relative "job_two"

begin
  conn = PG.connect(ENV.fetch("DATABASE_URL"))

  interval = 10
  puts "Polling every #{interval} seconds..."

  loop do
    sleep interval

    t = Process.clock_gettime(Process::CLOCK_MONOTONIC)

    job = conn.exec(<<~SQL).first
      SELECT
        id,
        name,
        data
      FROM
        job_queue
      WHERE
        worked_at IS NULL
      ORDER BY
        created_at ASC
      LIMIT 1
    SQL
    if !job
      next
    end

    status = case job["name"]
    when "JobOne"
      JobOne.call(JSON.parse(job["data"]))
    when "JobTwo"
      JobTwo.call(JSON.parse(job["data"]))
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
ensure
  conn&.close
end
# enddoc: all
