# cmd / kill-pid-on-port

To kill a process listening on a given port:

```bash
kill-pid-on-port 3000
```

Script:

```bash
#!/bin/bash
#
# https://github.com/croaky/laptop/blob/main/bin/kill-pid-on-port

set -euo pipefail

lsof -n -i :"$1" | grep LISTEN | awk '{ print $2 }' | xargs kill
```

Why not `lsof -ti :3000 | xargs kill`?

`lsof -ti` returns all processes on the port—servers and clients.
`grep LISTEN` ensures we only kill the server.
