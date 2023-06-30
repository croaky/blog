# Four-phase test

Four-phase test is a testing pattern, applicable to all programming languages
and unit tests (but not integration tests).

It takes the following general form:

```ruby
test do
  setup
  exercise
  verify
  teardown
end
```

There are four distinct phases of the test,
executed sequentially.

Set up system under test (usually a class, object, or method):

```ruby
user = User.new(password: "password")
```

Exercise the system under test:

```ruby
user.save
```

Verify the result of the exercise against the developer's expectations:

```ruby
expect(user.encrypted_password).to_not be_nil
```

Tear down the system under test to its pre-setup state.
This is usually handled implicitly by the language (releasing memory)
or test framework (running inside a database transaction).

The four phases are wrapped into a named test.

A related style guideline is to
separate setup, exercise, verification, and teardown phases with newlines:

```ruby
describe "#save" do
  it "encrypts the password" do
    user = User.new(password: "password")

    user.save

    expect(user.encrypted_password).to_not be_nil
  end
end
```
