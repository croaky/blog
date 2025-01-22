# Thin API clients

I've sometimes seen open source libraries
used as an API client for a SaaS service
when a better choice would be
a "thin" client written by the application developers.

Here's an example of a "thin" API client:

```ruby
require "bundler/inline"

gemfile do
  source "https://rubygems.org"
  gem "http"
end

HTTP.post(ENV["SLACK_WEBHOOK"], json: {
  text: "Hello, world!"
})
```

This triggers a
[Slack incoming webhook](https://slack.com/apps/A0F7XDUAZ-incoming-webhooks)
from Ruby.

I would rather write and maintain this code
with a generic HTTP dependency
than depend on a Slack-specific third-party library.

Third-party libraries have costs.
They need to be managed by a package manager.
They may need to be upgraded to patch security issues
or resolve competing requirements in the dependency graph.
They can become unmaintained.
They are an additional interface for the team to learn.

In the Ruby example,
I used a third-party [HTTP](https://github.com/httprb/http) library
because I prefer its interface to the Ruby standard library's interface.
A "thinner" client would use only the language's standard library,
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
	url := os.Getenv("SLACK_WEBHOOK")
	if url == "" {
		log.Fatalln("no webhook provided")
	}

	reqBody, err := json.Marshal(map[string]string{
		"text": "Hello, world!",
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(respBody))
}
```

If my only needs are a `POST` request with JSON body,
I would write a thin client.
If there is some lightweight authentication with a header token,
an idempotency key, or retry logic,
I would still choose to write a thin client.
