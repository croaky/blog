# Process Model

Define Unix processes in a manifest named `Procfile`
and use tools to manage those processes.

Rails app:

```
web: ./bin/rails server
webpacker: ./bin/webpack-dev-server
worker: QUEUES=mailers ./bin/rake jobs:work
```

Rails API with React frontend:

```
client: cd client && npm start
server: cd server && bundle exec puma -C config/puma.rb
```

Sinatra app:

```
web: cd canary && bundle exec ruby web.rb
```

Go API with React frontend:

```
client: cd client && npm start
server: cd serverd && go install && serverd
```

In development,
tools like [Foreman](http://ddollar.github.io/foreman/)
interleave output streams,
respond to crashed processes,
and handle user-initiated restarts and shutdowns.

```
foreman start
```

In production,
[Heroku automatically uses the `Procfile`][Heroku] to specify the app's dynos.
Foreman can [export] the `Procfile`'s process definitions
to other formats such as `systemd`:

[Heroku]: https://devcenter.heroku.com/articles/procfile
[export]: https://ddollar.github.io/foreman/#EXPORTING

```
foreman export systemd .
```
