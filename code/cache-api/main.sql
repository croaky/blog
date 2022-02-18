-- begindoc: schema
-- createdb venues
-- psql -d venues -f schema.sql
CREATE TABLE cache_foursquare (
  req_url text NOT NULL,
  resp_body jsonb NOT NULL,
  fetched_at timestamp NOT NULL,
  UNIQUE (req_url)
);

CREATE INDEX req_url_idx ON cache_foursquare (req_url);
-- enddoc: schema

-- begindoc: sweep
DELETE FROM cache_foursquare
WHERE fetched_at < now() - '24 hours'::interval;
-- enddoc: sweep
