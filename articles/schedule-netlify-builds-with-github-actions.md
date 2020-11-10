# Schedule Netlify Builds with GitHub Actions

I use a [custom static site generator](https://github.com/croaky/blog)
to publish this blog.
It automatically deploys to Netlify
when I merge new articles into my Git repository's `main` branch.
To support a "scheduled article" feature,
I have configured a
[GitHub Actions scheduled workflow](https://help.github.com/en/actions/reference/workflow-syntax-for-github-actions#onschedule):

```yaml
name: daily publish

on:
  schedule:
    - cron: "0 0 * * *" # every day at midnight UTC

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Trigger Netlify build
        shell: bash
        env:
          NETLIFY_BUILD_HOOK: ${{ secrets.NETLIFY_BUILD_HOOK }}
        run: curl -X POST -d {} "$NETLIFY_BUILD_HOOK"
```

Every day at midnight UTC, GitHub runs the workflow,
deploying the site using a
[Netlify build hook](https://docs.netlify.com/configure-builds/build-hooks/).
The build hook is a URL which I've stored as a
[GitHub encrypted secret](https://help.github.com/en/actions/configuring-and-managing-workflows/creating-and-storing-encrypted-secrets#using-encrypted-secrets-in-a-workflow).

When there are articles whose scheduled date matches the new UTC date,
they are automatically published by this workflow.

Scheduled workflows are disabled automatically
after 60 days of repository inactivity.
