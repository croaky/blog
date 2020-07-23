# Cache API Requests

This article covers a case where
I ingested data from Foursquare's Places API, whose
[terms of use](https://developer.foursquare.com/docs/usage-guidelines/) include:

* the data can not be cached longer than 24 hours
* an hourly rate limit and a daily call quota, whichever comes first

Database schema:

```embed
code/cache-api-requests/schema.sql all
```

I created a database table to store API requests and responses
and wrote a client interface for the API endpoints I needed like:

```ruby
Foursquare.explore("tacos", near: "San Francisco, CA")
```

The first time this code runs, an HTTP request will be made.
The request URL will be saved to a Postgres database.
When it runs within Foursquare's cache policy,
no HTTP request will be made.

```embed
code/cache-api-requests/main.rb all
```

Data older than 30 days needs to additionally be deleted.
It can be swept via a clock process or
[pg_cron](https://github.com/citusdata/pg_cron):

```sql
DELETE FROM cache_foursquare
WHERE fetched_at < now() - interval '24 hours';
```
