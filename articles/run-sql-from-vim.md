# Run SQL from Vim

When I run `<Leader>r` from a `.sql` file in Vim,
the file's contents are run against my Postgres database through `psql`
and the output is printed to my screen.

The configuration for this is in `~/.vim/ftplugin/sql.vim`:

```vim
" Run current file
nmap <buffer> <Leader>r :!clear && psql -d $(cat .db) -f %<CR>
```

I also have a `.db` file that contains only the local database name:

```
example_development
```

See `man psql` for more detail on the `-d` and `-f` flags.
