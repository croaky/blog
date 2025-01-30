# vim / search project

Projects can be searched for specific text within Vim:

```
:grep sometext
```

`grep` is a built-in command of Vim.
By default, it will use the system's `grep` command.
I override it to use
[ripgrep](https://github.com/BurntSushi/ripgrep/)'s
`rg` command.

## Search word under cursor

In `~/.vimrc`, this config helps me search for the text under the cursor
and show the results in a "quickfix" window:

```lua
vim.opt.grepprg = "rg --vimgrep"
map("n", "K", ':grep! "\\b<C-R><C-W>\\b"<CR>:cw<CR>')
```

It looks like this when `K`
is typed with the cursor over `processMD`:

![Vim quickfix under cursor](/images/quickfix-under-cursor.png)

Cursor over each search result, hit `Enter`, and the file will be opened.

## Search contents of files in project

Also in `~/.vimrc`, this defines a new command `Rg` mapped to `\` to search for
the provided text and open a quickfix window:

```lua
vim.api.nvim_set_keymap("n", "\\", ":Rg<SPACE>", { noremap = true })
```

When `\` is pressed, Vim waits for input:

```vim
:Rg
```

It looks like this when I search for `:Rg lua`:

![Vim qickfix window with search results](/images/quickfix-custom-command.png)
