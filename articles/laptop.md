# Laptop

I maintain a [laptop](https://github.com/croaky/laptop) repo
which sets up a macOS machine
as a software development environment.

## Install

Set the `LAPTOP` environment variable to a directory of your choice:

```
export LAPTOP="$HOME/laptop"
```

Clone onto laptop:

```
git clone https://github.com/croaky/laptop.git $LAPTOP
cd $LAPTOP
```

Review:

```
less laptop.sh
```

Run:

```
./laptop.sh
```

## What it sets up

The script is tested on macOS Big Sur (11.3)
on arm64 (Apple Silicon) and x86_64 (Intel) chips.
It:

- installs system packages with Homebrew package manager
- changes shell to Z shell (zsh)
- creates symlinks for dotfiles to `$HOME`
- installs programming language runtimes
- installs or updates Vim plugins

It can safely be run multiple times.
I run it most workday mornings.
