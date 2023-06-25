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

The script can safely be run multiple times. I run it most workday mornings.
It is tested on the latest version of macOS on a arm64 (Apple Silicon) chip.
It:

- installs system packages with Homebrew package manager
- configures the shell (zsh)
- creates symlinks for dotfiles to `$HOME`
- installs programming language runtimes
- installs or updates text editor plugins

It also adds a few scripts to `$PATH` whose details are described in these articles:

- [`adblock`](/block-with-etc-hosts)
- [`kill-pid-on-port` and `kill-pid-running`](/kill-pid)
- [`replace`](/find-and-replace)
