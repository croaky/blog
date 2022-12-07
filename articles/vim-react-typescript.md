# Configure Vim for React and TypeScript

This article describes a macOS setup and Vim 8 configuration
for React and TypeScript with syntax highlighting,
tab completion, auto-formatting, linting, and jump to definition.

Refresh dependencies using [Homebrew](https://brew.sh/):

```
brew update-reset
brew install node
brew install vim
brew upgrade
brew cleanup
```

Example `~/.vimrc` using
[Vim Plug](https://github.com/junegunn/vim-plug) (plugin manager),
[ALE](https://github.com/sbdchd/neoformat) (linting, auto-fixing/fmt'ing) and
[COC](https://github.com/neoclide/coc.nvim)
(Language Server Protocol for tab completion, jump to definition):

```vim
call plug#begin('~/.vim/plugged')
  " :help ale
  Plug 'dense-analysis/ale'

  " Language Server Protocol
  Plug 'neoclide/coc.nvim', { 'branch': 'release' }

  " Frontend
  Plug 'leafgarland/typescript-vim'
  Plug 'mxw/vim-jsx'
  Plug 'pangloss/vim-javascript'
call plug#end()

" Lint with ALE
augroup ale
  autocmd!

  autocmd VimEnter *
    \ let g:ale_lint_on_enter = 1 |
    \ let g:ale_lint_on_text_changed = 0
augroup END

" https://github.com/neoclide/coc.nvim/wiki/Using-coc-extensions
let g:coc_global_extensions = [
  \ 'coc-prettier',
  \ 'coc-tsserver'
  \ ]

" Tab complete with COC
inoremap <silent><expr> <TAB>
  \ pumvisible() ? "\<C-n>" :
  \ <SID>check_back_space() ? "\<TAB>" :
  \ coc#refresh()
inoremap <expr><S-TAB> pumvisible() ? "\<C-p>" : "\<C-h>"

" Jump to definition, implementation, or call sites
nnoremap <silent> gd <Plug>(coc-definition)
nnoremap <silent> gi <Plug>(coc-implementation)
nnoremap <silent> gr <Plug>(coc-references)

" When all else fails, grep word under cursor with "K"
nnoremap K :grep! "\b<C-R><C-W>\b"<CR>:cw<CR>
```

Example `~/.vim/ftplugin/typescript.vim`:

```vim
" Auto-fix
let b:ale_fixers = ['prettier']
let g:ale_fix_on_save = 1

" Lint
let b:ale_linters = ['tsserver']
```

[Install Prettier](https://prettier.io/docs/en/install.html)
in your project directory:

```bash
# with NPM
npm install --save-dev --save-exact prettier

# or, with Yarn
yarn add --dev --exact prettier
```
