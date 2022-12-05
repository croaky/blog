# Fly.io Multi-Region Ruby with Postgres Read Replicas

For latency reasons,
you might want to deploy your app to a specific region.
[Fly.io has 20+ regions](https://fly.io/docs/reference/regions/)
on their own Points of Presence (not on public clouds like AWS).

If you have a read-heavy application, look at Fly's
[multi-region Postgres read replicas](https://fly.io/docs/getting-started/multi-region-databases/).
The caveat is
[Fly Postgres is not managed](https://fly.io/docs/rails/getting-started/migrate-from-heroku/#databases).

## Ruby app

In this example, we'll write and deploy a Ruby app to Fly.io:

```embed
code/fly-read-replica/main.rb
```

## Develop

Run locally:

```bash
createdb db
chmod +x main.rb
DATABASE_URL=postgres:///db ./main.rb
```

## Fly

Install the `flyctl` CLI, which is symlinked as `fly`:

```bash
brew install flyctl
```

Connect to your account:

```
fly auth login
```

If you don't have Docker installed,
use Fly's `--remote-only` option,
which will use a remote Docker builder in your Fly account:

```bash
fly launch --remote-only
```

If you have Docker installed,
deploy with Fly's `--nixpacks` option,
which will use [nixpacks](https://github.com/railwayapp/nixpacks) as the builder.
It deploys faster:

```bash
fly launch --nixpacks
```

Configure region the primary database location,
such as Sunnyvale, California:

```bash
fly secrets set PRIMARY_REGION=sjc
```

Add the regions where the read replicas will go,
such as Frankfurt, Germany and Sydney, Australia:

```bash
fly regions add fra syd
```

Create read replicas in those regions:

```bash
fly volumes create pg_data -a example-read-replicas-db --size 1 --region fra
fly volumes create pg_data -a example-read-replicas-db --size 1 --region syd
```

Scale 1 instance per region:

```bash
fly autoscale standard min=3 max=3
```

Done!
