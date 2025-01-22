# vim / ai

From a Markdown file in Vim, I type `<Leader>r` to send the contents to an LLM
such as OpenAI o1 or o1-mini and have it display the output in a vertical split:

![using an LLM from Neovim](/images/vim-ai.gif)

I can embed files from my codebase into my prompt because the Markdown
is passed through [mdembed](/cmd/mdembed) before it is sent to the LLM:

![embedded code via mdembed](/images/vim-ai-embed.gif)

When I use `<Leader>r`, a new conversation is started.
If I want to continue the last conversation, I use `<Leader>c`.

## Vim config

In `laptop.sh`:

```bash
# https://github.com/croaky/laptop/blob/main/laptop.sh

# AI via CLI
go install github.com/charmbracelet/mods@latest
go install github.com/croaky/mdembed@latest
```

In `init.lua`:

```lua
-- https://github.com/croaky/laptop/blob/main/vim/init.lua

-- Helper functions
local function filetype_autocmd(ft, callback)
	vim.api.nvim_create_autocmd("FileType", {
		pattern = ft,
		callback = callback,
	})
end

-- Markdown
filetype_autocmd("markdown", function()
	-- Run through LLM
	run_file("<Leader>r", "cat % | mdembed | mods", "vsplit")
	run_file("<Leader>c", "cat % | mdembed | mods -C", "vsplit")
end)
```

## LLM config

Configure [mods](https://github.com/charmbracelet/mods) with your preferred LLM:

```bash
mods --settings
```

Example:

```yaml
default-model: o1-preview
apis:
  openai:
    base-url: https://api.openai.com/v1
    api-key:
    api-key-env: OPENAI_API_KEY
    models: # https://platform.openai.com/docs/models
      o1-preview:
        max-input-chars: 792000
        fallback: o1-mini
      o1-mini:
        max-input-chars: 500000
        fallback: gpt-4o
      gpt-4o:
        max-input-chars: 392000
  anthropic:
    base-url: https://api.anthropic.com/v1
    api-key:
    api-key-env: ANTHROPIC_API_KEY
    models: # https://docs.anthropic.com/en/docs/about-claude/models
      claude-3-5-sonnet-20240620:
        aliases: ["claude3.5-sonnet", "claude-3-5-sonnet", "sonnet-3.5"]
        max-input-chars: 680000
      claude-3-opus-20240229:
        aliases: ["claude3-opus", "opus"]
        max-input-chars: 680000
  ollama:
    base-url: http://localhost:11434/api
    models: # https://ollama.com/library
      "llama3:70b":
        aliases: ["llama3"]
        max-input-chars: 650000
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

I also agree with
<a href="https://crawshaw.io/blog/programming-with-llms"
target="_blank">this</a> from David Crawshaw:

> Give an LLM a specific objective and all the background material it needs so
> it can craft a well-contained code review packet. Avoid creating a situation
> with so much complexity and ambiguity that the LLM gets confused and produces
> bad results. I want a blank slate on which to craft a well-contained request.

The best results I've had using LLMs for programming are when I prepare a good
prompt and use a reasoning model like o1. The longer a conversation goes
back-and-forth with follow-up requests, the more the LLM gets confused and won't
produce what I want.

I generally prefer to start with a fresh prompt for each request (`<Leader>r`),
editing my `tmp.md` and embedding files from my codebase,
some of which I may have edited,
based on anything I learned from earlier LLM requests.
