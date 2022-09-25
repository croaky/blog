# SSO with WorkOS in Ruby

[WorkOS lets you adds Single Sign-On (SSO)](https://workos.com/docs/sso/guide)
for all popular Identity Providers (e.g. Okta, SAML, Azure, Duo, LastPass)
with a single integration.
I've been happily using it in production for ~2 years.

## Template

Install gems:

```ruby
gem "rails"
gem "sentry"
gem "workos"
```

Store the WorkOS API key in your app's environment
and set it for the Ruby SDK:

```embed
code/workos/all.rb init
```

Set up HTTP endpoints:

```embed
code/workos/all.rb routes
```

Set up controller:

```embed
code/workos/all.rb controller
```
