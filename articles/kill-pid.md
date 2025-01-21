# Kill PID scripts

I use a [`kill-pid-on-port` script](https://github.com/croaky/laptop/blob/main/bin/kill-pid-on-port)
to kill processes listening on a given port:

```bash
kill-pid-on-port 3000
```

The implementation is:

```bash
lsof -n -i :"$1" | grep LISTEN | awk '{ print $2 }' | xargs kill
```

I use a [`kill-pid-running`
script](https://github.com/croaky/laptop/blob/main/bin/kill-pid-running)
to kills running process by its name:

```bash
kill-pid-running sqls
```
