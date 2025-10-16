# cmd / hostctl

I manage `/etc/hosts` for local development and DNS-level ad blocking
with a CLI tool:

```sh
hostctl
```

## Local apps

Browsers treat `http://localhost` and `http://*.localhost` as
[potentially trustworthy origins](https://developer.mozilla.org/en-US/docs/Web/Security/Secure_Contexts),
providing Secure Context features without TLS.

Running multiple apps on different subdomains and ports
(e.g., `blog.localhost:3000`, `htmz.localhost:3001`)
lets me test
[cross-origin security](https://www.alexedwards.net/blog/preventing-csrf-in-go)
without TLS:

- `SameSite=Strict` cookies
- `Sec-Fetch-Site: same-origin` and `Origin` headers
- `Referrer-Policy: strict-origin-when-cross-origin` header

## Ad blocking

Unlike browser extension ad blockers,
DNS-level blocking works on all apps on my device (not only web browsers).

Unlike DNS sinkholes like [Pi-hole](https://pi-hole.net/),
it only works on my laptop (not phones or tablets on the network)
but it does not require an additional always-on device such as a Raspberry Pi
and it works reliably when using the laptop away from home.

Ad/tracker blocking is enabled by default:

```sh
hostctl
```

To disable blocking while keeping local app entries intact:

```sh
hostctl --unblock
```

Script:

```sh
#!/bin/bash
#
# https://github.com/croaky/laptop/blob/main/bin/hostctl

set -euo pipefail

custom_hosts_entries() {
  cat <<EOF
# macOS defaults
255.255.255.255 broadcasthost

# local apps
127.0.0.1 blog.localhost
127.0.0.1 bsfeeds.localhost
127.0.0.1 eds.localhost
127.0.0.1 htmz.localhost
127.0.0.1 neogit.localhost
EOF
}

if [[ "${1:-}" == "--unblock" ]]; then
  cat <<EOF | sudo tee /etc/hosts >/dev/null
# IPv4
127.0.0.1 localhost

# IPv6
::1 localhost

$(custom_hosts_entries)
EOF
else
  # Block ads, trackers, and malicious websites at the DNS host level.
  curl -s https://winhelp2002.mvps.org/hosts.txt | tr -d '\r' | sudo tee /etc/hosts >/dev/null

  # Append custom entries
  (echo && custom_hosts_entries) | sudo tee -a /etc/hosts >/dev/null
fi

# Flush DNS cache
sudo killall -HUP mDNSResponder
```
