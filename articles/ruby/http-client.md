# ruby / http client

I would rather write and maintain code like this
than depend on a library specific to the API:

```ruby
require "bundler/inline"

gemfile do
  source "https://rubygems.org"
  gem "http"
end

url = ENV["API_URL"]
if url.to_s.strip == ""
  puts "err: API_URL environment variable is not set"
  exit 1
end

begin
  resp = HTTP.post(url, json: { text: "hi" })
rescue => e
  puts "err: #{e.message}"
  exit 1
end

puts resp.body.to_s
```

SDK-style libraries have costs.
They may need to be upgraded to patch security issues
or resolve competing requirements in the dependency graph.
They can be slow to be updated or become unmaintained.
They are an additional interface for the team to learn.

In the above example,
I used the [HTTP](https://github.com/httprb/http) gem.
I prefer it to the Ruby standard library `net/http`'s interfaces
but ideally, I'd use the standard library,
such as this Go version of the same program:

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
)

func main() {
    url := os.Getenv("API_URL")
    if url == "" {
        log.Fatalln("err: API_URL environment variable is not set")
    }

    reqBody, err := json.Marshal(map[string]string{
        "text": "Hello, world!",
    })
    if err != nil {
        log.Fatalf("err: %v\n", err)
    }

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        log.Fatalf("err: %v\n", err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("err: %v\n", err)
    }

    fmt.Println(string(respBody))
}
```

As I use an API, I will build up my custom client as needed.
For example, I might add retries with exponential backoff.
