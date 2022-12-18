# Deno

[Deno](https://deno.land/) is a secure runtime for JavaScript and TypeScript.
[v1](https://deno.land/v1) was released May 13, 2020.
Its core team includes Ryan Dahl, the creator of Node.

## Why

Deno is the first runtime that has me interested in running
JavaScript/TypeScript on a server (or in serverless functions).

Deno's main innovations are:

1. Secure by default
2. No package manager

## Install

`deno` is shipped as a single executable. Install:

```bash
curl -fsSL https://deno.land/x/install/install.sh | sh
```

Add to `~/.zshrc`:

```bash
export DENO_INSTALL="$HOME/.deno"
export PATH="$DENO_INSTALL/bin:$PATH"
```

## Run

Run untrusted, third-party programs safely
from the command line by referencing the source URL:

```
% deno run https://deno.land/std/examples/welcome.ts
Download https://deno.land/std/examples/welcome.ts
Compile https://deno.land/std/examples/welcome.ts
Welcome to Deno 🦕
```

By default, the program does not have access to
disk, network, subprocesses, or environment variables.
Like browsers, it runs in a secure sandbox.
You can't open files or sockets.

The user has to opt in to those behaviors with flags:

```
--allow-read=/tmp
--allow-write
--allow-net=google.com
--allow-env
```

## Modules

There is no `package.json` or centralized package management server.
The core team feels that style of package manager has caused more harm than good.

Instead, modules are imported explicitly from a server using URLs:

```ts
import { serve } from "https://deno.land/std/http/server.ts";
```

Deno treats modules and files as the same concept.
This is how browsers think about
[ES modules](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import).
In Node, this is not the case but
in Deno, this is explicit.

[Pika CDN](https://www.pika.dev/cdn) is one of the best module servers for Deno.
It manages NPM packages that are bundled as ES modules
and respects Semantic Versioning:

```ts
import { Component, render } from "https://cdn.pika.dev/preact@^10.0.0";
```

To handle private modules, set up a web server.

See your program's module dependencies:

```
% deno info https://deno.land/std/testing/asserts.ts
https://deno.land/std/testing/asserts.ts
  ├── https://deno.land/std/fmt/colors.ts
  └── https://deno.land/std/testing/diff.ts
```

## Standard Library

The [`https://deno.land/std/`](https://deno.land/std) modules
are the standard lib that the JavaScript community has wanted for a long time.

> [The Deno Standard Library] is a loose port of Go's standard library.
> When in doubt, simply port Go's source code, documentation, and tests.
> There are many times when the nature of JavaScript, TypeScript, or Deno itself
> justifies diverging from Go,
> but if possible we want to leverage the energy that went into building Go.
> We generally welcome direct ports of Go's code.

## Types

The TypeScript compiler is compiled into Deno. The team is
[rewriting type checking in Rust](https://github.com/denoland/deno/issues/5432)
for better performance.
Deno was originally prototyped in Go but is now written in Rust.

There are many TypeScript available on npm through `@types/` and
[DefinitelyTyped](https://definitelytyped.org/).
If the remote module server sends an HTTP header
[`X-TypeScript-Types`](https://dev.to/pika/introducing-pika-cdn-deno-p8b)
and it has the content of a URL,
Deno will use those types for type checking.

Even better, because the code has already been bundled,
Deno accesses those types without the overhead of
[transpilation](https://scotch.io/tutorials/javascript-transpilers-what-they-are-why-we-need-them).

Deno ships type definitions for the runtime.
Print them:

```
% deno types
```

## CLIs

Install command line programs:

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

`deno` is a Swiss Army knife of tools installed as an executable
on your machine, a CI machine, or any machine.
