# cmd / mdembed

`mdembed` embeds programming file contents in Markdown.

It is [open source](https://github.com/croaky/mdembed).

Install:

```bash
go install github.com/croaky/mdembed@latest
```

## Use

In `prompt.md`:

````md
# Context

```embed
cmd/server.go
lib/**/*.rb
```
````

Run:

```bash
cat prompt.md | mdembed
```

The output will contain code blocks with the file contents,
each with a header comment showing the filename:

````md
# Context

```go
// cmd/server.go
// file contents...
```

```rb
# lib/foo.rb
# file contents...
```

```rb
# lib/bar.rb
# file contents...
```
````

## Block markers

Embed specific sections of a file using `emdo` and `emdone` markers:

```go
func main() {
    // emdo setup
    cfg := loadConfig()
    db := connectDB(cfg)
    // emdone setup

    // ... rest of code
}
```

Reference the block by name:

````md
```embed
main.go setup
```
````

## Recursive embeds

Referencing another Markdown file embeds its contents directly,
recursively processing any embed blocks it contains.

## Why

I wanted a workflow in Vim:

1. Open `tmp.md` in my project.
2. Write a prompt for an LLM, referencing files in my project.
3. Hit a key combo to send the contents to an LLM.

`mdembed` handles the Markdown parsing.
[mods](https://github.com/charmbracelet/mods) handles the LLM:

```bash
cat prompt.md | mdembed | mods
```
