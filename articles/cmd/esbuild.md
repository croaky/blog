# cmd / esbuild

I use [esbuild](https://esbuild.github.io/)
to bundle, minify, and transpile JavaScript and CSS.

It's fast, has sensible defaults,
and works well without a complex configuration.

## Setup

```bash
npm install --save-dev esbuild esbuild-sass-plugin typescript
```

Example `package.json`:

```json
{
  "name": "app",
  "private": "true",
  "dependencies": {
    "esbuild": "^0.17.15",
    "esbuild-sass-plugin": "^2.8.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "scripts": {
    "build": "node build.mjs"
  },
  "devDependencies": {
    "typescript": "^5.0.3"
  }
}
```

## Build script

Example `build.mjs`:

```js
import * as esbuild from "esbuild";
import { sassPlugin } from "esbuild-sass-plugin";

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

await esbuild.build({
  ...opts,
  minify: true,
  keepNames: true,
});
```

This:

- Bundles TypeScript and Sass entry points
- Minifies output for production
- Outputs to `public/app.js` and `public/app.css`
- Preserves function names for better stack traces

Run with:

```bash
npm run build
```

## Development mode

For development with watch mode:

```js
let ctx = await esbuild.context({
  ...opts,
  sourcemap: true,
});

await ctx.watch();
console.log("Watching...");
```

Or use esbuild's built-in dev server:

```js
await ctx.serve({
  servedir: "public",
  port: 3000,
});
```

## Deployment

During deployment,
build before starting the web server.
Example [Render build command](https://render.com/docs/deploys#build-command):

```bash
npm install && npm run build
```

To bust caches for [/web/cdn](CDNs),
fingerprint the output files.
See [go/fingerprint](/go/fingerprint)
or [ruby/fingerprint](/ruby/fingerprint).
