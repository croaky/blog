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
