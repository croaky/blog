# Dev data

I frequently download and restore my production database to my laptop
using the following scripts, which are placed in my project's Git repo.

They depend on standard Unix tools, Postgres, and [Crunchy
Bridge](https://docs.crunchybridge.com/concepts/cli/) CLIs.

The `db-download-prod` script
downloads the backup to `tmp/latest.backup`:

```embed
code/postgres/db-download-prod.sh
```

The `db-restore-dev` script restores from the `tmp/latest.backup` file
and does custom post-processing as needed for the project:

```embed
code/postgres/db-restore-dev.sh
```
