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
[`~/.gitconfig`](https://github.com/croaky/laptop/blob/master/dotfiles/git/gitconfig):

```
[merge]
  ff = only
[push]
  default = current
```

## Update

Update the fork with changes in the upstream repo:

```
git checkout master
git fetch upstream
git merge upstream/master
git push
```

If there are no conflicts,
the merge will apply cleanly and
the fork's `master` branch will be sync'ed
with the `upstream` repository.

If there are conflicts, I'll have to rebase and force push:

```
git checkout master
git fetch upstream
git rebase upstream/master
[fix conflicts]
git add .
git rebase --continue
git push -f
```
