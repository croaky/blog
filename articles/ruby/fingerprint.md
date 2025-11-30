# ruby / fingerprint

I use file-based asset fingerprinting in Ruby web apps
to enable aggressive caching with [CDNs](/web/cdn).

## The approach

As part of the deployment build,
after [building assets with esbuild](/cmd/esbuild),
a Rake task fingerprints the files:

```ruby
require "digest"

namespace :assets do
  task :precompile do
    ["public/css/app.css", "public/js/app.js"].each do |old_path|
      hash = Digest::MD5.file(File.expand_path(old_path, __dir__))
      ext = File.extname(old_path)
      base = old_path.chomp(ext)
      new_path = "#{base}-#{hash}#{ext}"
      system "mv #{old_path} #{new_path}"
    end
  end
end
```

Rails will serve the renamed files:

```
public/app.css  -> public/app-a1b2c3d4.css
public/app.js   -> public/app-a1b2c3d5.js
```

## Rails configuration

In `config/environments/production.rb`:

```ruby
config.public_file_server.enabled = true
config.public_file_server.headers = {
  "Cache-Control" => "public, max-age=31536000, immutable"
}
```

Since filenames include content hashes,
each URL is immutable.
Browsers and CDNs can cache aggressively (1 year)
without risk of serving stale content.

The `immutable` directive eliminates revalidation requests
even on page reload.

## Deployment

Example build command for [Render](https://render.com):

```bash
npm install && \
npm run build && \
bundle install && \
bundle exec rake db:migrate && \
bundle exec rake assets:precompile
```

This:

1. Installs JavaScript dependencies
2. Builds and bundles with [esbuild](/cmd/esbuild)
3. Installs Ruby dependencies
4. Migrates the database
5. Fingerprints static assets

## Template integration

In `config/initializers/assets.rb`:

```ruby
# see rake assets:precompile definition in Rakefile
# and app/views/layouts/application.haml

app_css_path = "/css/app.css"
app_js_path = "/js/app.js"

if ["staging", "production"].include?(ENV.fetch("APP_ENV"))
  path = Dir.glob("#{Rails.root}/public/css/app*.css")&.first
  if path
    app_css_path = path.split("public")[1]
  end

  path = Dir.glob("#{Rails.root}/public/js/app*.js")&.first
  if path
    app_js_path = path.split("public")[1]
  end
end

APP_CSS_PATH = app_css_path.freeze
APP_JS_PATH = app_js_path.freeze
```

In views such as `app/views/layouts/application.haml`:

```haml
%link{ rel: "stylesheet", href: APP_CSS_PATH }
%script{ src: APP_JS_PATH }
```
