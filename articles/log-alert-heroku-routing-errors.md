# Log and Alert on Heroku Routing Errors

Most application errors can be sent to a third-party error service
such as [Sentry](https://sentry.io/).

On Heroku, routing errors never reach these services
because they don't occur in the application's processes, they occur in
[Heroku's routing layer](https://devcenter.heroku.com/articles/http-routing).

To observe and alert on these errors,
send them to a logging service such as
[Papertrail](https://devcenter.heroku.com/articles/papertrail).

## Heroku Routing Errors

You may have seen
[Heroku Platform Error Codes][codes] in your logs:

* `H12 - Request Timeout`
* `H15 - Idle connection`
* `H18 - Request Interrupted`

These can be caused by many different factors ranging from
misconfigured web server concurrency
to slow clients (mobile phones on weak cell connections).

Before we can tune our app,
we need to first know these errors are occurring.

The bad news is we can't use our usual error tracking systems.
The good news is that Heroku
reliably includes the text `status=503` in the logs for these errors.

[codes]: https://devcenter.heroku.com/articles/error-codes

## Papertrail

Papertrail is a logger with a Heroku add-on.
Add and open it:

```bash
heroku addons:add papertrail
heroku addons:open papertrail
```

* Search for `status=503` and click "Save Search".
* Give it a name and click "Save & Setup an Alert".
* Click "Slack".
* Command-click "new Papertrail integration" to open Slack in a new browser tab.
* Select your [Slack] channel for the Papertrail integration.
* Click "Copy URL" to get the webhook URL.
* Paste it back in the Papertrail add-on.
* Click "Create Alert".

[Slack]: https://slack.com

## Slack

Going forward,
you'll receive alerts for any Heroku routing errors
in your project's Slack channel.

Each log entry has a `request_id` that you can copy
and paste into Papertrail to see the contextual requests
before and after the 503.
