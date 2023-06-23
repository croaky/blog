# Webstack

From 2009-2023, my web stack most often included Heroku
but I wanted to move to a new provider.
I created a [webstack](https://github.com/croaky/webstack) repo
to prototype options.

Based on the users of the primary web app I work on,
I set up API checks via [Checkly](https://www.checklyhq.com/) from Northern California and London.

The hosting providers I tested the most were:

- [Fly.io](https://fly.io)
- [Heroku](https://heroku.com)
- [Northflank](https://northflank.com/)
- [Railway](https://railway.app)
- [Render](https://render.com)

The Postgres databases I tested the most were:

- [Aiven](https://aiven.com/)
- [Crunchy Bridge](https://crunchybridge.com/)
- Fly.io
- Heroku
- [Neon](https://neon.tech)
- Northflank
- Railway
- Render

I also spent some time with Vercel, PlanetScale, Supabase, Cockroach, and Fly.io SQLite.
See the repo for details on those; they aren't covered in this article.

Each stack served a healthcheck-style HTTP API endpoint that executed a
`SELECT 1` to a SQL database and responded with JSON `{"status":"ok"}`.
Each stack used a lightweight router, a SQL database driver (no ORM),
and a database connection pool. Most stacks were written in Go. Example:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// env
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		dbUrl = "postgres:///webstack_dev"
	}

	// db
	db, err := pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var col int
		db.QueryRow(r.Context(), "SELECT 1").Scan(&col)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"status\":\"ok\"}")
	})

	// listen
	log.Println("Listening at http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

Neon and Crunchy Bridge offer connection pooling via PgBouncer, which I enabled.
See [Neon docs](https://neon.tech/docs/get-started-with-neon/connection-pooling/),
and [Crunchy Bridge docs](https://docs.crunchybridge.com/how-to/pgbouncer/).
In Neon, you enable it with a web UI toggle.
In Crunchy Bridge, you enable it by \`psql\`'ing into your cluster and running `CREATE EXTENSION crunchy_pooler;`.

## Heroku

My baseline Heroku workflow has historically been:

- Version code in a GitHub repo
- Open GitHub pull requests w/ [CI](https://www.thoughtworks.com/continuous-integration)
- Merge into `main` branch to auto-deploy to staging environment
- Promote staging manually to production environment in a
  [pipeline](https://devcenter.heroku.com/articles/pipelines)
- Configure production database with a
  [high availability
  follower](https://devcenter.heroku.com/articles/heroku-postgres-follower-databases)
  and continuous backups

While performance and reliability have been very good on Heroku,
I wanted to migrate off due to the platform's offerings beginning to stagnate.
My confidence was particularly shaken when its GitHub integration broke in April
2022 and it took over a month to be resolved.

Heroku lacked some long overdue features such as HTTP/2 support and
it was not possible to restrict access to its Postgres database
from public internet without a major increase in cost.

## Fly.io

Pros:

- I was able to best performance of any platform I tested.
- HTTP/2 support.
- The CLI is fantastic.
- SOC2 certified.

Cons:

- [Fly's Postgres databases are not
  managed](https://fly.io/docs/rails/getting-started/migrate-from-heroku/#databases),
  which requires extra monitoring and operating work.
- They are undergoing a platform architecture migration to "Fly Apps v2" that
  makes it awkward to navigate the docs and CLI output. While I had good uptime in
  my testing, there have been many reports from others on their forum that
  reliability hasn't been great during the Apps v2 migration.
- The web UI is not as aesthetic as Heroku, Northflank, Railway, or Render.
- Slower builds than Railway.

## Northflank

Pros:

- Good customer support.
- Comprehensive feature set.
- Excellent UI.
- Build with Docker, Heroku buildpacks, or Paketo buildpacks.
- Deploy to US Central or Europe West.
- Nice organization of infra into projects and deployment via staging-production pipeline.
- Managed Postgres.
- HA Postgres databases.
- Hourly Postgres backups.
- Postgres read replicas.
- Postgres can be networked to be hidden from the internet.
- HTTP/2 works, TLS terminates at Northflank's routing layer.

Cons:

- Performance didn't seem as good in my testing as Fly or Render.
- I could get `heroku/buildpacks:20` working, but not `heroku/builder:22`.
- Unlike Fly and Render, not SOC2 certified.
- Unlike Fly, no multi-region deployments.
- Unlike Render, no built-in DDoS protection.
- Slower builds than Railway.

## Railway

Pros:

- The web UI is absolutely gorgeous.
- Fast build times via [Nixpacks](https://docs.railway.app/deploy/builds).
- Web services and Postgres are easy to set up and connect to each other. The UI
  provides a great visualization of their relationship.
- Fast pull request environments.
- Nice in-browser SQL editor / console.
- Philosophically aligned with my general mentality of wanting to not think in containers.

Cons:

- The only region supported so far is US West. [Other regions are
  planned](https://feedback.railway.app/feature-requests/p/configurable-deployment-region)
- Postgres databases are exposed to the public internet like Heroku Postgres on
  Heroku's Common Runtime. [Private networking is
  planned](https://feedback.railway.app/feature-requests/p/internal-networking)
- HA Postgres databases are not available.

## Render

Pros:

- Good customer support.
- Managed Postgres.
- [IP access control on Postgres databases](https://render.com/docs/databases)
- [DDoS protection](https://render.com/docs/ddos-protection)
- HTTP/3 and HTTP/2 support.
- [Zero-downtime deploys via health checks](https://render.com/docs/deploys#zero-downtime-deploys)
- SOC2 certified.

Cons:

- HA Postgres databases are not available.
- Slower builds than Railway.
- Unlike Fly, no multi-region deployments.

## My stack choice

For the app I work on, I ended up choosing Render and Crunchy Bridge Postgres.

It felt like the smallest step from Heroku.
I felt both companies were mature organizations and reliability for my software would be good.
I am "all-in" on Postgres; I use it as my queuing system and have no other databases such as Redis.

## Recommendations for others

Consider Fly if your priorities are multi-region, low-latency, or lowest cost.

Consider Northflank if your app and users are mainly in Europe, or if your
mindset is particularly oriented around Docker or Kubernetes.

Consider Railway if your priorities are fast build times or you want to interact
with your infrastructure primarily via web UI.

Consider Aiven if you have multiple database types you want managed
such as some combination of Postgres, Redis, Kafka, or ElasticSearch.

Consider Crunchy Bridge if you are "all-in" on Postgres and want your database
tooling and support to be optimized around Postgres.
