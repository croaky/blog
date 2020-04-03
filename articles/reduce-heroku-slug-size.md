# Reduce Heroku Slug Size

Have you ever seen this error when you tried to deploy to your Heroku app?

> Compiled slug size: 300.6M is too large (max is 300M)

You can't deploy but it can be fixed.
Here are ways to reduce the slug size.

## Clean up the Git repository

Install the [Heroku Repo plugin][repo]:

[repo]: https://github.com/heroku/heroku-repo

```
heroku plugins:install https://github.com/heroku/heroku-repo.git
```

Then run:

```
heroku repo:gc --app your-app-name
heroku repo:purge_cache --app your-app-name
```

These commands will execute `git gc --agressive` ([git-clean]) and
delete the contents of the Heroku build cache stored in
your application's Git repository.

[git-clean]: https://git-scm.com/docs/git-clean

Doing this on a recent app reduced my slug size by 100M.

## Move some files out of the repo

Move internal design documents (`.sketch` files) somewhere like Dropbox.
Move user-facing media (`.mp3`, `.mpg` files) somewhere like Amazon S3.

## Ignore some files that have to be in the repo

[Follow Heroku's instructions][slugignore] to
ignore files such as unit tests with `.slugignore`.

[slugignore]: https://devcenter.heroku.com/articles/slug-compiler#ignoring-files-with-slugignore

## Remove unused dependencies

Hopefully each Ruby gem in [Bundler groups] is being used by the app.
Do an audit and remove any that aren't used.

[Bundler groups]: http://bundler.io/v1.5/groups.html

## Bundle only what you need for the environment

Heroku accepts an [environment variable][env]
to limit the Ruby gems that it bundles and caches.

[env]: https://devcenter.heroku.com/articles/config-vars

```
heroku config:set BUNDLE_WITHOUT="development:test" --app your-app-name
```
