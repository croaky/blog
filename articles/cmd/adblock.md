# cmd / adblock

To improve speed, privacy, and safety on my laptop,
I block ads, trackers, and malicious websites at the DNS host level:

```bash
adblock
```

Unlike browser extension ad blockers,
it works on all apps on my device (not only web browsers).

Unlike DNS sinkholes,
it only works on my laptop (not phones or tablets on the network)
but it does not require an additional always-on device such as a Raspberry Pi
and it works reliably when using the laptop away from home.

To disable and re-enable it:

```bash
adblock undo
adblock
```

Script:

```bash
#!/bin/bash
#
# https://github.com/croaky/laptop/blob/main/bin/adblock

set -euo pipefail

if [[ "$1" == "undo" ]]; then
  echo -e '127.0.0.1\tlocalhost\n# MacOS default\n255.255.255.255\tbroadcasthost' | sudo tee /etc/hosts > /dev/null
else
  # Creative Commons Attribution-NonCommercial-ShareAlike
  curl -s https://winhelp2002.mvps.org/hosts.txt > /tmp/adblock

  # Re-write Windows to Unix line endings
  tr -d '\r' < /tmp/adblock > /tmp/etchosts

  # Restore macOS system defaults
  echo -e '# MacOS default\n255.255.255.255\tbroadcasthost' >> /tmp/etchosts

  # Apply to /etc/hosts
  sudo mv /tmp/etchosts /etc/hosts
fi

# Flush DNS cache
sudo killall -HUP mDNSResponder
```
