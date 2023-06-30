# Heroku Dataclips clone

This is a minimum viable [Heroku Dataclips](https://devcenter.heroku.com/articles/dataclips) clone.

```embed
code/heroku-dataclips-clone/main.rb
```

Run locally:

```bash
createdb db
chmod +x main.rb
DATABASE_URL=postgres:///db ./main.rb
```
