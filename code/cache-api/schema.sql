CREATE TABLE cache_foursquare (
  req_url text NOT NULL,
  resp_body jsonb NOT NULL,
  fetched_at timestamp NOT NULL,
  UNIQUE (req_url)
);

CREATE INDEX req_url_idx ON cache_foursquare (req_url);
