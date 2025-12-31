# ai / warp

[Warp](https://www.warp.dev/) is a terminal rebuilt around AI agents.
I use it as my primary terminal for software development.

## How I got here

My AI coding progression has been:

1. **No AI**: just me and Vim.
2. **ChatGPT in browser**: I'd prompt the LLM in the browser, then copy-paste
   code into my editor.
3. **AI chat in Vim**: I'd type `<Leader>r` in Vim to pipe the current buffer's
   Markdown prompt through [mdembed](/cmd/mdembed) and
   [mods](https://github.com/charmbracelet/mods).
   The result would display in a vertical split and I'd yank code into files.
4. **AI agent in terminal**: I now prompt the Warp agent from the command line,
   which applies edits directly to files.
   I review in a diff view. The feedback loop is much tighter.

![Warp AI agent and code review screenshot](/images/warp-code-review.png)

## Why a terminal and not IDE

I'm a long-time Vim user and don't want to give up my editor
for an AI IDE like Cursor or Zed.
I do use [GitHub Copilot for Neovim](https://github.com/github/copilot.vim)
for autocompletion.

I was using [Ghostty](https://ghostty.org/) before Warp.
Ghostty is excellent (fast, native, by Mitchell Hashimoto)
but it's not trying to be an AI tool.

## Why a terminal and not a CLI

Claude Code, Amp, and other AI CLI tools are popular
but lack UI affordances that make me feel especially productive in Warp.

Some things are **nicer with a native UI**:

- Reviewing diffs inline with point-and-click to accept or reject hunks
- Seeing the agent's task list
- Switching between prompting and normal shell

Some things **aren't possible in a raw CLI**:

- Hold a key to speak to the agent via voice transcription
- Drag a screenshot into the prompt

## Pricing

I pay for an Enterprise plan at work and a Build plan at home.
The free plan can give you a taste.
