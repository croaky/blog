# Configure Vim for React Native

Why build iOS and Android apps with React Native?

* Prototype rapidly.
* Share code, often 80-90% of the codebase.
* The same designers and developers
  can contribute to the web app and mobile app.
* Save a file and see app reload.
* Use existing text editor.
* The compiled apps are native and performant.

This article describes a macOS setup and Vim configuration
for React Native with Expo and Prettier.

## Set up

Install dependencies using
[Homebrew](http://brew.sh/) and [NPM](https://www.npmjs.org/):

```
brew update-reset
brew install node
brew install watchman
brew upgrade
brew cleanup
brew cask cleanup
npm install prettier --global
```

For Vim, install these [plugins](https://github.com/junegunn/vim-plug):

```
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
```

Configure [ALE](https://github.com/sbdchd/neoformat) and
[COC](https://github.com/neoclide/coc.nvim):

```vim
" https://github.com/neoclide/coc.nvim/wiki/Using-coc-extensions
let g:coc_global_extensions = [
  \ 'coc-prettier',
  \ 'coc-tsserver'
  \ ]

" Jump to definition
nmap <silent> gd <Plug>(coc-definition)
nmap <silent> gy <Plug>(coc-type-definition)
nmap <silent> gi <Plug>(coc-implementation)
nmap <silent> gr <Plug>(coc-references)

" Tab complete
inoremap <silent><expr> <TAB>
  \ pumvisible() ? "\<C-n>" :
  \ <SID>check_back_space() ? "\<TAB>" :
  \ coc#refresh()
inoremap <expr><S-TAB> pumvisible() ? "\<C-p>" : "\<C-h>"

" Auto-fix
let b:ale_fixers = ['prettier']
let g:ale_fix_on_save = 1

" Lint
let b:ale_linters = ['prettier']

augroup ale
  autocmd!

  autocmd VimEnter *
    \ let g:ale_lint_on_enter = 1 |
    \ let g:ale_lint_on_text_changed = 0
augroup END
```

## Develop

[Expo](https://expo.io) is a set of free tools that
allows us to work on JavaScript-only React Native apps
without installing XCode or Android Studio.
Create an account or sign in.

Create a new project from their dropdown menu.
This will download the appropriate files and build the project.

Download the [Expo Client iPhone app](https://itunes.com/apps/exponent).
Open it.
From within the Expo Client app,
scan the project's QR code in the Expo XDE.
Expo Client reloads the app on your phone!

Edit a file such as `screens/HomeScreen.js` in Vim.
Save the file.
Prettier re-formats it on save (no linting necessary)
and Expo Client reloads the app on your phone.
