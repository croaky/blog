# Kanban

Here's an example kanban board in [Notion](https://notion.so):

![Kanban board](/images/kanban-board.png)

Add a card in the "To do" column.
It might be a feature, bug, or chore.
Cards are sorted by priority.

To start a new card,
I put my face on the top unassigned card in "To do",
move it to "Doing",
and make a branch:

```bash
git checkout -b my-branch
```

I make my changes and then commit them to version control:

```bash
git add --all
git commit --verbose
```

I push the feature to a remote branch:

```bash
git push
```

I open a pull request from the command line
via [GitHub CLI](https://cli.github.com/):

```bash
gh pr create --fill
gh pr view --web
```

This opens a new pull request in a web browser.

A GitHub webhook starts a
[CI](https://www.martinfowler.com/articles/continuousIntegration.html) build.
Another GitHub webhook posts the pull request to a team
[Slack](https://slack.com) channel.

A teammate clicks the link in the Slack channel.
The teammate comments in-line on the code,
[offers feedback, and approves it][pr].

[pr]: https://help.github.com/articles/about-pull-request-reviews/

Code review before code lands in `main` offers these benefits:

- The whole team learns about new code as it is written.
- Mistakes are caught earlier.
- Coding standards are likely to be established and followed.
- Feedback is likely to be applied.
- Context ("Why did we write this?") is less likely to be forgotten.

I make the suggested changes and commit them:

```bash
git add --all
git commit --verbose
git push
```

We have branch protection rules enabled:
"Require pull request reviews before merging",
"Require status checks to pass before merging",
and "Require branches to be up to date before merging".

Once the pull request has been approved, feedback addressed, and CI has passed,
I press the "Squash and merge" button.
We have the repo settings for commit message set to
"Default to pull request title and description".

After the pull request merges cleanly,
back on the command line in `my-branch`, I run
[this script](https://github.com/croaky/laptop/blob/main/bin/git-post-land):

```bash
git post-land
```

It runs some cleanup and moves me back to `main`:

```bash
git checkout main
git fetch origin
git merge --ff-only origin/main
git branch -D "$branch"
git remote prune origin
```

At this point,
web apps are continuously delivered to a staging environment,
mobile apps are continuously delivered as ad-hoc builds,
and team members are acceptance testing.

When everything looks good,
the code is deployed to production and the card moves to "Done".
