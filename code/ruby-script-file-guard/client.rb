# frozen_string_literal: true

require "dotenv/load"
require "http"
require "json"

module Service
  class Client
    def get(path)
      resp = HTTP
        .headers(
          accept: "application/json",
          apikey: ENV.fetch("API_KEY")
        )
        .timeout(10)
        .get("https://api.example.com#{path}")
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
  pp Service::Client.new.get("/v1/users/1")
end
