require "bundler/inline"
require "digest"
require "json"

gemfile do
  source "https://rubygems.org"
  gem "dotenv"
  gem "http"
  gem "pg"
end

Dotenv.load
DB = PG.connect("postgres://postgres:postgres@localhost:5432/venues")

class Foursquare
  def self.explore(query, near:)
    http = HTTP::Client.new
    req = http.build_request(
      :get,
      "https://api.foursquare.com/v2/venues/explore",
      params: {
        client_id: ENV.fetch("FSQ_ID"),
        client_secret: ENV.fetch("FSQ_SECRET"),
        v: "20180323",
        query: query,
        near: near
      }
    )
    req_url = Digest::MD5.hexdigest(req.uri)

    # lookup cache
    cache = DB.exec_params(<<~SQL, [req_url]).first
      SELECT
        resp_body
      FROM
        cache_foursquare
      WHERE
        req_url = $1
        AND fetched_at > now() - '24 hours'::interval
    SQL
    if cache
      # return cache if fresh
      return [200, cache["resp_body"]]
    end

    # GET req.uri
    resp = http.perform(req, HTTP::Options.new({}))
    if resp.code != 200
      return [resp.code, JSON.parse(resp.body)]
    end

    # add to cache, or update stale cache
    DB.exec_params(<<~SQL, [req_url, resp.body])
      INSERT INTO cache_foursquare (fetched_at, req_url, resp_body)
        VALUES (now(), $1, $2)
      ON CONFLICT (req_url)
        DO UPDATE SET
          fetched_at = EXCLUDED.fetched_at, resp_body = EXCLUDED.resp_body
    SQL

    # return fresh data
    [200, JSON.parse(resp.body)]
  end
end

if $0 == __FILE__
  Foursquare.explore("tacos", near: "San Francisco, CA")
end
