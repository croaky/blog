# Deno

[Deno](https://deno.land/) is a secure runtime for JavaScript and TypeScript.
[v1](https://deno.land/v1) was released May 13, 2020.
Its core team includes Ryan Dahl, the creator of Node.

## Why

With apologies to Node,
Deno is the first runtime that has me interested in running JS/TS on a server
(or in serverless functions).

Deno's main innovations are:

1. Secure by default
2. No package manager

Number 2 is also security-related: the core team feels
that centralized package managers have caused more harm than good.
We could say "If it's on NPM, I can trust it." but that's wrong.

## How

`deno` is shipped as a single executable. Install on macOS:

```
% brew install deno
```

Run untrusted, third-party programs safely
from the command line by referencing the source URL:

```
% deno run https://deno.land/std/examples/welcome.ts
Download https://deno.land/std/examples/welcome.ts
Compile https://deno.land/std/examples/welcome.ts
Welcome to Deno 🦕
```

By default, the program does not have access to
disk, network, subprocesses, or environment vars.
Like browsers, it runs in a secure sandbox.
You can't open files or sockets.

The user has to opt in to those behaviors with flags:

```
--allow-read=/tmp
--allow-write
--allow-net=google.com
--allow-env
```

There is no `package.json` or a centralized server.
Modules are imported explicitly from a server using URLs:

```ts
import { serve } from "https://deno.land/std/http/server.ts"
```

Deno treats modules and files as the same concept.
This is how browsers think about
[ES modules](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import).
In Node, this is not the case but
in Deno, this is explicit.

Deno can be thought of as "a browser for ES modules".
One server right now that is particularly useful for Deno:
[Pika CDN](https://www.pika.dev/cdn).
Pika CDN only deals with NPM packages that are bundled as ES modules
and respects Semantic Versioning:

```ts
import { Component, render } from "https://cdn.pika.dev/preact@^10.0.0"
```

There's lots of TypeScript types out in the world,
available on npm through `@types/` and
[DefinitelyTyped](https://definitelytyped.org/).
How can developers access those types but not have the overhead
of doing the
[transpilation](https://scotch.io/tutorials/javascript-transpilers-what-they-are-why-we-need-them),
take advantage of the fact that the code has already been bundled?

If the remote server sends an HTTP header
[`X-TypeScript-Types`](https://dev.to/pika/introducing-pika-cdn-deno-p8b),
and it has the content of a URL,
Deno will use those types to type-check the package.

If we have these solutions (including Semantic Versioning),
why do we need a package manager?

Users can import URLs to TypeScript code directly.
Deno ships type definitions for the runtime, which can be seen here:

```
% deno types
```

The TypeScript compiler is compiled into Deno. The team is
[rewriting type checking in Rust](https://github.com/denoland/deno/issues/5432)
for better performance.
Deno was originally prototyped in Go but is now written in Rust.

Analyze its dependency tree (also works on any ES module):

```
% deno info https://bit.ly/deno-bronto
```

Install command line programs like this:

```
% deno install --allow-read https://deno.land/std/examples/catj.ts
Download https://deno.land/std/examples/catj.ts
...
Compile https://deno.land/std/examples/catj.ts
Successfully installed catj
/Users/croaky/.deno/bin/catj
Add /Users/croaky/.deno/bin to PATH
  export PATH="/Users/croaky/.deno/bin:$PATH"
```

`deno` is supposed to be a Swiss Army knife of tooling
without much beyond one executable installed on any machine:
yours, CI, etc.
