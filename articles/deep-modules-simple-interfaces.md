# Deep modules, simple interfaces

One of my favorite books on software is
[A Philosophy of Software Design](https://amzn.to/2OQkBEQ),
especially its parts about "deep modules."

Deep modules have powerful functionality and simple interfaces:

```txt
===========
|         |
|         |
|         |
|         |   =============================
|         |   |                           |
|         |   |                           |
-----------   -----------------------------
Deep module   Shallow module
```

The breadth `===` of interface is a cost.
The depth `|` of functionality is a benefit.

I believe this principle to be true at the micro-level of programming such as
defining structs, classes, or packages.

For example, the garbage collector in a language such as [Go](https://go.dev/)
is as deep module. This module has no interface at all; it works invisibly
behind the scenes to reclaim unused memory. Adding garbage collection to a
system actually shrinks its overall interface, since it eliminates the interface
for freeing objects.

The effect is magnified at the level of APIs.

Send an SMS over a global telephony network
with the [Twilio](https://www.twilio.com) API:

```ruby
require "twilio-ruby"

client = Twilio::REST::Client.new(
  ENV["TWILIO_ACCOUNT_SID"],
  ENV["TWILIO_AUTH_TOKEN"]
)

client.messages.create(
  from: "+15551234567",
  to: "+15555555555",
  body: "Hey friend!"
)
```

Charge a credit card over a global financial network
with the [Stripe](https://stripe.com/) API:

```ruby
require "stripe"

Stripe.api_key = ENV["STRIPE_API_KEY"]

Stripe::Charge.create(
  amount: 2000,
  currency: "usd",
)
```

Chat with a 1.76 trillion parameter large language model
with the [OpenAI](https://openai.com/) API:

```ruby
require "openai"

client = OpenAI::Client.new(access_token: ENV["OPENAI_ACCESS_TOKEN"])

response = client.chat(
  parameters: {
  model: "gpt-4",
  messages: [
    {
      role: "user",
      content: "What is the weather like in San Francisco?"
    }
  ],
  }
)
puts response.dig("choices", 0, "message", "content")
```

[Cloudflare](https://www.cloudflare.com/)'s interface is almost invisible:
manage a few DNS records to protect your web apps and databases from DDoS attacks.

These are the deepest modules I can think of. With only a few
lines of code or configuration, we can utilize hugely powerful systems.
However, they can fail by including unnecessary details or omitting crucial
ones.

Software is composed of layers, each providing a different abstraction. It's
important to pull complexity downwards, making modules simpler for users at the
expense of more complex implementations. This approach minimizes overall system
complexity, as most modules have more users than developers.
