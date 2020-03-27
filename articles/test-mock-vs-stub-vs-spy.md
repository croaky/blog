# Test Mock vs. Stub vs. Spy

[Test mocks](http://xunitpatterns.com/Mock%20Object.html),
[test stubs](http://xunitpatterns.com/Test%20Stub.html), and
[test spies](http://xunitpatterns.com/Test%20Spy.html)
can be used in unit tests, with tradeoffs.
Consider this system under test:

```ruby
def fetch(url)
  HTTP.get(url).body
rescue => err
  ErrorLogger.log(err, { url: url })
  nil
end
```

The method interacts with two collaborators:
`HTTP` for HTTP requests
and `ErrorLogger` for logging errors.

## Mock

A unit test using a mock:

```ruby
describe "fetch" do
  it "logs errors" do
    allow(HTTP).to receive(:get).and_raise("error")
    expect(ErrorLogger).to receive(:log) # mock

    result = fetch("https://example.com")

    expect(result).to be_nil
  end
end
```

Nothing is duplicated but
the phases of the test are "setup, verify, exercise, verify",
which can be confusing to read.

## Stub

A unit test using a stub:

```ruby
describe "fetch" do
  it "logs errors" do
    allow(HTTP).to receive(:get).and_raise("error")
    allow(ErrorLogger).to receive(:log) # stub

    result = fetch("https://example.com")

    expect(ErrorLogger).to have_received(:log)
    expect(result).to be_nil
  end
end
```

[`allow`](https://github.com/rspec/rspec-mocks#method-stubs) stubs
the collaborator.
[`expect`](https://github.com/rspec/rspec-mocks#test-spies)
asserts an expectation was met on the stub.

This style keeps a [Four-Phase Test](/four-phase-test) order,
emphasized by newlines separating
setup, exercise, and verification phases.

## Spy

A unit test using a
[spy](https://relishapp.com/rspec/rspec-mocks/docs/basics/spies):

```ruby
describe "fetch" do
  it "logs errors" do
    allow(HTTP).to receive(:get).and_raise("error")
    spy(ErrorLogger) # stub

    result = fetch("https://example.com")

    expect(ErrorLogger).to have_received(:log)
    expect(result).to be_nil
  end
end
```

The `spy` is more informative about the test double's purpose
and it removes duplicated references to `log`.
