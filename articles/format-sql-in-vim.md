# Format SQL in Vim

When I save a `.sql` file in Vim,
it auto-formats it with `pgformatter`
using the flags from my `ftplugin/sql.vim`.

Install and configure
[ALE](https://github.com/dense-analysis/ale) in `~/.vimrc`:

```vim
call plug#begin('~/.vim/plugged')
  Plug 'dense-analysis/ale' " :help ale
call plug#end()
```

In `~/.vim/ftplugin/sql.vim`:

```vim
" Auto-fix
let b:ale_fixers = ['pgformatter']
let g:ale_fix_on_save = 1
let b:ale_sql_pgformatter_options = '--function-case 1 --keyword-case 2 --spaces 2 --no-extra-line'
```

In my [laptop.sh](https://github.com/croaky/laptop),
I have a variant of the following that is idempotent;
it will install or update [Homebrew](https://brew.sh/)
and [pgFormatter](https://github.com/darold/pgFormatter):

```bash
# pgformatter
brew install pgformatter

# Vim plugins
curl -fLo "$HOME/.vim/autoload/plug.vim" --create-dirs \
  https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
vim -u "$HOME/.vimrc" +PlugUpdate +PlugClean! +qa
```
