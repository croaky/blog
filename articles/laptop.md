# Laptop

I maintain a [laptop](https://github.com/croaky/laptop) repo
which sets up a macOS machine
as a software development environment.

## Install

Clone onto laptop:

```
export LAPTOP="$HOME/laptop"
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

The script can safely be run multiple times.
I run it most workday mornings.
It is tested on the latest version of macOS on a arm64 (Apple Silicon) chip.
It:

- installs system packages with Homebrew
- sets the shell to [zsh](https://www.zsh.org/)
- sets the terminal to [kitty](https://sw.kovidgoyal.net/kitty/)
- symlinks dotfiles from the repo to `$HOME`
- installs Go and Ruby programming languages
- sets up [Neovim](https://neovim.io/)
- configures Neovim with LSP, completion, fuzzy-finding
- adds a few scripts to `$PATH`

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

## Kill PID scripts

[The `kill-pid-on-port` script](https://github.com/croaky/laptop/blob/main/bin/kill-pid-on-port)
kills processes listening on a given port:

```bash
kill-pid-on-port 3000
```

[The `kill-pid-running`
script](https://github.com/croaky/laptop/blob/main/bin/kill-pid-running)
kills running process by its name:

```bash
kill-pid-running sqls
```

## Find and replace script

[The `replace`
script](https://github.com/croaky/laptop/blob/main/bin/replace)
finds and replaces code/text by a file glob:

```bash
replace foo bar **/*.rb
```
