# Fly.io Multi-Region Ruby with Postgres Read Replicas

For latency (or, less relevant to this article, data residency) reasons,
you might want to deploy your app to a specific region.
[Fly.io has 20+ regions](https://fly.io/docs/reference/regions/)
on their own Points of Presence
(not on public clouds like AWS).

If you have a read-heavy application, look at Fly's
[multi-region Postgres read replicas](https://fly.io/docs/getting-started/multi-region-databases/).
The big caveat is
[Fly, unlike Heroku, is not managed Postgres](https://fly.io/docs/rails/getting-started/migrate-from-heroku/#databases).

## Ruby app

In this example, we'll write and deploy a Ruby app to Fly.io using popular gems:

* [Connection Pool](https://rubygems.org/gems/connection_pool)
* [PG](https://rubygems.org/gems/pg)
* [Puma](https://rubygems.org/gems/puma)
* [Sinatra](https://rubygems.org/gems/sinatra)

The `Gemfile` looks like this:

```embed
code/fly-ruby-read-replica/Gemfile all
```

Install the gems:

```bash
bundle
```

The `api.rb` looks like this:

```embed
code/fly-ruby-read-replica/api.rb all
```

## Develop

Develop locally:

```bash
bundle
createdb example_dev
DATABASE_URL="postgres:///example_dev" bundle exec ruby api.rb
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

Deploy using Fly's `--remote-only` option,
which will use a remote Docker builder Fly sets up in your account.
If you previously had a Docker installation, avoid
[a gotcha](https://community.fly.io/t/failed-to-fetch-builder-image-when-deploying-ruby-example/3726/7),
by deleting `~/.docker` before you deploy:

```bash
rm -rf ~/.docker
fly launch --remote-only
```

Configure region where the primary database is:

```bash
fly secrets set PRIMARY_REGION=sjc
```

Create the regions where the read replicas will go:

```bash
fly regions add ams lhr syd yul
```

Create read replicas in those regions:

```bash
fly volumes create pg_data -a example-read-replicas-db --size 1 --region ams
fly volumes create pg_data -a example-read-replicas-db --size 1 --region lhr
fly volumes create pg_data -a example-read-replicas-db --size 1 --region syd
fly volumes create pg_data -a example-read-replicas-db --size 1 --region yul
```

Scale 1 instance per region:

```bash
fly autoscale standard min=5 max=5
```

Done!
