# laptop.sh

I maintain a [laptop](https://github.com/croaky/laptop) repo
which sets up a macOS machine
as a software development environment.

## Install

Set the `OK` environment variable to a directory of your choice:

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

The script is tested on macOS Mojave (10.14).
It:

* installs system packages with Homebrew package manager
* changes shell to Z shell (zsh)
* installs SF Mono font for Kitty terminal and Vim editor
* creates symlinks from `$LAPTOP/dotfiles` to `$HOME`
* installs or updates Vim plugins
* installs programming language runtimes

This script can be run safely multiple times.
I run it most working mornings.
