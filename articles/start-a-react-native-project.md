# Start a React Native Project

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
brew update
brew install node
brew install watchman
brew upgrade
brew cleanup
brew cask cleanup
npm install prettier --global
```

For Vim, install these [plugins](https://github.com/junegunn/vim-plug):

```
Plug 'mxw/vim-jsx'
Plug 'pangloss/vim-javascript'
Plug 'sbdchd/neoformat'
```

Configure [Neoformat](https://github.com/sbdchd/neoformat)
and [Prettier](https://github.com/prettier/prettier):

```vim
" Auto-format on save
augroup fmt
  autocmd!
  autocmd BufWritePre *.js,*.jsx Neoformat prettier
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
