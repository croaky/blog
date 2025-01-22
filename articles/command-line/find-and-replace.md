# Command line find and replace

To find and replace code/text by a file glob:

```bash
replace foo bar **/*.rb
```

Script:

```bash
#!/bin/bash
#
# https://github.com/croaky/laptop/blob/main/bin/replace

find_this="$1"
shift
replace_with="$1"
shift

items=$(ag -l --nocolor "$find_this" "$@")
IFS=$'\n'
for item in $items; do
  sed "s/$find_this/$replace_with/g" "$item" > tmpfile && mv tmpfile "$item"
done
```
