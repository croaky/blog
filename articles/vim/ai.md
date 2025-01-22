# vim / ai

From a Markdown file in Vim, I type `<Leader>r` to send the contents to an LLM
such as OpenAI o1 or o1-mini and have it display the output in a vertical split:

![using an LLM from Neovim](/images/vim-ai.gif)

I can also embed files from my codebase into my prompt using
[mdembed](/cmd/mdembed):

![embedded code via mdembed](/images/vim-ai-embed.gif)
