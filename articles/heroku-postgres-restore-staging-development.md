# Heroku Postgres Restore to Staging and Development

I use these scripts when it is okay
to copy Heroku Postgres data between production, staging,
and development environments:

```
db-restore-stag-from-prod-backup
db-backup-stag
db-download-stag-backup
db-restore-dev-from-downloaded-backup
```

## db-restore-stag-from-prod-backup

```embed
code/heroku/db-restore-stag-from-prod-backup content
```

Specific to each project, that script can sanitize data or
disable users or feature flags to help prevent accidental notifications.

## db-backup-stag

```embed
code/heroku/db-backup-stag content
```

## db-download-stag-backup

```embed
code/heroku/db-download-stag-backup content
```

I keep the `tmp/latest.backup` file on my filesystem
so I can restore locally at any time using the next script.

## db-restore-dev-from-downloaded-backup

```embed
code/heroku/db-restore-dev-from-downloaded-backup content
```

This script is another opportunity to run SQL commands
to flip feature flags, make my user an admin user, etc.

## Why not Parity?

I am the original author of the
[Parity Ruby gem](https://github.com/thoughtbot/parity),
which provided similar functionality.
I switched to these scripts because:

* some projects shouldn't have a Ruby dependency
* shell scripts can be customized easily for post-processing
* additional security and bug avoidance by hard-coding app names instead of
  going through the indirection of named `staging` and `production` Git remotes

For the same reasons,
I switched from Parity's `production` and `staging` commands
to these scripts:

```
stag
prod
```

## stag

```embed
code/heroku/stag content
```

## prod

```embed
code/heroku/prod content
```
