# Heroku Postgres Restore to Staging and Development

I use [these scripts](https://github.com/croaky/blog/tree/main/code/heroku)
to restore Heroku Postgres data from production to staging
and from production to development environments.

## Restore to staging

```embed
code/heroku/db-restore-stag-from-prod-backup content
```

Specific to each project, that script can sanitize data or
disable users or feature flags to help prevent accidental notifications.

## Restore to development

```embed
code/heroku/db-download-prod-backup content
```

I keep the `tmp/latest.backup` file on my filesystem
so I can restore locally at any time using the next script.

```embed
code/heroku/db-restore-dev-from-downloaded-backup content
```

This script is another opportunity to run SQL commands
to flip feature flags, make my user an admin user, etc.

## Why not Parity?

I authored [Parity](https://github.com/thoughtbot/parity)
but switched to these scripts because:

* some projects shouldn't have a Ruby dependency
* shell scripts can be customized easily for post-processing
* improve security and avoid bugs by hard-coding Heroku app names
  instead of indirectly using `staging` and `production` Git remotes

For similar reasons,
I switched from Parity's `production` and `staging` commands
to these scripts:

```embed
code/heroku/stag content
```

```embed
code/heroku/prod content
```
