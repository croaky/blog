# web / static assets

Content Distribution Networks (CDNs)
pull content from their [origin server] during HTTP requests to cache them:

[origin server]: https://www.rfc-editor.org/rfc/rfc9110.html#name-origin-server

```
DNS -> CDN -> Origin
```

Example:

```
Cloudflare DNS -> Cloudflare CDN -> Render
```

## Without an asset host

If a `CNAME` record for a domain name points to a Rails app on
[Render](https://render.com):

```
www.example.com -> example.onrender.com
```

The first HTTP request for a static asset:

- is received by Render.com's load balancers, which terminates TLS
- is forwarded to the [Render.com web service](https://render.com/docs/web-services),
  which needs to bind to host `0.0.0.0` on a port, usually specified by a `PORT` variable.
- passed to one of the running web server processes.
- routed by the web server to the asset (CSS, JS, img, font file)

The logs will contain lines like this:

```
200 GET /assets/app-ql4h2308y.js
200 GET /assets/app-ql4h2308y.css
```

This isn't the best use of web processes;
they should be reserved for handling application logic.
Response time is degraded by waiting for processes
to finish their work.

## With a CDN as an asset host

In production,
esbuild (more info below) appends a hash of each asset's contents
to the asset's name.
When the file changes,
the browser requests the latest version.

The first time a user requests an asset, it will look like this:

```
200 GET /assets/app-ql4h2308y.css
```

A Cloudflare cache miss "pulls from the origin",
making a `GET` request to the origin,
stores the result in their cache,
and serves the result.

Future `GET` and `HEAD` requests
to the Cloudflare URL within the cache duration
will be cached, with no second HTTP request to the origin.

All HTTP requests using verbs other than `GET` and `HEAD`
proxy through to the origin.

## esbuild config

I recommend using [esbuild](https://esbuild.github.io/) instead of the Rails
asset pipeline.

Example `package.json` configuring React, Sass, and TypeScript:

```json
{
  "name": "app",
  "private": "true",
  "dependencies": {
    "esbuild": "^0.17.15",
    "esbuild-sass-plugin": "^2.8.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-select": "^5.7.2"
  },
  "scripts": {
    "build": "node build.mjs"
  },
  "devDependencies": {
    "typescript": "^5.0.3"
  }
}
```

`build` is run as part of the Render "Build command" during deployment.

Example `build.mjs`:

```js
import * as esbuild from "esbuild";
import { sassPlugin } from "esbuild-sass-plugin";

const args = process.argv.slice(2);

let opts = {
  entryPoints: ["js/app.ts", "css/app.scss"],
  plugins: [sassPlugin()],
  bundle: true,
  external: ["fonts/*"],
  loader: {
    ".woff2": "dataurl",
  },
  tsconfig: "tsconfig.json",
  outdir: "public",
};

// deploy
await esbuild.build({
  ...opts,
  minify: true,
  keepNames: true,
});
```

## Rails config

In `config/environments/production.rb`:

```ruby
config.public_file_server.enabled = true
config.public_file_server.headers = {
  "Cache-Control" => "public, max-age=31536000, immutable"
}
```

Since filenames include content hashes,
each URL is immutable.
When content changes, the filename changes.
Browsers and CDNs can cache aggressively (1 year)
without risk of serving stale content.

The [immutable directive](https://code.facebook.com/posts/557147474482256/this-browser-tweak-saved-60-of-requests-to-facebook/)
eliminates revalidation requests even on page reload.

In `Rakefile`:

```ruby
require "digest"

# ...

namespace :assets do
  task :precompile do
    ["public/css/app.css", "public/js/app.js"].each do |old_path|
      hash = Digest::MD5.file(File.expand_path(old_path, __dir__))
      ext = File.extname(old_path)
      base = old_path.chomp(ext)
      new_path = "#{base}-#{hash}#{ext}"
      system "mv #{old_path} #{new_path}"
    end
  end
end
```

## Render config

In our production web service on
<a href="https://render.com" target="_blank">Render</a>,
our build command will look like this:

```
npm install && npm run build && bundle install && bundle exec rake db:migrate && bundle exec rake assets:precompile
```

In order, this:

1. Builds JavaScript dependencies via NPM
2. Transpiles, bundles, minifies via esbuild
3. Builds Ruby dependencies via Bundler
4. Migrates the database
5. Fingerprints the static assets using the above rake task
