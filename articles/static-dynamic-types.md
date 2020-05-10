# Blurry Line Between Static and Dynamic Typing

From an [interview with Rob Pike](https://evrone.com/rob-pike-interview):

> Evrone: With "gradual typing" being introduced into "dynamically typed"
> languages and "type inference" into "statically typed", the line between two
> is now blurred. What is your opinion on a type system for a modern programming
> language?

> Rob: I am a big fan of static typing because of the stability and safety it
> brings.

> I am a big fan of dynamic typing because of the fun and lightweight feel it
> brings. (As a side note, the big push for integrated unit testing can be
> credited to languages like Python, which drove testing to demonstrate
> correctness that the typing system failed to provide.)

> I am not a fan of type-driven programming, type hierarchies and classes and
> inheritance. Although many hugely successful projects have been built that
> way, I feel the approach pushes important decisions too early into the design
> phase, before experience can influence it. In other words, I prefer
> composition to inheritance.

> However, I say to those who are comfortable using inheritance to structure
> their programs: pay no attention and please continue to use what works for
> you.

My taste preferences match this sentiment on almost every point.

Ruby can feel fun and lightweight:

```ruby
require "http"
require "pg"

db_conn_url = ENV.fetch("DATABASE_URL")

# publish
fork do
  sleep 1
  conn = PG.connect(db_conn_url)
  conn.exec "NOTIFY slackhttp, 'data'"
end

# subscribe
begin
  conn = PG.connect(db_conn_url)
  conn.exec "LISTEN slackhttp"

  loop do
    conn.wait_for_notify do |event, id, data|
      HTTP.post(ENV.fetch("SLACK_WEBHOOK_URL"), body: data)
    end
  end
ensure
  conn.exec "UNLISTEN slackhttp"
end
```

Or, it can feel confusing and heavyweight as object hierarchies
descend from `ActiveRecord`.

I know where everything is defined in Go and TypeScript...

```go

```
