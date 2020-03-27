# Intercept Email from Staging

Intercept email from a Rails app
and deliver it to a whitelist of email addresses
so the team can preview emails
without accidentally delivering staging email to production customers.

Integrate <https://github.com/croaky/recipient_interceptor>:

```ruby
# Gemfile
gem "recipient_interceptor"

# config/environments/production.rb
My::Application.configure do
  config.action_mailer.default_url_options = { host: ENV.fetch("HOST") }
  config.action_mailer.delivery_method = :smtp
  config.action_mailer.smtp_settings = {
    address: "smtp.sendgrid.net",
    authentication: :plain,
    domain: "heroku.com",
    password: ENV.fetch("SENDGRID_PASSWORD"),
    port: "587",
    user_name: ENV.fetch("SENDGRID_USERNAME")
  }
end

# config/environments/staging.rb
require_relative "production"
Mail.register_interceptor(
  RecipientInterceptor.new(ENV.fetch("EMAIL_RECIPIENTS"))
)
```

Use the `EMAIL_RECIPIENTS` environment variable
to update the comma-separated list of email addresses
that should receive staging emails.
For example:

```
heroku config:add EMAIL_RECIPIENTS="staging@example.com"
```
