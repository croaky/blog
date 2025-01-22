# cmd / adblock

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
