# postgres / dev test clusters

I run separate Postgres clusters for development and test environments.
The test cluster has durability settings disabled
for fast automated test execution.

## Setup

I automate cluster setup in [cmd / laptop](/cmd/laptop):

```bash
brew install postgresql@17

export PATH="$BREW/opt/postgresql@17/bin:$PATH"
if ! command -v initdb >/dev/null || ! command -v pg_ctl >/dev/null; then
  echo "initdb and/or pg_ctl not found in PATH"
  exit 1
fi

start_postgres_cluster() {
  local port="$1"
  local data_dir="$2"
  local log_file="$3"
  local opts="$4"

  mkdir -p "$(dirname "$data_dir")"
  mkdir -p "$(dirname "$log_file")"

  if [ ! -f "$data_dir/PG_VERSION" ]; then
    initdb -D "$data_dir" -U postgres -c maintenance_work_mem=2GB
  fi

  if pg_ctl -D "$data_dir" status >/dev/null 2>&1; then
    echo "Postgres is already running for data directory $data_dir"
    return
  fi

  if lsof -i "tcp:$port" >/dev/null 2>&1; then
    echo "Port $port is already in use"
    return
  fi

  pg_ctl -D "$data_dir" -l "$log_file" -o "-p $port $opts" start
}

# dev databases
start_postgres_cluster 5432 \
  "$HOME/.local/share/postgres/data_dev" \
  "$HOME/.local/share/postgres/log_dev.log" \
  ""

# test databases
start_postgres_cluster 5433 \
  "$HOME/.local/share/postgres/data_test" \
  "$HOME/.local/share/postgres/log_test.log" \
  "-c fsync=off -c synchronous_commit=off -c full_page_writes=off"
```

The clusters start automatically when I run the laptop script.

## Test speed optimizations

The test cluster disables durability settings:

- `fsync=off`: No forced disk synchronization
- `synchronous_commit=off`: Transactions commit before writing to disk
- `full_page_writes=off`: Reduces write overhead

These settings trade data safety for speed, which is fine for test data.

## Usage

Set test database URLs as needed in test suites:

```
postgres://postgres@localhost:5433/app_test
```
