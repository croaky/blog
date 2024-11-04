# Git workflow

When I think of a change to my web app,
I add a card to a kanban board in [Notion](https://notion.com):

![Kanban board](/images/kanban-board.png)

The card might refer to a feature, bug, or chore.
Cards are sorted by priority in each column.

When I'm ready to work on the change,
I assign myself to the card,
move it to "Doing",
and make a Git branch:

```bash
git checkout -b my-branch
```

I edit the codebase and commit the changed files to version control:

```bash
git add --all
git commit --verbose
```

I push the feature to a remote branch:

```bash
git push
```

This only pushes `my-branch` to GitHub because I have this setting in
my `~/.gitconfig`:

```
[push]
  default = current
```

I open a pull request (PR) from the command line
via [GitHub CLI](https://cli.github.com/):

```bash
gh pr create --fill
```

This triggers webhooks that create:

1. a [CI](https://www.martinfowler.com/articles/continuousIntegration.html) build
2. a [Slack](https://slack.com) message in my team's channel

I open new pull request in a web browser:

```bash
gh pr view --web
```

I review the code again.
I may push follow-up changes or edit the PR description.

When CI passes,
I open the Slack thread and ask a teammate
to review:

```
@buddy PTAL
```

"PTAL" means "Please Take A Look".

When they are ready to review,
they add an 👀 emoji to the thread
and open the PR in a browser.

They comment in-line on the code,
[offer feedback, and approve it](https://help.github.com/articles/about-pull-request-reviews/).

I make any suggested changes and commit them.

My repo has these settings:

1. Require pull request reviews before merging
2. Require status checks to pass before merging
3. Require branches to be up to date before merging
4. Default commit message to pull request title and description

I press the "Squash and merge" button.

GitHub triggers a webhook to deploy the `main` branch
to my staging environment on [Render](https://render.com).

I acceptance test on staging.
When everything looks good,
I move back to the command line.

In `my-branch`, I run
[this script](https://github.com/croaky/laptop/blob/main/bin/git-post-land),
which runs some cleanup and moves me back to `main`:

```bash
git post-land
```

I deploy to production with deploy script:

```bash
deploy-prod
```

I move the card on the kanban board to "Done".
