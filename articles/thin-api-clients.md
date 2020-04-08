# Thin API Clients

I've sometimes seen open source libraries
used as a special-purpose API client for a SaaS service
when a "thin" client written by the application developers
would be a better choice.

Here's an example of a "thin" API client:

```embed
code/slack_webhook.rb all
```

This triggers a
[Slack incoming webhook](https://slack.com/apps/A0F7XDUAZ-incoming-webhooks)
from Ruby.

I would rather write and maintain this code,
which uses a generic [HTTP](https://github.com/httprb/http) dependency,
than depend on a third-party library or package,
such as a specific Slack client gem.

Third-party packages have costs.
They need to be managed by a package manager.
They may need to be upgraded in the future
to patch security issues
or to resolve competing requirements in the dependency graph.
They can become unmaintained.
They are an additional interface for the team to learn.

In the Ruby example, I used a third-party HTTP package
because I like its interface better than the Ruby standard library's interface.
A "thinner" client would use only the language's standard library.
Here's a Go version of the same program:

```embed
code/slackwebhook.go all
```

If all I'm dealing with is a `POST` request,
JSON body, and some lightweight authentication such as a header token,
I would choose the thin client approach.

If there are some lightweight idempotency key and retry logic,
I may still choose to write a thin client.

I would use a third-party package if the authentication mechanism
is more complex or we need to use a large surface area of the API
(many endpoints that may be stitched together).

A middle ground is choosing the lightest weight option amongst
multiple open source packages. For example, when I've written a feature to
add a new credit card to a React Native app, I've used the
[stripe-client](https://www.npmjs.com/package/stripe-client) Node package
instead of alternatives that pull in iOS and Android dependencies.
When all we need to do is get a token from Stripe's API,
then send the token to our backend for processing,
the heavier weight option is not worth the costs.

```embed
code/stripe-card-token.ts all
```
