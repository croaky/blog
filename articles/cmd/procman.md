# cmd / procman

`procman` is a process manager for local development on macOS.

It is [open source](https://github.com/croaky/procman).

Install:

```bash
go install github.com/croaky/procman@latest
```

Define your processes in `Procfile.dev`. For example:

```txt
clock: bundle exec ruby schedule/clock.rb
esbuild: npm run buildwatch
queues: bundle exec ruby queues/poll.rb
web: bundle exec puma -p 3000 -C ./config/puma.rb
```

Run processes by naming them in a comma-delimited list:

```bash
procman esbuild,web
```

You must list at least one process.
There are no options or flags.

The output from multiple processes will be combined. For example:

```txt
esbuild |
esbuild | > buildwatch
esbuild | > node build.mjs --watch
esbuild |
esbuild | watching...
web     | [67296] Puma starting in cluster mode...
web     | [67296] * Puma version: 6.4.2 (ruby 3.3.0-p0) ("The Eagle of Durango")
web     | [67296] *  Min threads: 12
web     | [67296] *  Max threads: 12
web     | [67296] *  Environment: development
web     | [67296] *   Master PID: 67296
web     | [67296] *      Workers: 2
web     | [67296] *     Restarts: (✔) hot (✖) phased
web     | [67296] * Preloading application
web     | [67296] * Listening on http://0.0.0.0:3000
web     | [67296] Use Ctrl-C to stop
web     | [67296] - Worker 1 (PID: 67331) booted in 0.9s, phase: 0
web     | [67296] - Worker 0 (PID: 67330) booted in 0.9s, phase: 0
```

Processes until `procman` receives a SIGINT (`Ctrl+C`), `SIGTERM`, or `SIGHUP`.

If one of the processes finishes, it will send a `SIGINT` to all remaining
running processes wait 5s, and then send a `SIGKILL` all remaining processes.
