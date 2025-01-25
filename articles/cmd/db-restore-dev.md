# cmd / db-restore-dev

I frequently download and restore my production database to my laptop
using scripts which are placed in my project's Git repo.

They depend on Unix, Postgres, and [Crunchy
Bridge](https://docs.crunchybridge.com/concepts/cli/) CLIs.

The [db-download-prod](/cmd/db-download-prod) script
downloads the backup to `tmp/latest.backup`.

The `db-restore-dev` script restores from the `tmp/latest.backup` file
and does custom post-processing as needed:

```bash
#!/bin/bash
set -euo pipefail

db="app_dev"

dropdb --if-exists "$db"
createdb "$db"
psql "$db" <<SQL
  CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
  CREATE EXTENSION IF NOT EXISTS pg_trgm;
  CREATE EXTENSION IF NOT EXISTS plpgsql;
SQL

# Same directory defined in `bin/db-download-prod`
backup_dir="tmp/latest_backup_dir"

# Detect the number of CPU cores
case "$(uname -s)" in
    Linux*)     cores=$(nproc);;
    Darwin*)    cores=$(sysctl -n hw.ncpu);;
    *)          cores=1;;
esac

# Use one less than the total number of cores, but ensure at least 1 is used
(( jobs = cores - 1 ))
if (( jobs < 1 )); then
    jobs=1
fi

echo "Restoring with $jobs parallel job(s)"

# Restore from directory
pg_restore -d "$db" --verbose --no-acl --no-owner -j "$jobs" "$backup_dir"

# Post-process
psql "$db" <<SQL
  UPDATE ar_internal_metadata
  SET value = 'development'
  WHERE key = 'environment';

  -- Avoid re-running incomplete jobs
  DELETE FROM jobs
  WHERE status IN ('pending', 'started');

  -- Avoid emailing production users
  UPDATE users
  SET active = false;

  -- Turn on flags for developers
  UPDATE
    users
  SET
    active = true,
    admin = true
  WHERE
    email IN (
      'dev1@example.com',
      'dev2@example.com'
    );
SQL
```
