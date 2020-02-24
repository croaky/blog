# Block Ads with /etc/hosts

I run a small shell script
on my laptop to block ads at the DNS host level.
I also use DNS 1.1.1.1 on my laptop and mobile phone.
This article describes why, alternatives, and trade-offs.

## Browser extensions

Ad blockers such as Adblock and AdBlock Plus (different companies)
are installed as browser extensions.
They must be installed in each browser on each of your devices.

Like any browser extension,
they can read the HTML of every site you browse.
You might trust them to sensitively handle all that personal content
but there's no need to, given other options.

## DNS sinkholes

Ad blockers such as [Pi-hole](https://pi-hole.net/)
are installed as DNS sinkholes.
They block ads on more devices
(laptops, mobile phones, tablets, smart TVs)
and from more apps than only web browsers.

DNS sinkholes must be run as a server on an always-on device somewhere.
Many users run them on their LAN or a Raspberry Pi ($35) at home.
When using a device away from home,
ad-blocking may work (with additional setup),
work partially (depending on caching)
or not all.

This is the best overall solution
but requires the most technical ability, time, and cost.

## Modify /etc/hosts

Here's [the script](https://github.com/croaky/laptop/blob/master/bin/adblock)
I use:

```bash
#!/bin/bash

set -eo pipefail

# Get /etc/hosts file to block ads at the networking level
curl -s http://winhelp2002.mvps.org/hosts.txt > /tmp/adblock

# Re-write Windows to Unix line endings
tr -d '\r' < /tmp/adblock > /tmp/etchosts

comment() {
  replace "0.0.0.0 $1" "# 0.0.0.0 $1" /tmp/etchosts
}

# Comment-out used hosts
comment 'ads.twitter.com'
comment 'api.segment.io'
comment 'static.ads-twitter.com'

# Add /etc/hosts overrides
echo '# MacOS default
255.255.255.255 broadcasthost' >> /tmp/etchosts

# Update /etc/hosts
sudo mv /tmp/etchosts /etc/hosts

# Flush DNS cache
sudo killall -HUP mDNSResponder
```

The data source is [MVPS HOSTS](http://winhelp2002.mvps.org/hosts.txt).
It is free to use for personal use and licensed under
the Creative Commons Attribution-NonCommercial-ShareAlike License.

It is fast to set up, run, and edit.
It can be modified to allow specific hosts,
which I sometimes need depending on my work.
It works at home or when traveling.

While free, fast, and flexible, it only works on my laptop, not my mobile phone.

## DNS to 1.1.1.1

I also have my DNS on my laptop and phone set to `1.1.1.1`,
a fast, privacy-focused consumer DNS service from Cloudflare.
Read [their 2018 announcement](https://blog.cloudflare.com/announcing-1111/)
to learn about its encrypted DNS benefits
and [1.1.1.1](https://1.1.1.1) to see its speed benefits
compared to other DNS resolvers.

On macOS, this setting can be controlled by going to
"System Preferences > Network > Advanced... > DNS",
clicking "+", entering "1.1.1.1", clicking "OK",
and clicking "Apply".

On iPhone, this setting can be controlled via a beautiful
[iPhone app](https://apps.apple.com/us/app/1-1-1-1-faster-internet/id1423538627).

## Conclusion

I've landed on this approach as comfortable middle ground.
My internet access is faster and more private than it would otherwise be.
On my laptop, it is mostly ad-free.
It doesn't cost anything and there is little to maintain.
