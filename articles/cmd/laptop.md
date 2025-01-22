# cmd / laptop

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
