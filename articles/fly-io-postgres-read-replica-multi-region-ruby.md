# Fly.io Multi-Region Ruby with Postgres Read Replicas

For latency or data residency reasons,
you might want to deploy your web app to a specific region.
[Fly.io has 20+ regions](https://fly.io/docs/reference/regions/)
on their own Points of Presence
(not on public clouds like AWS).

If you have a read-heavy application, look at Fly's
[multi-region Postgres read
replicas](https://fly.io/docs/getting-started/multi-region-databases/).
The big caveat is that [Fly, unlike Heroku, is not managed
Postgres](https://fly.io/docs/rails/getting-started/migrate-from-heroku/#databases).

## Ruby app

In this example, we'll build a Ruby web app using three Ruby gems:

* [Connection Pool](https://rubygems.org/gems/connection_pool)
* [PG](https://rubygems.org/gems/pg)
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
flyctl auth login
```

I prefer to not install Docker on my local machine.
One "gotcha" is a leftover `~/.docker` directory can cause errors on Fly.
Remove it:

```bash
rm -rf ~/.docker
```

Now we can use Fly's `--remote-only` option.
It will use a remote Docker builder that they set up in our account:

```bash
fly launch --remote-only
```

Configure primary region where the primary database is:

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
