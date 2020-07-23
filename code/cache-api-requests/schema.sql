-- begindoc: all
-- createdb venues
-- psql -d venues -f schema.sql

CREATE TABLE cache_foursquare (
  hashed_uri text NOT NULL,
  resp_body jsonb NOT NULL,
  fetched_at timestamp NOT NULL,
  UNIQUE (hashed_uri)
);

CREATE INDEX hashed_uri_idx ON cache_foursquare (hashed_uri);
-- enddoc: all
