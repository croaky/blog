# .gitignore

Configure `~/.gitignore` for all projects on a machine
and `project/.gitignore` for project-specific overrides.

Example:

```
*.lock
*.log
*.pyc
*.sw[nop]
.DS_Store
.bundle
node_modules
public
tmp
vendor
```

Directories and files matching these patterns will be ignored for
[git](https://git-scm.com/docs/gitignore) commits,
[ag](https://github.com/ggreer/the_silver_searcher/wiki/Advanced-Usage)
searches,
and [fzf](https://github.com/junegunn/fzf#respecting-gitignore) searches.
