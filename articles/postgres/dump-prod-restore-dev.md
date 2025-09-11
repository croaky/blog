# postgres / dump prod restore dev

I frequently dump my production Postgres database
and restore to my development machine using scripts
stored in my project's Git repo.

The current version depends on Unix, Postgres,
and [Crunchy Bridge](https://docs.crunchybridge.com/concepts/cli/) CLIs.

The `db-download-prod` script
downloads the backup to `tmp/latest.backup`:

```bash
#!/bin/bash
set -euo pipefail

# Delete/create target directory
backup_dir="tmp/latest_backup_dir"
rm -rf "$backup_dir"
mkdir -p "$backup_dir"

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

echo "Downloading with $jobs parallel job(s)"

# Use the directory format and specify the number of jobs for parallel dumping
pg_dump -Fd "$(cb uri app-prod --role application)" -j "$jobs" -f "$backup_dir"
```

The `db-restore-dev` script restores from backup files
and post-processes:

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

I separate the scripts so I can restore
a recent backup without re-downloading.
