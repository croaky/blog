# Cache API Calls

To ingest data from Foursquare's Places API,
[terms of use](https://developer.foursquare.com/docs/usage-guidelines/) include:

* the data can not be cached longer than 24 hours
* an hourly rate limit and a daily call quota, whichever comes first

A database table stores API calls:

```embed
code/cache-api/main.sql schema
```

A Ruby client makes API calls:

```ruby
Foursquare.explore("tacos", near: "San Francisco, CA")
```

The first time this code runs, an HTTP request is made,
the request URL (hashed), response body, and timestamp
are saved to a Postgres database.

When it runs again within Foursquare's cache policy,
no HTTP request is made.

```embed
code/cache-api/main.rb all
```

Old data can be deleted via a clock process or
[pg_cron](https://github.com/citusdata/pg_cron):

```embed
code/cache-api/main.sql sweep
```
