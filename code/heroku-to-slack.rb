# begindoc: all
require "json"
require "net/http"
require "time"
require "uri"

def lambda_handler(event:, context:)
  json = JSON.parse(event["body"])
  puts json

  if !["dyno", "collaborator"].include?(json["resource"])
    return {
      statusCode: 400,
      headers: {"Content-Type": "application/json"},
      body: "webhook event type not supported"
    }
  end

  localtime = Time
    .parse(json.dig("created_at"))
    .getlocal("-08:00")
    .strftime("%H:%M:%S")

  if json.dig("resource") == "dyno"
    # https://devcenter.heroku.com/articles/webhook-events#api-dyno
    name = json.dig("data", "name") || ""
    state = json.dig("data", "state") || ""

    ignored_states = ["starting", "down"].include?(state)
    term_one_off = state == "crashed" && ["scheduler", "run"].any? { |p| name.include?(p) }

    if ignored_states || term_one_off || name.include?("release")
      return {
        statusCode: 200,
        headers: {"Content-Type": "application/json"},
        body: "ok"
      }
    end

    text = [
      localtime,
      name,
      "`#{json.dig("data", "command")}`", # backticks for code formatting in Slack
      state
    ].compact.join(" ")
  end

  if json.dig("resource") == "collaborator"
    # https://devcenter.heroku.com/articles/webhook-events#api-collaborator

    text = [
      localtime,
      json.dig("actor", "email"),
      "#{json.dig("action")}'d collaborator",
      json.dig("data", "user", "email")
    ].compact.join(" ")
  end

  uri = URI.parse(ENV.fetch("SLACK_URL"))
  http = Net::HTTP.new(uri.host, uri.port)
  http.use_ssl = true

  req = Net::HTTP::Post.new(uri.request_uri)
  req["Content-Type"] = "application/json"
  req.body = {text: text}.to_json

  res = http.request(req)

  {
    statusCode: 200,
    headers: {"Content-Type": "application/json"},
    body: res.body
  }
end
# enddoc: all
