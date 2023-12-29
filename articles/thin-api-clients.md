# Thin API clients

I've sometimes seen open source libraries
used as an API client for a SaaS service
when a better choice would be
a "thin" client written by the application developers.

Here's an example of a "thin" API client:

```embed
code/thin-api-clients/slack_webhook.rb all
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

```embed
code/thin-api-clients/slackwebhook.go all
```

If my only needs are a `POST` request with JSON body,
I would write a thin client.
If there is some lightweight authentication with a header token,
an idempotency key, or retry logic,
I would still choose to write a thin client.

Another "thin" variation is
to choose the lightest option of open source libraries.
For example, to add a new credit card to a React Native app, I've used the
[stripe-client](https://www.npmjs.com/package/stripe-client) Node package
instead of alternatives that have iOS and Android dependencies.
In this case, we only need to get a token from Stripe's API
and send the token to our backend for processing.

```embed
code/thin-api-clients/stripe-card-token.ts all
```
