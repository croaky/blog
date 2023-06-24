# Postgres Download Prod, Restore Dev

I restore my development database from a production copy daily using these bash
scripts. They depend on standard Unix tools, Postgres, and [Crunchy
Bridge](https://docs.crunchybridge.com/concepts/cli/) CLIs. They can be
customized per-project for pre- and post-processing.

The `db-download-prod` script
downloads to the `tmp/latest.backup` file on my filesystem:

```embed
code/postgres/db-download-prod content
```

I restore `tmp/latest.backup` and post-process in
the `db-restore-dev` script:

```embed
code/postgres/db-restore-dev content
```
