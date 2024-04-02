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
