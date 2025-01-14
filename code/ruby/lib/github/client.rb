require "http"
require "json"

module GitHub
  class Client
    def get(path) # => (json: Hash | nil, err: String | nil)
      resp = HTTP.timeout(1).get("https://api.github.com#{path}")
      if resp.code / 100 != 2
        return [nil, resp.status]
      end

      [JSON.parse(resp.body), nil]
    rescue => err
      [nil, err.to_s]
    end
  end
end

if $0 == __FILE__
  pp GitHub::Client.new.get("/orgs/thoughtbot/repos")
end
