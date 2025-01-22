# ruby / file guard

I commonly use this stanza at the bottom of my Ruby files:

```ruby
if $0 == __FILE__
  # pp something
end
```

`$0` contains the file name of the Ruby script currently being run. If the
script is being executed directly, `$0` will hold the name of the script file.
If the script is being required or loaded as a module, `$0` will hold the name
of the main script that initiated the loading.

`__FILE__` is a built-in constant in Ruby that holds the current file's path. It
represents the path to the file that contains the current line of code.

Here's a longer example:

```ruby
require "http"
require "json"

module GitHub
  class Client
    def get(path) # => (json: Hash | nil, err: String | nil)
      resp = HTTP.timeout(1).get("https://api.github.com#{path}")
      if resp.code / 100 != 2
        return [nil, resp.status]
      end

      [JSON.parse(resp.body), nil]
    rescue => err
      [nil, err.to_s]
    end
  end
end

if $0 == __FILE__
  pp GitHub::Client.new.get("/orgs/thoughtbot/repos")
end
```

I can reference `GitHub::Client` from other files in the project,
which will not run the code in the `if $0 == __FILE__` guard.

I can also run the file directly for a tight feedback loop during development.
