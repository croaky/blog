# Update a GitHub Fork

After forking a GitHub repo to my personal account,
I want to update the fork with changes in the "upstream" repo.

## Fork

Get the [GitHub CLI](https://cli.github.com/) tool:

```
brew install github/gh/gh
brew update
brew upgrade gh
```

Fork the upstream repo:

```
gh repo fork org/project
```

The `upstream` remote is automatically set to the upstream repository.

I also have these relevant settings in my
[`~/.gitconfig`](https://github.com/croaky/laptop/blob/main/dotfiles/git/gitconfig):

```
[merge]
  ff = only
[push]
  default = current
```

## Update

Update the fork with changes in the upstream repo:

```
git checkout main
git fetch upstream
git merge upstream/main
git push
```
