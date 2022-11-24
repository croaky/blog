# Postgres concat_ws() Function

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

Using `concat_ws()` instead of `concat()` avoids output like `AB / / CD`.
