# Twitter Clone

This tutorial builds a Twitter clone with
React, Next.js, TypeScript, Deno, and Postgres.
The app is deployed to Vercel and Heroku.

## Monorepo

Top-level directory structure:

```
twitter-clone
├── api
└── web
```

## API

Set up:

```
% brew install deno
% deno --version
deno 1.0.5
v8 8.4.300
typescript 3.9.2
```

Edit `api/main.ts`:

```ts
import { serve } from "https://deno.land/std/http/server.ts";

const server = serve({ port: 8000 });
console.log("http://localhost:8000/");

for await (const req of server) {
  req.respond({ body: "Hello World\n" });
}
```

Run:

```
% deno run --allow-net --allow-read api/main.ts
Listening at http://localhost:8000/
```

## Web

Set up:

```
% npx create-next-app
✔ What is your project named? … web
✔ Pick a template › Default starter app
```

Run:

```
% (cd web && yarn dev)
ready - started server on http://localhost:3000
```
