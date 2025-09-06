# ruby / db

I wrap the [pg](https://github.com/ged/ruby-pg) Postgres driver
with light methods instead of using an ORM.

## Configuration

Add the pg gem:

```ruby
gem "connection_pool"
gem "pg"
```

Configure the pool on process boot:

```ruby
# config/puma.rb
on_worker_boot do
  DB.configure do |c|
    c.pool_size = workers * threads
    c.reap = true
  end
end
```

## Wrapper

The DB class manages connection pooling and provides a simple interface:

```ruby
# lib/db.rb
require "connection_pool"
require "pg"

class DB
  Config = Struct.new(:pool_size, :reap)

  class << self
    def configure
      @config = Config.new
      yield(@config)
    end

    def pool
      @pool ||= new(
        pool_size: @config.pool_size,
        reap: @config.reap
      )
    end
  end

  def initialize(pool_size: 1, reap: false)
    @pool = ConnectionPool.new(size: pool_size, timeout: 5) { build_pg_conn }

    if reap
      start_reaper_thread
    end
  end

  def exec(sql, params = [])
    @pool.with do |conn|
      iconn = InstrumentedConnection.new(conn, @env)
      iconn.exec(sql, params)
    end
  end

  def transaction
    @pool.with do |conn|
      iconn = InstrumentedConnection.new(conn, @env)
      iconn.transaction { yield(iconn) }
    end
  end

  def fuzzy_like(string)
    pattern = Regexp.union("%", "_")
    escape_char = "\\"
    escaped = string
      .gsub(pattern) { |x| [escape_char, x].join }
      .tr(" ", "%")

    "%#{escaped}%"
  end

  private def build_pg_conn
    conn = PG.connect(db_url_for(@env))

    map = PG::BasicTypeMapForResults.new(conn)
    map.default_type_map = PG::TypeMapAllStrings.new

    conn.type_map_for_results = map
    conn.type_map_for_queries = PG::BasicTypeMapForQueries.new(conn)

    conn
  end

  private def start_reaper_thread
    Thread.new do
      Thread.current.name = "db-reaper"
      loop do
        @pool.reap(300) { |conn| conn&.close }
        sleep 60
      end
    end
  end
end
```

Each connection is wrapped with instrumentation for observability:

```ruby
# lib/instrumented_connection.rb
require "forwardable"
require "pg"
require "sentry-ruby"

class InstrumentedConnection
  extend Forwardable

  def_delegators :@conn, *(
    PG::Connection.public_instance_methods(false) - [
      :exec,
      :exec_params,
      :transaction
    ]
  )

  def initialize(conn, env)
    @conn = conn
    @env = env
  end

  def exec(sql, params = [])
    rows = []

    with_sentry(sql, params) do
      result = execute_query(sql, params)
      rows = result.to_a
    end

    rows
  rescue PG::ConnectionBad
    sleep 5 # give HA backup time to come online
    retry
  rescue => err
    if ["development", "test"].include?(@env)
      raise err
    else
      Sentry.capture_exception(err)
      []
    end
  end

  def transaction
    @conn.transaction { yield(self) }
  rescue PG::ConnectionBad
    sleep 5
    retry
  end

  private def execute_query(sql, params)
    if params.empty?
      @conn.exec(sql)
    else
      params = params.map { |p| p.is_a?(Hash) ? p.to_json : p }
      @conn.exec_params(sql, params)
    end
  end

  private def with_sentry(sql, params = [])
    tx = Sentry.get_current_scope&.get_span ||
      Sentry.get_current_scope&.get_transaction ||
      Sentry.start_transaction(name: "DB#exec")

    if tx
      tx.with_child_span(op: "db.sql.execute", description: sql) do |span|
        span.set_data("SQL Params", params)
        yield
      end
    else
      yield
    end
  end
end
```

## Usage

In controllers, access the pool:

```ruby
class ApplicationController < ActionController::Base
  private def db
    DB.pool
  end
end

class SuggestionsController < ApplicationController
  def new
    render json: Search::SuggestCompany.new(db).call(
      query: params[:query]
    )
  end
end
```

Pass `db` through initializers and use `<<~SQL` heredocs:

```ruby
module Search
  class SuggestCompany
    def initialize(db)
      @db = db
    end

    def call(query:)
      @db.exec(<<~SQL, [@db.fuzzy_like(query)])
        SELECT
          companies.id,
          companies.name,
          companies.status
        FROM
          companies
        WHERE
          companies.name ILIKE $1
          OR companies.also_known_as ILIKE $1
        ORDER BY
          companies.score DESC
        LIMIT 50
      SQL
    end
  end
end
```

For scripts, instantiate directly:

```ruby
if $0 == __FILE__
  require_relative "../db"
  pp Search::SuggestCompany.new(DB.new).call(query: "Data")
end
```

## Background processes

Initialize DB connections in background processes:

```ruby
# Each forked worker gets its own connection
children = workers.map do |worker|
  fork { worker.new(DB.new).poll }
end
```

See [ruby / job queues](/ruby/job-queues) and [ruby / clock](/ruby/clock)
for complete examples.

## Benefits

- Direct SQL with parameterized queries
- Connection pooling with automatic reaping
- Automatic retries for connection failures
- Sentry integration for error tracking and performance monitoring
- JSON serialization for hash parameters
- No DSL to learn or maintain
- Full control over every query
