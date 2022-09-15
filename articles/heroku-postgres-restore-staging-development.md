# Heroku Postgres Restore to Staging and Development

I use [these bash scripts](https://github.com/croaky/blog/tree/main/code/heroku)
to restore Heroku Postgres data from production to staging
and from production to development environments.

They depend only on standard Unix tools, Postgres CLIs, and the Heroku CLI.
They can be customized per-project for pre- and post-processing.

## Restore to staging

The `db-restore-stag-from-prod-backup` script
pre- and post-processes data to prevent accidents:

```embed
code/heroku/db-restore-stag-from-prod-backup content
```

## Restore to development

The `db-download-prod-backup` script
downloads to the `tmp/latest.backup` file on my filesystem:

```embed
code/heroku/db-download-prod-backup content
```

I restore `tmp/latest.backup`, pre- and post-process in
the `db-restore-dev-from-downloaded-backup` script:

```embed
code/heroku/db-restore-dev-from-downloaded-backup content
```
