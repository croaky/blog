# Postgres tips

Here are a few quick Postgres tips

## concat_ws

Postgres provides a
[`concat_ws()`](https://www.postgresql.org/docs/current/functions-string.html)
string function with this signature:

```
concat_ws(sep text, str "any" [, str "any" [, ...] ])
```

Consider a schema where `projects.second_user_id`
and `project.third_user_id` columns can be `NULL`:

```
              Table "public.projects"
     Column     |            Type        | Nullable |
----------------+------------------------+-----------
 id             | bigint                 | not null |
 name           | character varying(255) | not null |
 lead_user_id   | bigint                 | not null |
 second_user_id | bigint                 |          |
 third_user_id  | bigint                 |          |

                Table "public.users"
     Column     |            Type        | Nullable |
----------------+------------------------+-----------
 id             | bigint                 | not null |
 initials       | character varying(255) | not null |
```

A query to get the team for a project by each user's initials:

```sql
SELECT
  projects.name,
  concat_ws(' / ', u1.initials, u2.initials, u3.initials) AS team
FROM
  projects
  LEFT JOIN users u1 ON u1.id = projects.lead_user_id
  LEFT JOIN users u2 ON u2.id = projects.second_user_id
  LEFT JOIN users u3 ON u3.id = projects.third_user_id
GROUP BY
  projects.name,
  team
ORDER BY
  projects.name ASC;
```

Example output:

```
     name     |     team
--------------+--------------
Private Beta  | AB / CD
SLA           | EF / GH / IJ
Treasury Ops  | KL
(3 rows)

Time: 1 ms
```

Using `concat_ws()` instead of `concat()` prevents `AB / / CD`.

## Create indexes concurrently

By default,
Postgres' `CREATE INDEX` locks writes (but not reads) to a table.
That can be unacceptable during a production deploy.
On a large table, indexing can take hours.

Postgres has a [`CONCURRENTLY` option for `CREATE INDEX`](https://www.postgresql.org/docs/current/sql-createindex.html)
that creates the index without preventing concurrent
`INSERT`s, `UPDATE`s, or `DELETE`s on the table.

One caveat is that
[concurrent indexes must be created outside a transaction](https://www.postgresql.org/docs/current/sql-createindex.html#SQL-CREATEINDEX-CONCURRENTLY).

If you want to do this in ActiveRecord:

```ruby
class AddIndexToAsksActive < ActiveRecord::Migration
  disable_ddl_transaction!

  def change
    add_index :asks, :active, algorithm: :concurrently
  end
end
```

The `disable_ddl_transaction!` method applies only to that migration file.
Adjacent migrations still run in their own transactions
and roll back automatically if they fail.
Therefore, it's a good idea to isolate concurrent index migrations
to their own migration files.
