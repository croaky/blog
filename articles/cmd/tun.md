# cmd / tun

`tun` tunnels local services to the public internet.

It is [open source](https://github.com/croaky/tun).

## Why

I wanted an ngrok‑like tunnel I could self‑host for a Slack Events API
integration. Fewer moving parts, no vendor accounts, data stays in our
infrastructure.

## Set up the server

Deploy the server to a host like [Render](https://render.com/):

- Public git repository: `https://github.com/croaky/tun`
- Build command: `go build -o tund ./cmd/tund`
- Start command: `./tund`
- Env var: `TUN_TOKEN` (shared secret)
- Health check path: `/health`

Render provides HTTPS automatically.
Render and others (e.g. Railway) offer "scale to zero" plans,
which keeps costs low for occasional tunnels.

Server logs look like:

```
[croaky] tunnel connected
200 POST /slack/events 147.33ms
[croaky] tunnel disconnected
```

## Set up the client

Install:

```sh
go install github.com/croaky/tun/cmd/tun@latest
```

Create a `.env` file in the directory you run `tun`:

```
TUN_SERVER=wss://your-service.onrender.com/tunnel
TUN_LOCAL=http://localhost:3000
TUN_ALLOW="POST /slack/events GET /health"
TUN_TOKEN=your-shared-secret
```

`TUN_ALLOW` accepts space-separated `METHOD /path` pairs (exact match, no wildcards).
Requests not matching a rule return `403 Forbidden`.

`TUN_TOKEN` authenticates the client to the server
via `Authorization: Bearer <token>`.

Run:

```sh
tun
```

Client logs look like:

```
[croaky] connected to wss://your-service.onrender.com/tunnel, forwarding to http://localhost:3000
POST /slack/events
```

The username comes from `git config github.user` (falls back to `$USER`),
helping identify who has the tunnel when teammates share a server.

The client auto-reconnects with exponential backoff (500ms to 30s).
Requests timeout after 30 seconds.
The server accepts one active tunnel at a time;
a new connection closes the previous one.
