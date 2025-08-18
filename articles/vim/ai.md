# vim / ai

I type `<Leader>r` from a Markdown file in Vim to send the contents to an LLM
such as <a href="https://platform.openai.com/docs/models#o1" target="_blank">o1</a>
and display the output in a vertical split:

![using an LLM from Neovim](/images/vim-ai.gif)

I can embed files from my codebase into my prompt because the Markdown
is passed through [mdembed](/cmd/mdembed) before it is sent to the LLM:

![embedded code via mdembed](/images/vim-ai-embed.gif)

I can type `<Leader>c` to continue the last conversation.

## Vim config

I install Go via Homebrew, and `mods` and `mdembed` via Go::

```sh
# https://github.com/croaky/laptop/blob/main/laptop.sh

brew "go"

# AI via CLI
go install github.com/charmbracelet/mods@latest
go install github.com/croaky/mdembed@latest
```

I make `mods` and `mdembed` available to Vim by adding
<a href="https://go.dev/wiki/SettingGOPATH" target="_blank">`$(go env GOPATH)/bin`</a>
to my `$PATH`:

```sh
# https://github.com/croaky/laptop/blob/main/shell/zshrc

export PATH="$HOME/go/bin:$PATH"
```

I map `<Leader>r` and `<Leader>c` to the `mdembed` and `mods` pipelines:

```lua
-- https://github.com/croaky/laptop/blob/main/vim/init.lua

-- Helper functions
local function map(mode, lhs, rhs, opts)
  opts = vim.tbl_extend("keep", opts or {}, { noremap = true, silent = false })
  vim.keymap.set(mode, lhs, rhs, opts)
end

local function run_file(key, cmd_template, split_cmd)
  map("n", key, function()
    local cmd = cmd_template:gsub("%%", vim.fn.expand("%:p"))
    vim.cmd(split_cmd)
    vim.cmd("terminal " .. cmd)
  end, { buffer = 0 })
end

-- Markdown
vim.api.nvim_create_autocmd("FileType", {
  pattern = "markdown",
  callback = function()
    run_file("<Leader>r", "cat % | mdembed | mods", "vsplit")
    run_file("<Leader>c", "cat % | mdembed | mods -C", "vsplit")
  end,
})
```

## LLM config

Configure [mods](https://github.com/charmbracelet/mods) with your preferred LLM:

```bash
mods --settings
```

Example:

```yaml
default-model: o1
apis:
  openai:
    base-url: https://api.openai.com/v1
    api-key:
    api-key-env: OPENAI_API_KEY
    models: # https://platform.openai.com/docs/models
      o1-preview:
        aliases: ["o1"]
        max-input-chars: 792000
        fallback: o1-mini
      o1-mini:
        max-input-chars: 500000
  anthropic:
    base-url: https://api.anthropic.com/v1
    api-key:
    api-key-env: ANTHROPIC_API_KEY
    models: # https://docs.anthropic.com/en/docs/about-claude/models
      claude-3-5-sonnet-20240620:
        aliases: ["claude", "sonnet"]
        max-input-chars: 680000
  ollama:
    base-url: http://localhost:11434/api
    models: # https://ollama.com/library
      "deepseek-r1:latest":
        aliases: ["deepseek"]
        max-input-chars: 500000
```

## So what?

I like
<a href="https://sankalp.bearblog.dev/evolution-of-ai-assisted-coding-features-and-developer-interaction-patterns/" target="_blank">this</a>
analogy:

> The lower the gear in a car, the more control you have over the engine but you
> can go with less speed. If you feel in control, go to a higher gear. If you
> are overwhelmed or stuck, go to a lower gear. AI assisted coding is all about
> grokking when you need to gain more granular control and when you need to let
> go of control to move faster. Higher level gears leave more room for errors
> and trust issues.

When I want to work in a lower gear, I use small chunks of code in my Markdown
prompt. When I want to work in a higher gear, I embed one or more files in my
Markdown prompt.

The best results I've had using LLMs are when I prepare a good, isolated prompt.
The longer a conversation goes back-and-forth with follow-ups, the more the LLM
gets confused and won't produce what I want.

I prefer to start with a fresh prompt for each request (`<Leader>r`), editing my
`tmp.md` based on anything I learned from earlier LLM requests,
embedding files from the codebase (some of which I may have edited).
