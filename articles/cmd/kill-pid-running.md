# cmd / kill-pid-running

To kill a running process by its name:

```bash
kill-pid-running sqls
```

Script:

```bash
#!/bin/bash
#
# https://github.com/croaky/laptop/blob/main/bin/kill-pid-running

ps aux | ag "$1" | awk '/$1/ && !/awk/ { print $2 }' | xargs kill
```
