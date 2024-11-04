require "dotenv/load"
require "http"
require "json"

module Github
  class Client
    def get(path)
      resp = HTTP
        .basic_auth(
          user: ENV.fetch("GITHUB_CLIENT_ID"),
          pass: ENV.fetch("GITHUB_CLIENT_SECRET")
        )
        .timeout(10)
        .get("https://api.github.com#{path}")
      if resp.code / 100 != 2
        return {"err" => "response code #{resp.code}"}
      end

      JSON.parse(resp.body)
    rescue HTTP::TimeoutError
      {"err" => "10s timeout"}
    end
  end
end

if $0 == __FILE__
  pp Github::Client.new.get("/orgs/thoughtbot/repos")
end
