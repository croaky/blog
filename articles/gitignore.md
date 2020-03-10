# .gitignore

Configure `~/.gitignore` for all projects on a machine
and `project/.gitignore` for project-specific overrides.
[Example](https://github.com/croaky/laptop/blob/master/dotfiles/git/gitignore):

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
[Git](https://git-scm.com/docs/gitignore) commits,
[ag](https://github.com/ggreer/the_silver_searcher/wiki/Advanced-Usage)
searches,
and [fzf](https://github.com/junegunn/fzf#respecting-gitignore) searches.
