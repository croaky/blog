# Restricted Postgres Access Between Render and Crunchy Bridge

I recently migrated a Ruby on Rails app from Heroku's Common Runtime and Heroku
Postgres to [Render and Crunchy Bridge](/webstack).

One reason to do this is improved networking security for the Postgres database.
Render and Crunchy Bridge support two networking options for free, or low extra
cost:

1. IP address firewall
2. Tailscale

## The attack vector

The problem I'm trying to solve is if my `DATABASE_URL` connection string
somehow leaks, I'd like extra protection that an attacker can't use it to access
my database without doing additional work.

Heroku Postgres' databases are available to the public internet (`0.0.0./0`),
which provides no additional protection in this scenario.
They offer solutions to this in [Heroku Private Spaces](https://www.heroku.com/private-spaces),
part of their enterprise offering, which can be a big increase in spending.

## Render

Services on Render send
[outbound traffic through a group of static IP addresses](https://render.com/docs/static-outbound-ip-addresses).

The IPs are very visible in the Render UI for each service.
You click "Connect > Outbound" and will see a group of 3 IP addresses
and a button to copy them.

The IPs are the same across services if the services are on the same region.
The IPs are not specific to services or environments.

For my purposes, this basic offering was an improvement over the previous setup
on Heroku. It replaced my use of the [Heroku Fixie](https://usefixie.com/)
add-on that I needed for an integration to an SFTP server in my app. And, it
allowed me to take my Postgres instance off the public internet.

Render, for now, does not offer outbound IPs that are unique-to-your project,
service, or environment. Something like Fixie or
[Quotaguard](https://render.com/docs/quotaguard) might work, but I have not
tested it. The documentation makes me think it will work for outbound HTTP
requests but not necessarily Postgres TCP connections.

## Render vs. Crunchy vs. Aiven for Postgres

The simplest option to secure networking between Render services and a Postgres
database is to host the Postgres instance on Render itself.
In a few clicks in Render's [Access Control forms](https://render.com/docs/databases#access-control),
you can disable the external database URL and prevent `0.0.0.0/0` connections.

For my purposes, I chose Crunchy over Render primarily because, as of June 2023,
Render does not have a HA (High Availability) Postgres option.
I also strongly considered [Aiven](https://aiven.io/) for managed Postgres
and would have chosen them if I had any other databases in my system
such as Redis, Kafka, or ElasticSearch.

Crunchy also has a few quality-of-life features like a built-in dashboard
showing index usage, a Heroku Dataclips-like feature called
[Queries](https://docs.crunchybridge.com/concepts/saved-queries),
and a Tailscale option (more below).

## Crunchy IP restrictions

The [Crunchy firewall rules](https://docs.crunchybridge.com/how-to/firewall/)
default to `0.0.0.0/0` and `::/0`.

Remove the Crunchy defaults
and add the 3 IPs from Render,
appending `/32` as a suffix to each
meaning "only allow this single, static IP address."

This can be done in the Crunchy UI or in one CLI command:

```
cb firewall --cluster CLUSTER_ID --remove 0.0.0.0/0 --remove ::/0 --add RENDER_IP_1/32 --add RENDER_IP_2/32 --add RENDER_IP_3/32
```

It's not really an "atomic" operation since they are still posted individually
but they are applied one right after the other and it's done faster than adding
each rule via the UI.

Existing database sessions connected will be unaffected by this change
but new sessions will be subject to the state of the firewall rules
when the connection was established.

## Crunchy Tailscale

This IP restriction setup is nice but still limited:

1. An attacker might know to try to access the database from the same Render
   region's IP addresses.
2. What if I want to add our corporate firewall's IP to the Crunchy firewall
   rules so I can connect to the database from our development machine? This IP
   address may also be known by an attacker.
3. IP addresses can be spoofed.

What I want is to have a private network across:

1. Render services
2. Crunchy Postgres databases
3. Laptop (development machine)

I like the idea of using Tailscale for this purpose. Here are docs
from both [Tailscale](https://tailscale.com/kb/1231/crunchy-bridge/)
and [Crunchy Bridge](https://docs.crunchybridge.com/how-to/tailscale).

I have tested accessing a Crunchy Postgres database using Tailscale
from my laptop and from [Fly.io](https://fly.io) and it seemed to work
exactly as I wanted.

It is possible to use Tailscale from Render,
but I need to specifically test Postgres-connections-over-Tailscale.

I have seen issues on [Heroku w/ Tailscale](https://github.com/croaky/webstack/tree/main/heroku-go-crunchy)
and [Railway w/
Tailscale](https://github.com/croaky/webstack/tree/main/railway-go-tailscale)
where I can install Tailscale and add the services to my Tailnet,
but I'm unable to connect to the Crunchy Postgres database.

I believe this may be due to how those platforms don't provide a `/dev/net/tun`
device, so they require Tailscale's [userspace networking
mode](https://tailscale.com/kb/1107/heroku/).

I might be able to change my application code to do something like deal with SOCKS5 proxies
but I don't think my Ruby Postgres driver supports this well
and it makes me very nervous changing such a critical point in my stack,

So, I have some open questions like:

- Will Render's Tailscale implementation work as smoothly as Fly.io
  where I can continue to use my simple `DATABASE_URL` connection in my
  application code, untouched?
- Will I need to switch off Render's [Native Runtimes](https://render.com/docs/native-runtimes)
  to a Docker deployment? This might be fine, or even provide some kind of
  benefits (faster builds?) but requires some exploration. There are no official
  Render + Tailscale docs published yet.
