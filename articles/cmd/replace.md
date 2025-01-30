# cmd / replace

To find and replace code/text by a file glob:

```bash
replace foo bar **/*.rb
```

Script:

```bash
#!/bin/bash
#
# https://github.com/croaky/laptop/blob/main/bin/replace

set -euo pipefail

find_this="$1"
shift
replace_with="$1"
shift

items=$(rg -l "$find_this" "$@")

IFS=$'\n'
for item in $items; do
  sed -i '' "s/$find_this/$replace_with/g" "$item"
done
```
