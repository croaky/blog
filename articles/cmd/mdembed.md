# cmd / mdembed

`mdembed` is a command line tool to embed file contents in Markdown.

Install:

```bash
go install github.com/croaky/mdembed
```

In `input.md`:

    Embed a whole file:

    ```embed
    f1.rb
    ```

In `f1.rb`:

    puts "hi"

Run:

```bash
cat input.md | mdembed
```

The output will be:

    Embed a whole file:

    ```rb
    # f1.rb
    puts "hi"
    ```

See [examples](https://github.com/croaky/mdembed/tree/main/examples) for more.

## So what?

I wanted the following workflow in Vim:

1. Open `tmp.md` in my project.
2. Write a prompt for an LLM (Large Language Model).
3. Reference other files, or subsets of files, in my project.
4. Hit a key combo (`Space+r` for me) to send all the contents to an LLM.
5. Open a visual split to render the LLM's response.

`mdembed` handles step 3.
I use [mods](https://github.com/charmbracelet/mods) for the LLM steps:

```bash
go install github.com/charmbracelet/mods@latest
```

So, my Vim config runs the following Unix pipeline in a visual split:

```bash
cat example.md | mdembed | mods
```
