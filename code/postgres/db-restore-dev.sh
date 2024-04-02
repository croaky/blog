#!/bin/bash
set -euo pipefail

db="app_development"
dropdb --if-exists "$db"
createdb "$db"
psql "$db" -c "
  CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
  CREATE EXTENSION IF NOT EXISTS pg_trgm;
  CREATE EXTENSION IF NOT EXISTS plpgsql;
"
pg_restore tmp/latest.backup --verbose --no-acl --no-owner --dbname "$db"
psql "$db" -c "
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
"
