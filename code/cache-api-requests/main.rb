# begindoc: all
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
        client_id: ENV["FSQ_ID"],
        client_secret: ENV["FSQ_SECRET"],
        v: "20180323",
        query: query,
        near: near
      }
    )
    hashed_uri = Digest::MD5.hexdigest(req.uri)

    # lookup cache
    result = DB.exec_params(<<~SQL, [hashed_uri])
      SELECT resp_body
      FROM cache_foursquare
      WHERE hashed_uri = $1
      AND fetched_at > now() - interval '24 hours';
    SQL

    if result.num_tuples == 1
      # return cache if fresh in database
      return [200, result[0]["resp_body"]]
    end

    # GET req.uri
    resp = http.perform(req, HTTP::Options.new({}))

    if resp.code != 200
      # return error
      return [resp.code, JSON.parse(resp.body)]
    end

    # add to cache, or update stale cache
    DB.exec_params(<<~SQL, [hashed_uri, resp.body])
      INSERT INTO cache_foursquare (fetched_at, hashed_uri, resp_body)
      VALUES (now(), $1, $2)
      ON CONFLICT (hashed_uri) DO UPDATE
      SET fetched_at = EXCLUDED.fetched_at, resp_body = EXCLUDED.resp_body;
    SQL

    # return fresh data
    [200, JSON.parse(resp.body)]
  end
end

Foursquare.explore("tacos", near: "San Francisco, CA")
# enddoc: all
