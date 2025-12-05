# cmd / tun

[tun](https://github.com/croaky/tun) tunnels local services
to the public internet.
It is a self-hosted [ngrok](https://ngrok.com/) alternative.

## Install

Install the client:

```sh
go install github.com/croaky/tun/cmd/tun@latest
```

Deploy the server (`tund`) to a host like [Render](https://render.com/):

- Public Git Repository: `https://github.com/croaky/tun`
- Build command: `go build -o tund ./cmd/tund`
- Start command: `./tund`
- Environment variable: `TUN_TOKEN` (shared secret)
- Health Check Path: `/health`

Render provides HTTPS automatically.

## Configure

Create a `.env` file in the directory you run `tun`:

```
TUN_SERVER=wss://your-service.onrender.com/tunnel
TUN_LOCAL=http://localhost:3000
TUN_ALLOW=POST /slack/events GET /health
TUN_TOKEN=your-shared-secret
```

`TUN_ALLOW` accepts space-separated `METHOD /path` pairs (exact match).
Requests not matching a rule return `403 Forbidden`.

`TUN_TOKEN` authenticates the client to the server
via `Authorization: Bearer <token>`.

## Run

```sh
tun
```

Client output:

```
[croaky] connected to wss://your-service.onrender.com/tunnel, forwarding to http://localhost:3000
POST /slack/events
```

Server output:

```
[croaky] tunnel connected
200 POST /slack/events 147.33ms
[croaky] tunnel disconnected
```

The client auto-reconnects with exponential backoff (500ms to 30s).
Requests timeout after 30 seconds.
The server accepts one active tunnel at a time;
a new connection closes the previous one.

## Example: Slack Events API

Configure the Slack app's "Event Subscriptions URL" to
`https://your-service.onrender.com/slack/events`.

Run the local Slack bot on port 3000.
Run `tun` to forward Slack events through the tunnel.
