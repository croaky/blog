# Heroku to Slack with AWS Lambda

When my production Heroku app processes change,
I want to be notified in Slack:

![Screenshot of Slack notification](images/heroku-to-slack.png)

Other examples:

```
17:38:52 clock.1 `bundle exec ruby schedule/clock.rb` up
17:39:03 web.1 `bundle exec puma -p $PORT -C ./config/puma.rb` up
17:39:05 web.2 `bundle exec puma -p $PORT -C ./config/puma.rb` up
```

Pager-notifying events:

```
17:38:52 queuenote.1 `bundle exec ruby queue/note.rb` crashed
```

[Heroku has webhooks](https://devcenter.heroku.com/articles/app-webhooks)
for these events but their payloads aren't in the format needed for
[Slack incoming webhooks](https://api.slack.com/messaging/webhooks).

AWS Lambda is the perfect glue to transform the Heroku webhook's JSON payload
into a useful JSON payload for Slack's incoming webhook.

## Slack config

Create an [incoming webhook](https://api.slack.com/messaging/webhooks).
Copy the URL.

## Lambda config

Create a [Lambda function](https://docs.aws.amazon.com/lambda/latest/dg/getting-started.html).
AWS' supported runtimes include Node, Python, Ruby, and Go.
You can alternatively implement a custom runtime.
Here's an example in Ruby:

```embed
code/heroku-to-slack.rb all
```

Paste the Slack incoming webhook URL as an environment variable,
which is encrypted at rest.

Create an [API Gateway](https://docs.aws.amazon.com/apigateway/latest/developerguide/welcome.html)
to make the Lambda function accessible in the Heroku web UI.

## Heroku config

Go to:

```
https://dashboard.heroku.com/apps/YOUR-APP-NAME/webhooks
```

Create a webhook with event type "dyno".
Paste the API Gateway URL as the Payload URL.

## Modify to taste

Edit and save the code in Lambda's web-based text editor.
Trigger a webhook to test the function.
View the auto-created CloudWatch logs for each function call.
