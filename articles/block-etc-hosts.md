# Block with /etc/hosts

I run a shell script on my laptop to block ads, trackers,
and malicious websites at the DNS host level.
I also use 1.1.1.1 as the DNS resolver on my laptop and phone.
This article describes why, alternatives, and trade-offs.

## Browser extensions

Ad blockers such as Adblock and AdBlock Plus (different companies)
are installed as browser extensions.
They are installed per-browser, per-device.
Like any browser extension,
they can read the code of every site you browse.

## DNS sinkholes

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

[NextDNS](https://nextdns.io/) is a hosted version that works
anywhere, including on a mobile phone.
It requires sending all your traffic through them.

## Modify /etc/hosts

Here's [the script](https://github.com/croaky/laptop/blob/main/bin/adblock)
I use (modified from a version by [@djcp](https://twitter.com/djcp)):

```bash
#!/bin/bash

set -eo pipefail

# Create file to block ads at the networking level
curl -s http://winhelp2002.mvps.org/hosts.txt > /tmp/adblock

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

# Flush DNS cache
sudo killall -HUP mDNSResponder
```

The data source is [MVPS HOSTS](http://winhelp2002.mvps.org/hosts.txt).
It is free to use for personal use and licensed under
the Creative Commons Attribution-NonCommercial-ShareAlike License.

It is fast to set up and run.
It can be edited to allow specific hosts,
such as those I need for work.
It works at home or when traveling.
It only works on my laptop, not my mobile phone.

## 1.1.1.1 as DNS resolver

I set my laptop's and phone's DNS resolver to [`1.1.1.1`](https://1.1.1.1),
a [fast, privacy-focused](https://blog.cloudflare.com/announcing-1111/)
consumer DNS service from Cloudflare.

On macOS, this setting can be controlled by going to
"System Preferences > Network > Advanced... > DNS",
clicking "+", entering "1.1.1.1", clicking "OK",
and clicking "Apply".

On iPhone, this setting can be controlled by a nice
[iPhone app](https://apps.apple.com/us/app/1-1-1-1-faster-internet/id1423538627).

## Conclusion

On my laptop, the internet is fast, private, and mostly ad-free.
On my phone, the internet is fast and private.
This approach is a middle ground without extra cost
and little to maintain.
