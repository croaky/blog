# Laptop

I maintain a [laptop](https://github.com/croaky/laptop) repo
which sets up a macOS machine
as a software development environment.

## Install

Clone onto laptop:

```bash
export LAPTOP="$HOME/laptop"
git clone https://github.com/croaky/laptop.git $LAPTOP
cd $LAPTOP
```

Review:

```bash
less laptop.sh
```

Run:

```bash
./laptop.sh
```

## What it sets up

The script can safely be run multiple times.
I run it most workday mornings.
It is tested on the latest version of macOS on a arm64 (Apple Silicon) chip.
It:

- installs system packages with Homebrew
- sets the shell to [zsh](https://www.zsh.org/)
- sets the terminal to [Ghostty](https://ghostty.org/)
- symlinks dotfiles from the repo to `$HOME`
- installs Go and Ruby programming languages
- sets up [Neovim](https://neovim.io/)
- configures Neovim with LSP, completion, fuzzy-finding
- adds a few scripts to `$PATH`

## Run SQL queries

When I run `<Leader>r` from a `.sql` file in Vim,
the file's contents are run against my Postgres database through `psql`
and the output is printed to my screen.

In `~/.vim/ftplugin/sql.vim`:

```vim
" Run current file
nmap <buffer> <Leader>r :redraw!<CR>:!psql -d $(cat .db) -f %<CR>
```

I also have a `.db` file that contains only the local database name:

```
example_development
```

See `man psql` for more detail on the `-d` and `-f` flags.

## Debug slow Postgres queries

If the query is slow, I add this to the top of the file:

```
EXPLAIN (ANALYZE, COSTS, VERBOSE, BUFFERS, FORMAT JSON)
```

Then, run:

```bash
:!psql -qAt -d $(cat .db) -f % | pbcopy
```

Paste into <http://tatiyants.com/pev/#/plans/new>
and delete the trailing line to make it valid JSON:

```
Time: 1111.111 ms (00:01.111)
```

The output is an interactive visualization that makes it
easy to identify which parts of the query are
slowest, largest, and costliest.

![EXPLAIN visualizer](/images/postgres-explain-visualizer.png)

## Run tests

Test-driven development thrives on a tight feedback loop
but switching from editor to shell
to manually run specs is inefficient.

The [vim-test](https://github.com/vim-test/vim-test) plugin
exposes commands such as `:TestNearest`, `:TestFile`, and `:TestLast`,
which I bind to `<Leader>s`, `<Leader>t`, and `<Leader>l`.

Cursor over any line within an RSpec spec like this:

```ruby
describe RecipientInterceptor do
  it 'overrides to/cc/bcc fields' do
    Mail.register_interceptor RecipientInterceptor.new(recipient_string)

    response = deliver_mail

    expect(response.to).to eq [recipient_string]
    expect(response.cc).to eq []
    expect(response.bcc).to eq []
  end
end
```

Type `<Leader>s`:

```
rspec spec/recipient_interceptor_spec.rb:4
.

Finished in 0.03059 seconds
1 example, 0 failures
```

The screen is overtaken by a shell that runs only the focused spec.

Feeling good that this new spec passes,
run the whole file's specs with `<Leader>t`
to make sure the class's entire functionality is still intact:

```
rspec spec/recipient_interceptor_spec.rb
......

Finished in 0.17752 seconds
6 examples, 0 failures
```

Red, green, refactor.
From the program:

```ruby
def delivering_email(message)
  add_custom_headers message
  add_subject_prefix message
  message.to = @recipients
  message.cc = []
  message.bcc = []
end
```

Run `<Leader>l` without having to switch back to the spec:

```
rspec spec/recipient_interceptor_spec.rb
......

Finished in 0.17752 seconds
6 examples, 0 failures
```

Running specs in tight feedback loops
reduces switching cost between editor and shell,
making test-driven development easier.

## Search projects in Vim

Projects can be searched for specific text within Vim:

```
:grep sometext
```

`grep` is a built-in command of Vim.
By default, it will use the system's `grep` command.
I override it to use
[The Silver Searcher](https://github.com/ggreer/the_silver_searcher)'s
`ag` command.

In `~/.vimrc`:

```vim
" The Silver Searcher
if executable('ag')
  " Use ag over grep
  set grepprg=ag\ --nogroup\ --nocolor

  " Use ag in CtrlP for listing files. Lightning fast and respects .gitignore
  let g:ctrlp_user_command = 'ag %s -l --nocolor -g ""'

  " ag is fast enough that CtrlP doesn't need to cache
  let g:ctrlp_use_caching = 0
endif
```

This searches for the text under the cursor
and shows the results in a "quickfix" window:

```vim
" bind <Leader>k to grep word under cursor
nnoremap <Leader>k :grep! "\b<C-R><C-W>\b"<CR>:cw<CR>
```

It looks like this when `<Leader>k`
is typed with the cursor over `SubscriptionMailer`:

![Vim quickfix under cursor](/images/quickfix-under-cursor.png)

Cursor over each search result, hit `Enter`, and the file will be opened.

This defines a new command `Ag` to search for the provided text
and open a quickfix window:

```vim
" bind \ (backward slash) to grep shortcut
command -nargs=+ -complete=file -bar Ag silent! grep! <args>|cwindow|redraw!
```

Map it to any character, such as `\`:

```vim
nnoremap \ :Ag<SPACE>
```

When `\` is pressed, Vim waits for input:

```vim
:Ag
```

Standard `ag` arguments may be passed in at this point:

```vim
:Ag -i Stripe app/models
```

Hitting `Enter` results in:

![Vim qickfix window with search results](/images/quickfix-custom-command.png)

## Ad block script

To improve speed, privacy, and safety on my laptop,
[the `adblock` script](https://github.com/croaky/laptop/blob/main/bin/adblock)
blocks ads, trackers, and malicious websites at the DNS host level:

```bash
adblock
```

Unlike browser extension ad blockers,
it works on all apps on my device (not only web browsers).

Unlike DNS sinkholes,
it only works on my laptop (not phones, tablets on the network)
but it does not require an additional always-on device such as a Raspberry Pi
and it works reliably when using the laptop away from home.

To disable and re-enable it:

```bash
adblock undo
adblock
```
