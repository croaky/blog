# Block with /etc/hosts

To improve speed, privacy, and safety on my laptop,
I run a shell script to block ads, trackers,
and malicious websites at the DNS host level.
I also use 1.1.1.1 as my DNS resolver.
This article describes why, alternatives, and trade-offs.

## Alternatives

Ad blockers such as Adblock and AdBlock Plus (different companies)
are installed as browser extensions.
They are installed per-browser, per-device.
Like any browser extension,
they can read the code of every site you browse.

Ad blockers such as [Pi-hole](https://pi-hole.net/)
are installed as DNS sinkholes.
They block ads on all apps (not only web browsers)
on all devices (laptops, phones, tablets) on the network.

DNS sinkholes require technical ability, time, and cost.
They must be run as a server on an always-on device somewhere
such as a Raspberry Pi ($35) at home.
On a device away from home,
ad-blocking may work with additional setup,
work partially depending on caching,
or not at all.

## Script

Here's [the script](https://github.com/croaky/laptop/blob/main/bin/adblock)
I use:

```bash
#!/bin/bash

set -eo pipefail

if [[ "$1" == "undo" ]]; then
  echo '# MacOS default
  255.255.255.255 broadcasthost' | sudo tee /etc/hosts > /dev/null
else
  # Create file to block ads at the networking level
  curl -s https://winhelp2002.mvps.org/hosts.txt > /tmp/adblock

  # Re-write Windows to Unix line endings
  tr -d '\r' < /tmp/adblock > /tmp/etchosts

  comment() {
    replace "0.0.0.0 $1" "# 0.0.0.0 $1" /tmp/etchosts
  }

  # Comment-out used hosts
  comment 'api.amplitude.com'
  comment 'api.segment.io'

  # Restore macOS system defaults
  echo '# MacOS default
  255.255.255.255 broadcasthost' >> /tmp/etchosts

  # Apply to /etc/hosts
  sudo mv /tmp/etchosts /etc/hosts
fi

# Flush DNS cache
sudo killall -HUP mDNSResponder
```

The data source is [MVPS HOSTS](https://winhelp2002.mvps.org/hosts.txt).
It is free to use for personal use and licensed under
the Creative Commons Attribution-NonCommercial-ShareAlike License.

It is fast to set up and run.
It works at home or when traveling.
It only works on my laptop, not my mobile phone.
It can be edited to allow specific hosts.
I can disable and re-enable it:

```bash
adblock undo
adblock
```

## 1.1.1.1 as DNS resolver

I also set my laptop's and phone's DNS resolver to [`1.1.1.1`](https://1.1.1.1),
a [fast, privacy-focused](https://blog.cloudflare.com/announcing-1111/)
consumer DNS service from Cloudflare.

On macOS, this setting can be controlled by going to
"System Preferences > Network > Advanced... > DNS",
clicking "+", entering "1.1.1.1", clicking "OK",
and clicking "Apply".
