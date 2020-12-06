# Postgres \set variable

Within `psql`,
you can `\set` variables and reference them with `:'var-name'`.
For example:

```sql
\set query '%SaaS%'

SELECT
  'https://example.com/companies/' || companies.id AS url,
  companies.name
FROM
  companies
  JOIN notes ON notes.company_id = companies.id
WHERE
  companies.name ILIKE :'query'
  OR companies.description ILIKE :'query'
  OR notes.comments ILIKE :'query'
GROUP BY
  url,
  companies.name
ORDER BY
  companies.name ASC;
```

When I run `<Leader>v` from a `.sql` file in Vim,
I get a prompt to bind my variables.

The configuration for this is in `~/.vim/ftplugin/sql.vim`:

```vim
" Prepare SQL command with var(s)
nmap <buffer> <Leader>v :!clear && psql -d $(cat .db) -f % -v<SPACE>
```

I also have a `.db` file that contains only the local database name:

```
example_development
```

When I press enter,
the variables are bound,
the file's contents are run against my Postgres database through `psql`,
and the output is printed to my screen.

See `man psql` for more detail on the `-d`, `-f`, and `-v` flags.
