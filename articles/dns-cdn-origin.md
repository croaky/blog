# DNS to CDN to origin

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

If a `CNAME` record for a domain name points to a Rails app on Render:

```
www.example.com -> example.onrender.com
```

The first HTTP request for a static asset:

- is received by Render.com's load balancers, which terminates TLS
- is forwarded to the [Render.com web service](https://render.com/docs/web-services),
  which needs to bind to host `0.0.0.0` on a port, usually specified by a `PORT` variable.
- passed to one of the running [Puma workers](https://github.com/puma/puma) (web server process)
- routed by Rails to the asset (CSS, JS, img, font file)

The logs will contain lines like this:

```
GET "/assets/application-ql4h2308y.js"
GET "/assets/application-ql4h2308y.css"
```

This isn't the best use of Ruby processes;
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
GET www.example.com/application-ql4h2308y.css
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

## esbuild configuation

I recommend not using the Rails asset pipeline and instead using
[esbuild](https://esbuild.github.io/).

Example `package.json` configuring React and TypeScript with linting and typechecking:

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
    "build": "node build.mjs",
    "buildwatch": "node build.mjs --watch",
    "lint": "eslint js",
    "typecheck": "tsc --noEmit"
  },
  "devDependencies": {
    "@tsconfig/recommended": "^1.0.2",
    "@types/react": "^18.0.28",
    "@types/react-dom": "^18.0.11",
    "@typescript-eslint/eslint-plugin": "^5.57.1",
    "@typescript-eslint/parser": "^5.57.1",
    "eslint": "^8.37.0",
    "typescript": "^5.0.3"
  }
}
```

Note the `build` and `buildwatch` scripts:

```
"build": "node build.mjs",
"buildwatch": "node build.mjs --watch",
```

`buildwatch` is run continuously in development.
`build` is run as part of the Render "Build command" during deployment.

Example `build.mjs`:

```js
import * as esbuild from "esbuild";
import { sassPlugin } from "esbuild-sass-plugin";

const args = process.argv.slice(2);
const watch = args.includes("--watch");

let opts = {
  entryPoints: ["js/application.ts", "css/application.scss"],
  plugins: [sassPlugin()],
  bundle: true,
  external: ["fonts/*"],
  loader: {
    ".woff2": "dataurl",
  },
  tsconfig: "tsconfig.json",
  outdir: "public",
};

if (watch) {
  // dev
  let ctx = await esbuild.context({
    ...opts,
    sourcemap: true,
  });
  await ctx.watch();
  console.log("watching...");
} else {
  // deploy
  await esbuild.build({
    ...opts,
    minify: true,
    keepNames: true,
  });
}
```

## Rails configuration

In `config/environments/production.rb`:

```ruby
config.action_controller.asset_host = ENV.fetch("APPLICATION_HOST")
config.public_file_server.enabled = true
  config.public_file_server.headers = {
    "Cache-Control" => "public, s-maxage=2592000, maxage=86400, immutable"
  }
```

The [`immutable` directive](https://code.facebook.com/posts/557147474482256/this-browser-tweak-saved-60-of-requests-to-facebook/)
eliminates revalidation requests.
