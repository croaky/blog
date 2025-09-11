# postgres / restore dev

I frequently download and restore
my production Postgres database to my laptop
using scripts which are placed in my project's Git repo.

They depend on Unix, Postgres,
and [Crunchy Bridge](https://docs.crunchybridge.com/concepts/cli/) CLIs.

The [`db-download-prod`](/postgres/download-prod) script
downloads the backup to `tmp/latest.backup`.
