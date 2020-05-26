# Deno

[Deno](https://deno.land/) is a secure runtime for JavaScript and TypeScript.
v1 is out.

`deno` is shipped as an executable. Install on macOS:

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
This is how browser users think about ES modules.
In Node, this is not the case.
In Deno, this is explicit.

Users can import URLs to TypeScript code directly.
Deno ships type definitions for the runtime, which can be seen here:

```
% deno types
```

The TypeScript compiler is compiled into Deno.
The team is [rewriting it in Rust](https://github.com/denoland/deno/issues/5432)
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
Warning Implicitly using master branch https://deno.land/std/examples/catj.ts
Download https://deno.land/std/flags/mod.ts
Download https://deno.land/std/fmt/colors.ts
Warning Implicitly using master branch https://deno.land/std/fmt/colors.ts
Warning Implicitly using master branch https://deno.land/std/flags/mod.ts
Download https://deno.land/std/testing/asserts.ts
Warning Implicitly using master branch https://deno.land/std/testing/asserts.ts
Download https://deno.land/std/testing/diff.ts
Warning Implicitly using master branch https://deno.land/std/testing/diff.ts
Compile https://deno.land/std/examples/catj.ts
✅ Successfully installed catj
/Users/croaky/.deno/bin/catj
ℹ️  Add /Users/croaky/.deno/bin to PATH
    export PATH="/Users/croaky/.deno/bin:$PATH"
```

`deno` is supposed to be a Swiss Army knife of tooling
without much beyond one executable installed on any machine:
yours, CI, etc.
