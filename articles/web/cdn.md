# web / cdn

Content Distribution Networks (CDNs)
cache content at edge locations close to users,
reducing latency and load on the origin server.

## Architecture

CDNs pull content from their [origin server] during HTTP requests:

[origin server]: https://www.rfc-editor.org/rfc/rfc9110.html#name-origin-server

```
DNS -> CDN -> Origin
```

Example:

```
Cloudflare DNS -> Cloudflare CDN -> Render
```

## Without a CDN

If a `CNAME` record points directly to an app server:

```
www.example.com -> app.onrender.com
```

Every HTTP request for a static asset:

- is received by the load balancer, which terminates TLS
- is forwarded to the web service
- is passed to a running web server process
- is routed by the web server to the asset (CSS, JS, image, font)

Logs will show:

```
200 GET /css/app.css
200 GET /js/app.js
```

This wastes web processes that should handle application logic,
not serve static files.
Response times degrade as processes queue up.

## With a CDN

The first time a user requests an asset:

```
200 GET /css/app-a1b2c3d4.css
```

A CDN cache miss "pulls from the origin",
making a `GET` request to the origin server,
storing the result in the CDN cache,
and serving the result to the user.

Future `GET` and `HEAD` requests
to the same URL within the cache duration
are served from the CDN cache
with no request to the origin.

All HTTP requests using verbs other than `GET` and `HEAD`
proxy through to the origin.

## Cache invalidation

To maximize cache efficiency,
set long cache headers (1 year)
and change the asset URL when the content changes.

Asset fingerprinting embeds a hash of the file's contents
in the URL or filename:

```
/assets/app-a1b2c3d4.css
```

When the file changes, the hash changes, creating a new URL.
The CDN cache misses and pulls the new version from origin.
Old URLs remain cached but are no longer requested.

Cache-Control header:

```
Cache-Control: public, max-age=31536000, immutable
```

The `immutable` directive
[eliminates revalidation requests](https://code.facebook.com/posts/557147474482256/this-browser-tweak-saved-60-of-requests-to-facebook/)
even on page reload.

## Implementations

See [ruby/fingerprint](/ruby/fingerprint)
or [go/fingerprint](/go/fingerprint)
for file-based fingerprinting implementations.
