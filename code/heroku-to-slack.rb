# begindoc: all
require "json"
require "net/http"
require "time"
require "uri"

def lambda_handler(event:, context:)
  dyno = JSON.parse(event["body"])
  name = dyno["data"]["name"] || ""
  state = dyno["data"]["state"] || ""

  ignored_states = ["starting", "down"].include?(state)
  term_one_off = state == "crashed" && ["scheduler", "run"].any? { |p| name.include?(p) }

  if ignored_states || term_one_off || name.include?("release")
    return {
      statusCode: 200,
      headers: { "Content-Type": "application/json" },
      body: "ok"
    }
  end

  uri = URI.parse(ENV.fetch("SLACK_URL"))
  http = Net::HTTP.new(uri.host, uri.port)
  http.use_ssl = true
  req = Net::HTTP::Post.new(uri.request_uri)
  req["Content-Type"] = "application/json"

  req.body = {
    text: [
      Time.parse(dyno["created_at"]).getlocal("-07:00").strftime('%H:%M:%S'),
      name,
      "`#{dyno["data"]["command"]}`",
      state
    ].compact.join(" ")
  }.to_json

  res = http.request(req)

  return {
    statusCode: 200,
    headers: { "Content-Type": "application/json" },
    body: res.body
  }
end
# enddoc: all
