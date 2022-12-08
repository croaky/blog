# Heroku Dataclips Clone

I pair programmed this with [OpenAI's ChatGPT](https://chat.openai.com/chat).
It's a minimum viable [Heroku Dataclips](https://devcenter.heroku.com/articles/dataclips) clone.

```embed
code/heroku-dataclips-clone/main.rb
```

Run locally:

```bash
createdb db
chmod +x main.rb
DATABASE_URL=postgres:///db ./main.rb
```
