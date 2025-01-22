# cloud / render

From 2009-2023, I most frequently deployed my database-backed web apps to
[Heroku](https://heroku.com).
In 2023, I wanted to move to a new provider.
I created a [webstack](https://github.com/croaky/webstack) repo
to prototype options.

Based on the users of the primary web app I work on,
I set up API checks via [Checkly](https://www.checklyhq.com/) from
Northern California and London.

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

While performance and reliability have been good on Heroku,
I wanted to migrate off due to the platform's offerings beginning to stagnate.
My confidence was particularly shaken when its GitHub integration broke in April
2022 and it took over a month to be resolved.

Heroku lacked some long overdue features such as HTTP/2 support and
it was not possible to restrict access to its Postgres database
from public internet without a major increase in cost.

## Render

For the app I work on, I ended up choosing
[Render](https://render.com)
and [Crunchy Bridge](https://crunchybridge.com/) Postgres.

This combo felt like the smallest step from Heroku.
I felt both companies were mature orgs and my software would be reliable.

Some things I liked about Render:

- Good customer support.
- [IP access control on Postgres databases](https://render.com/docs/databases)
- [DDoS protection](https://render.com/docs/ddos-protection)
- HTTP/3 and HTTP/2 support.
- [Zero-downtime deploys via health checks](https://render.com/docs/deploys#zero-downtime-deploys)
- SOC2 certified.

## Recommendations for others

Consider [Fly.io](https://fly.io) if your priorities are multi-region,
low-latency, or lowest cost.

Consider [Northflank](https://northflank.com/) if your users are mainly in
Europe, or if your mindset is particularly oriented around Docker or Kubernetes.

Consider [Railway](https://railway.app) if your priorities are fast build times
or you want to interact with your infrastructure primarily via web UI.

Consider [Aiven](https://aiven.com/) if you have multiple database types you
want managed such as some combination of Postgres, Redis, Kafka, or
ElasticSearch.

Consider Crunchy Bridge if you are "all-in" on Postgres and want your database
tooling and support to be optimized around Postgres.
