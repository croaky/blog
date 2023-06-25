# begindoc: all
require "bundler/inline"

gemfile do
  source "https://rubygems.org"
  gem "http"
end

HTTP.post(ENV["SLACK_WEBHOOK"], json: {
  text: "Hello, world!"
})
# enddoc: all
