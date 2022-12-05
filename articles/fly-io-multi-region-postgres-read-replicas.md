# Fly.io Multi-Region Postgres Read Replicas

For improved latency,
you might want to deploy your app to one or more of
[Fly.io's 20+ regions](https://fly.io/docs/reference/regions/)
on their own Points of Presence.

You can also pair your regional servers with Fly's
[multi-region Postgres read replicas](https://fly.io/docs/getting-started/multi-region-databases/).
The caveat is
[Fly Postgres is not managed](https://fly.io/docs/rails/getting-started/migrate-from-heroku/#databases).

## Develop

Copy this as `main.rb` and run the commands in the comments:

```embed
code/fly-read-replica/main.rb
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

Create a Fly app:

```bash
fly launch --remote-only
```

Tell the app which region will be primary,
such as Sunnyvale, California:

```bash
fly secrets set PRIMARY_REGION=sjc
```

Create Fly Postgres database in that region:

```bash
fly pg create --name db-replicas --region sjc
```

Set the `DATABASE_URL` environment variable
by attaching the database cluster to the app:

```bash
fly pg attach db-replicas
```

For improved latency to the app servers, add regions
such as Frankfurt, Germany and Sydney, Australia:

```bash
fly regions add fra syd
```

Create Postgres read replicas in those regions:

```bash
fly volumes create pg_data -a example-read-replicas-db --size 1 --region fra
fly volumes create pg_data -a example-read-replicas-db --size 1 --region syd
```

Scale 1 instance per region:

```bash
fly autoscale standard min=3 max=3
```

Done!
