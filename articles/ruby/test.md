# ruby / test

I wrote a custom test framework for Ruby.
I wanted:

- Fast: startup and runtime
- Simple: one assertion method, no DSL
- Features: db transactions, stubs, factories
- Flexible: plain Ruby, easy to extend
- Debuggable: small codebase, easy to understand

It has ~250 lines of code that
are included at the end of this article.

## Usage

Tests inherit from a `Test` base class:

```ruby
class DBTest < Test
  def test_exec_special_chars
    [
      "Company 100%",
      "O'Reilly Media",
      "Price: $100; Drop: 50%",
      "Robert'); DROP TABLE students;--",
      "Test\\Backslash",
      "Test'Quote\"DoubleQuote"
    ].each do |val|
      insert_company(name: val)

      rows = db.exec(<<~SQL, [val])
        SELECT
          name
        FROM
          companies
        WHERE
          name = $1
      SQL

      ok rows[0]["name"] == val
    end
  end

  def test_fuzzy_like_pattern
    got = db.fuzzy_like("test")

    ok got == "%test%"
  end
end
```

The `ok` method is the only assertion.
Pass a boolean expression.
If it's false, the test fails.
Add an optional message for context:

```ruby
ok user.valid?, "user is not valid"
ok rows.size == 3, "expected 3 rows, got #{rows.size}"
```

## Running tests

Run a test file directly:

```bash
ruby test/lib/db_test.rb
```

The framework randomizes test order and prints a seed:

```txt
seed 1234

DBTest
  test_exec_special_chars
  test_fuzzy_like_pattern

ok
```

Re-run with the same order using the seed:

```bash
ruby test/lib/db_test.rb --seed 1234
```

Run a single test method:

```bash
ruby test/lib/db_test.rb --name test_fuzzy_like_pattern
```

I trigger a single test method from Vim with
[a vim-test runner](https://github.com/croaky/laptop/commit/eb16cc13f6aaaf91436c5d3c97de50758b68e2de).

Run all tests:

```bash
ruby test/suite.rb
```

## Database transactions

Each test runs in a transaction that rolls back
to isolate each test:

```ruby
class TransactionTest < Test
  def test_insert
    co = insert_company(name: "Acme Inc")

    rows = db.exec("SELECT * FROM companies")
    ok rows.size == 1
  end

  def test_another_insert
    # Database is clean (previous test rolled back)
    rows = db.exec("SELECT * FROM companies")
    ok rows == []
  end
end
```

To test transaction behavior itself, set `@tx = false`:

```ruby
class TransactionBehaviorTest < Test
  def initialize
    @tx = false
  end

  def test_rollback
    # Test actual transaction behavior
    # Changes cleaned up with DELETE after test
  end
end
```

## Factories

I use factory methods for test data.
They work with [DB](/ruby/db):

```ruby
class CompaniesTest < Test
  def test_create
    co = insert_company(name: "Acme Inc", status: "Active")

    ok co.name == "Acme Inc"
    ok co.status == "Active"
  end

  def test_with_relationships
    co = insert_company
    per = insert_person(name: "Jane Doe")
    pos = insert_position(
      person_id: per.id,
      company_id: co.id,
      company_name: co.name,
      title: "CTO"
    )

    ok pos.person_id == per.id
  end
end
```

Factories provide defaults and return `Data` objects with attribute accessors.

## Stubs

The framework has built-in stubs for isolating collaborators.
All stubs are "spies" whose method calls can be verified.
Use them via dependency injection:

```ruby
module Companies
  class Import
    def initialize(db, client:)
      @db = db
      @client = client
    end

    def call(domain:)
      data, err = @client.fetch(domain)
      if err
        return "err: #{err}"
      end

      @db.exec(<<~SQL, [data["name"]])
        INSERT INTO companies (name)
        VALUES ($1)
      SQL

      "ok"
    end
  end
end

class CompaniesImportTest < Test
  def test_import
    client = stub(fetch: [{"name" => "Acme Inc"}, nil])

    status = Companies::Import.new(db, client: client).call(
      domain: "acme.com"
    )

    ok status == "ok"
    ok client.called?(:fetch)

    row = db.exec("SELECT * FROM companies").first
    ok row["name"] == "Acme Inc"
  end

  def test_api_error
    client = stub(fetch: [nil, "API rate limited"])

    got = Companies::Import.new(db, client: client).call(
      domain: "acme.com"
    )

    ok got == "err: API rate limited"
    ok db.exec("SELECT * FROM companies") == []
  end
end
```

Stubs support lambdas:

```ruby
client = stub(
  transform: ->(text) { text.upcase },
  calculate: ->(a, b) { a + b }
)
ok client.transform("hello") == "HELLO"
ok client.calculate(2, 3) == 5
```

For class methods, use `stub_class`.
Stubs are automatically restored after each test:

```ruby
class TimeTest < Test
  def test_frozen_time
    stub_class(Time, now: Time.at(0))

    ok Time.now == Time.at(0)
  end
end
```

## Implementation

The complete `test/test_helper.rb`:

```ruby
ENV["APP_ENV"] = "test"

require "webmock"

WebMock.enable!
WebMock.disable_net_connect!(allow_localhost: true)

require_relative "../lib/db"
require_relative "factories"

DB.configure do |c|
  c.pool_size = 1
  c.reap = false
end

class Test
  class Failure < StandardError; end

  include Factories
  include WebMock::API

  GREEN = "\e[32m"
  RED = "\e[31m"
  RESET = "\e[0m"

  @@groups = []
  @@seed = nil
  @@name = nil

  i = 0
  while i < ARGV.length
    case ARGV[i]
    when "--seed"
      @@seed = ARGV[i + 1].to_i if i + 1 < ARGV.length
      i += 2
    when "--name"
      @@name = ARGV[i + 1] if i + 1 < ARGV.length
      i += 2
    else
      i += 1
    end
  end

  def self.inherited(c)
    @@groups << c
  end

  def self.run_suite
    seed = @@seed || rand(1000..9999)
    srand seed
    puts "seed #{seed}\n"

    @@groups.shuffle.each do |c|
      c.run_group
    end

    print "\n#{GREEN}ok#{RESET}\n"
  end

  def self.run_group
    group = new

    tests = public_instance_methods(false)
      .grep(/^test_/)
      .shuffle

    if @@name
      tests = tests.select { |t| t.to_s == @@name }
      if tests == []
        return
      end
    end

    if tests == []
      return
    end

    puts "\n#{self}"
    tests.each { |test| group.run_test(test) }
  end

  def db
    DB.pool
  end

  def initialize
    @tx = true
    @stubs = []
  end

  def run_test(test)
    setup
    send(test)
    puts "  #{GREEN}#{test}#{RESET}"
  rescue => err
    puts "  #{RED}#{test}#{RESET}"
    lines = err.backtrace.reject { |l| l.include?(__FILE__) }.join("\n  ")
    puts "\n#{RED}fail: #{err}#{RESET}\n  #{lines}"
    exit 1
  ensure
    teardown
  end

  def ok(expression, m = nil)
    if !expression
      raise Test::Failure, m
    end
  end

  def stub(methods)
    obj = Object.new
    calls = Hash.new { |h, k| h[k] = [] }

    methods.each do |meth, return_value|
      obj.define_singleton_method(meth) do |*args, **kwargs, &block|
        calls[meth] << {args: args, kwargs: kwargs}
        if return_value.is_a?(Proc)
          return_value.call(*args, **kwargs, &block)
        else
          return_value
        end
      end
    end

    obj.define_singleton_method(:called?) do |meth|
      calls[meth] != []
    end

    obj.define_singleton_method(:calls) do
      calls
    end

    obj
  end

  def stub_class(klass, methods)
    methods.each do |meth, return_value|
      orig = klass.method(meth)
      @stubs << [klass, meth, orig]

      klass.define_singleton_method(meth) do |*args, **kwargs, &block|
        if return_value.is_a?(Proc)
          return_value.call(*args, **kwargs, &block)
        else
          return_value
        end
      end
    end
  end

  private def setup
    if @tx
      db.exec("BEGIN")
    end
  end

  private def teardown
    @stubs.reverse.each do |klass, meth, orig|
      klass.define_singleton_method(meth, orig)
    end
    @stubs = []

    if @tx
      db.exec("ROLLBACK")
    else
      tablenames = db.exec(<<~SQL).map { |row| row["tablename"] }
        SELECT
          tablename
        FROM
          pg_tables
        WHERE
          schemaname = 'public'
          AND tablename != 'users'
        ORDER BY
          tablename
      SQL

      tablenames.each do |t|
        db.exec("DELETE FROM #{t}")
      end

      db.exec("DELETE FROM users WHERE id != 1")
    end
  end
end

at_exit { Test.run_suite }
```

The `factories.rb` file defines factory methods:

```ruby
module Factories
  class Sequence
    def initialize
      @counter = 0
    end

    def next
      @counter += 1
    end
  end

  SEQ = Sequence.new

  def insert_company(o = {})
    insert_into("companies", {
      name: o[:name] || "Company #{SEQ.next}"
    }.merge(o))
  end

  def insert_person(o = {})
    insert_into("people", {
      name: o[:name] || "Person #{SEQ.next}"
    }.merge(o))
  end

  def insert_position(o = {})
    insert_into("positions", {
      company_id: o[:company_id] || insert_company.id,
      person_id: o[:person_id] || insert_person.id,
      title: o[:title] || "CEO, Founder"
    }.merge(o))
  end

  private def insert_into(table, attrs)
    row = db.exec(<<~SQL, attrs.values).first
      INSERT INTO #{table} (
        #{attrs.keys.join(", ")}
      ) VALUES (
        #{(1..attrs.size).map { |i| "$#{i}" }.join(", ")}
      )
      RETURNING *
    SQL

    Data.define(*row.keys.map(&:to_sym)).new(*row.values)
  end
end
```

The `test/suite.rb` file requires all tests
and the `at_exit` hook in `test_helper.rb` runs the suite:

```ruby
require_relative "test_helper"

Dir["#{__dir__}/**/*_test.rb"].each { |f| require f }
```
