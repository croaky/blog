DELETE FROM cache_foursquare
WHERE fetched_at < now() - '24 hours'::interval;
