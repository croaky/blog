# ruby / test

I use a custom test framework for Ruby. My goals are:

- Fast: startup and runtime
- Simple: one assertion method
- Flexible: plain Ruby, easy to extend
- Debuggable: small codebase, easy to understand
- Features: stubs, db transactions, db factories, etc.

It has ~250 lines of code that
are included at the end of this article.

## Test groups and test cases

Test groups inherit from a `Test` base class:

```ruby
class MathTest < Test
  def test_greater_than
    ok 10 > 5
  end

  def test_less_than
    ok 3 < 7
  end
end
```

Test cases are public instance methods
whose names start with `test_`.

One ore more test groups can be defined
in the same file.

## Assertions

The `ok` method is the only assertion.
It takes a boolean expression.

Add an optional message for context:

```ruby
class NilTest < Test
  def test_nil
    val = nil
    ok val == nil, "#{val} not nil"
  end

  def test_not_nil
    val = "value"
    ok val != nil, "val is nil"
  end
end

class RegexTest < Test
  def test_match
    got = "user@example.com"
    want = /\A[\w+\-.]+@[a-z\d\-.]+\.[a-z]+\z/i
    ok got =~ want, "#{got} not email format"
  end

  def test_no_match
    got = "text"
    ok got !~ /[<>]/, "#{got} contains HTML brackets"
  end
end

class ExceptionTest < Test
  def test_raised
    raised = false

    begin
      raise ArgumentError, "invalid argument"
    rescue ArgumentError
      raised = true
    end

    ok raised, "did not raise ArgumentError"
  end

  def test_not_raised
    raised = false
    err = nil

    begin
      10 / 2
    rescue => e
      err = e
      raised = true
    end

    ok !raised, "raised #{err}"
  end
end
```

If the expression passed to `ok` is true,
the assertion passes and the test continues.

If the expression is false, the assertion fails
and the test runner prints a backtrace
and exits immediately with a non-zero status code.

## Runner

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

The `test_` outputs are green for passing tests
and red for failing tests.
If `ENV["CI"]` is set, no color codes are used.

Re-run with the same order using the seed:

```bash
ruby test/lib/db_test.rb --seed 1234
```

Run a single test case from the command line:

```bash
ruby test/lib/db_test.rb --name test_fuzzy_like_pattern
```

Or, run a single test case from Vim with
[a vim-test runner](https://github.com/croaky/laptop/commit/eb16cc13f6aaaf91436c5d3c97de50758b68e2de).

To run multiple test files, create a file
that requires them all, e.g. `test/suite.rb`:

```ruby
require_relative "test_helper"

Dir["#{__dir__}/**/*_test.rb"].each { |f| require f }
```

Then run it:

```bash
ruby test/suite.rb
```

## Before suite

Before any tests run, the framework can execute
one-time setup code directly in `test_helper.rb`.
This runs once when the file is loaded:

```ruby
# before suite
DB.pool.exec(<<~SQL)
  INSERT INTO users (id, name, admin, email)
  VALUES (1, 'Admin', true, 'admin@example.com')
  ON CONFLICT DO NOTHING;

  ALTER SEQUENCE users_id_seq RESTART WITH 2;

  REFRESH MATERIALIZED VIEW cache_companies;
SQL
```

This is useful for:

- Creating shared fixtures
- Populating materialized views
- Other one-time expensive setup

## Database transactions

Each test case runs in a transaction that rolls back
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
    super
    @tx = false
  end

  def test_rollback
    # Test actual transaction behavior
    # Changes cleaned up with DELETE after test
  end
end
```

## Database factories

Factory methods for test data work with [DB](/ruby/db):

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

Factories provide defaults and return `Data` objects
with attribute accessors.

## State-based

Prefer state-based assertions whenever possible.
Assert results or side effects in the database:

```ruby
def test_create_company
  Companies::Create.new(db).call(name: "Acme Inc")

  row = db.exec("SELECT * FROM companies").first
  ok row["name"] == "Acme Inc"
end
```

## Object stubs

When state-based testing isn't practical,
object stubs can help isolate collaborators.
Use `stub` via dependency injection:

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

Stubs support lambdas for transformations:

```ruby
client = stub(
  transform: ->(text) { text.upcase },
  calculate: ->(a, b) { a + b }
)
ok client.transform("hello") == "HELLO"
ok client.calculate(2, 3) == 5
```

## Class method stubs

For class methods, use `stub_class`:

```ruby
class TimeTest < Test
  def test_frozen_time
    stub_class(Time, now: Time.at(0))

    ok Time.now == Time.at(0)
  end
end
```

Class method stubs are automatically restored after each test.

Class method stubs also support lambdas:

```ruby
# identity functions
stub_class(Convert::ExtractDomain, call: ->(host) { host })

# raise errors
stub_class(Aws::S3::Client, new: ->(*) {
  raise StandardError.new("auth failed")
})

# capture variables
err_msg = nil
stub_class(Sentry, capture_exception: ->(e) { err_msg = e.message })
some_code_that_raises
ok err_msg == "expected error"
```

## Asserting stub calls

All stubs are "spies" whose method calls can be asserted
with `called?`:

```ruby
client = stub(fetch: [{"name" => "Acme Inc"}, nil])

Companies::Import.new(db, client: client).call(domain: "acme.com")

ok client.called?(:fetch)
ok !client.called?(:delete)
```

Or, assert call count and arguments with `calls`:

```ruby
ok client.calls[:fetch].size == 2
ok client.calls[:fetch][0][:args] == ["acme.com"]
ok client.calls[:fetch][0][:kwargs] == {domain: "acme.com"}
```

`calls` returns a hash mapping method names
to arrays of call records.
Each call record is a hash with `:args` and `:kwargs` keys.

## Yielding stubs

For methods that yield, instead of `stub`,
use `Object.new` with `def`:

```ruby
client = Object.new
def client.get_data(_)
  yield "chunk1"
  yield "chunk2"
end

got = []
client.get_data("http://example.com") { |chunk| got << chunk }
ok got == ["chunk1", "chunk2"]
```

For `Object.new` stubs,
capture values with instance variables:

```ruby
client = Object.new
def client.process(_)
  @thread_ref = Thread.current
  yield "data"
end
def client.thread_ref
  @thread_ref
end

some_code_under_test(client)

ok !client.thread_ref.alive?
```

## Style guide

Prefer inlining code and avoiding unnecessary local variables.

When they clarify tests or improve failure messages,
use `got` and `want` variable names:

```ruby
def test_length
  got = "hello".length
  want = 5
  ok got == want, "#{got} != #{want}"
end
```

Typically, separate setup, exercise, and assertion phases with blank lines:

```ruby
def test_add
  a = 2
  b = 3

  got = a + b

  ok got == 5
end
```

When exercising the system under test multiple times,
group exercise and assertion together:

```ruby
def test_multiply
  got = 2 * 3
  ok got == 6

  got = 4 * 5
  ok got == 20

  got = 0 * 10
  ok got == 0
end
```

To reduce verbosity, name fresh fixtures with abbreviations (`co`, `per`, `u`).
When querying a changed fixture from the database,
prefix the variable with `db_` to distinguish it from the original:

```ruby
co = insert_company(name: "foo")

Something.new(db).call("bar")

db_co = db.exec("SELECT * FROM companies WHERE id = $1", [co.id]).first
ok db_co["name"] == "bar"
```

For SQL in assertions, use one-line `SELECT *` for simple predicates:

```ruby
note = db.exec("SELECT * FROM notes WHERE company_id = $1", [co.id]).first
```

For more complex queries, use heredocs:

```ruby
job = db.exec(<<~SQL, [per.id]).first
  SELECT
    jobs.*
  FROM
    notes
    JOIN jobs ON jobs.args ->> 'note_id' = notes.id::text
  WHERE
    notes.person_id = $1
    AND jobs.queue = 'slack'
SQL
```

Assert empty tables with `ok rows == []`, not `size == 0`.

```ruby
rows = db.exec("SELECT * FROM tracking WHERE company_id = $1", [co.id])
ok rows == []
```

For unordered comparisons, map and sort:

```ruby
rows = db.exec("SELECT * FROM list_items WHERE company_id = $1", [co.id])
ok rows.map { |r| r["list_id"] }.sort == [748, 541].sort
```

For error messages, prefer `include?` over full-array equality:

```ruby
ok got[:errs].include?("Company is required")
```

## Rails controller testing

For Rails controller tests, extend `Test` with
`Rack::Test` methods and helpers:

```ruby
require_relative "test_helper"
require File.expand_path("../config/environment", __dir__)
require "rackup"

class ControllerTest < Test
  include Rack::Test::Methods

  def app
    Rails.application
  end

  def sign_in
    set_cookie("remember_token=test")
  end

  def sign_in_as(user)
    set_cookie("remember_token=#{user.remember_token}")
  end

  def cookies
    @cookies ||= Cookies.new(rack_mock_session.cookie_jar)
  end

  class Cookies
    def initialize(jar)
      @jar = jar
    end

    def [](name)
      cookie = @jar.get_cookie(name.to_s)
      cookie&.value
    end
  end

  def flash
    Flash.new(last_request)
  end

  class Flash
    def initialize(request)
      @request = request
    end

    def [](key)
      rack_session = @request.env["rack.session"]
      if rack_session.nil?
        return nil
      end

      flash_hash = rack_session.dig("flash", "flashes")
      if flash_hash.nil?
        return nil
      end

      flash_hash[key.to_s]
    end
  end

  # override Rack::Test methods to return last_response
  def get(path, params = {}, headers = {})
    super
    last_response
  end

  def post(path, params = {}, headers = {})
    super
    last_response
  end

  private def teardown
    clear_cookies
    header "Ajax-Referer", nil
    super
  end
end
```

Use `ControllerTest` for Rails controller tests:

```ruby
class CompaniesControllerTest < ControllerTest
  def test_index
    sign_in
    co = insert_company(name: "Acme Inc")

    resp = get("/companies")

    ok resp.status == 200
    ok resp.body.include?("Acme Inc")
  end

  def test_create
    sign_in

    resp = post("/companies", {company: {name: "New Co"}})

    ok resp.status == 302
    ok flash[:notice] == "Company created"
    ok cookies["remember_token"] == "test"
  end
end
```

Separate controller tests from other tests
with different suite files:

```ruby
# test/ruby_suite.rb
require_relative "test_helper"

Dir.glob(File.join(__dir__, "**", "*_test.rb"))
  .reject { |f| f.include?("/controllers/") }
  .sort
  .each { |f| require f }

# test/rails_suite.rb
require_relative "rails_helper"

Dir.glob(File.join(__dir__, "controllers", "**", "*_test.rb"))
  .sort
  .each { |f| require f }
```

Run them separately:

```bash
ruby test/ruby_suite.rb  # fast, no Rails
ruby test/rails_suite.rb # slower, loads Rails
```

An `at_exit` hook in `test/test_helper.rb`
automatically runs each suite.

## Implementation

The `test/test_helper.rb` file:

```ruby
ENV["APP_ENV"] = "test"

require "webmock"
require_relative "../lib/db"
require_relative "factories"

WebMock.enable!
WebMock.disable_net_connect!(allow_localhost: true)

DB.configure do |c|
  c.pool_size = 1
  c.reap = false
end

class Test
  class Failure < StandardError; end

  include Factories
  include WebMock::API

  if ENV["CI"]
    GREEN = ""
    RED = ""
    RESET = ""
  else
    GREEN = "\e[32m"
    RED = "\e[31m"
    RESET = "\e[0m"
  end

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

      # app-specific cleanup of all users except admin fixture
      db.exec("DELETE FROM users WHERE id != 1")
    end
  end
end

at_exit { Test.run_suite }
```

The `test/factories.rb` file:

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

There are other Ruby testing frameworks available,
but this one is optimized for my happiness.
