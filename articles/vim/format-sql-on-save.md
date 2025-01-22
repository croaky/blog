# vim / format sql on save

When I save a `.sql` file in Vim,
it auto-formats it with [pgFormatter](https://github.com/darold/pgFormatter)
and [ALE](https://github.com/dense-analysis/ale).

In `laptop.sh`:

```bash
brew install pgformatter

if [ -e "$HOME/.vim/autoload/plug.vim" ]; then
  nvim --headless +PlugUpgrade +qa
else
  curl -fLo "$HOME/.vim/autoload/plug.vim" --create-dirs \
    https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
fi
nvim --headless +PlugUpdate +PlugClean! +qa
nvim --headless +TSUpdate +qa
```

In `~/.vimrc`:

```vim
call plug#begin('~/.vim/plugged')
  Plug 'dense-analysis/ale' " :help ale
call plug#end()
```

In ``~/.vim/ftplugin/sql.vim`

```vim
" Auto-fix
let b:ale_fixers = ['pgformatter'] " 'sqlfmt'
let g:ale_fix_on_save = 1
let b:ale_sql_pgformatter_options = '--function-case 1 --keyword-case 2 --spaces 2 --no-extra-line'

" Run current file
nmap <buffer> <Leader>r :redraw!<CR>:!psql -d $(cat .db) -f % \| less<CR>

" Prepare SQL command with var(s)
nmap <buffer> <Leader>v :redraw!<CR>:!psql -d $(cat .db) -f % -v \| less<SPACE>
```
