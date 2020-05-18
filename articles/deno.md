# Deno

[Deno](https://deno.land/) is a secure runtime for JavaScript and TypeScript.

```
brew install deno
```

You can run untrusted, third-party programs safely
from the command line by referencing the source URL:

```
deno run https://deno.land/std/examples/welcome.ts
Download https://deno.land/std/examples/welcome.ts
Compile https://deno.land/std/examples/welcome.ts
Welcome to Deno 🦕
```

By default, the program does not have access to
disk, network, subprocesses, or environment variables.
You can't open files or sockets.

The user has to opt in to those behaviors with flags:

```
--allow-read=/tmp
--allow-write
--allow-net=google.com
--allow-env
```

Modules are imported using URLs:

```
import { serve } from "https://deno.land/std/http/server.ts"
```

There is no `package.json`.

Be explicit about which server which you get a module from.
This is not dependent on a centralized server.

The TypeScript compiler is compiled into Deno.
V8 Snapshots starts the TS compiler quickly.
Users can import URLs to TypeScript code directly.
Deno ships type definitions for the runtime.

Deno treats modules and files as the same concept.
This is how browser users think about ES modules.
In Node, this is not the case.
In Deno, this is explicit.

Deno was originally prototyped in Go but now is solidly Rust.
