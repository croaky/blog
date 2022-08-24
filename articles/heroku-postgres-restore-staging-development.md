# Heroku Postgres Restore to Staging and Development

I use [these scripts](https://github.com/croaky/blog/tree/main/code/heroku)
to restore Heroku Postgres data from production to staging
and from production to development environments.

## Restore to staging

Specific to each project, the `db-restore-stag-from-prod-backup` script
can sanitize data or flip feature flags to help prevent accidents:

```embed
code/heroku/db-restore-stag-from-prod-backup content
```

## Restore to development

The `db-download-prod-backup` script
downloads to the `tmp/latest.backup` file on my filesystem:

```embed
code/heroku/db-download-prod-backup content
```

I can restore `tmp/latest.backup` at any time using
the `db-restore-dev-from-downloaded-backup` script,
which is another opportunity to flip feature flags,
make my user an admin user, etc.

```embed
code/heroku/db-restore-dev-from-downloaded-backup content
```

## Why not Parity?

I authored [Parity](https://github.com/thoughtbot/parity)
but switched to these scripts because:

* some projects shouldn't have a Ruby dependency
* shell scripts can be customized easily for post-processing
* separating the "download" and "restore" steps saves time
  when I only need to do one or the other
* improve security and avoid bugs by hard-coding Heroku app names
  instead of indirectly using `staging` and `production` Git remotes

For similar reasons,
I switched from Parity's `staging` and `production` commands
to the `stag` script:

```embed
code/heroku/stag content
```

And `prod` script:

```embed
code/heroku/prod content
```
