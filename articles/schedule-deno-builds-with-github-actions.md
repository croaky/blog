# Schedule Deno Builds with GitHub Actions

I use a custom static site generator to publish this blog.
It automatically deploys to Deno
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
  deploy:
    name: deploy
    runs-on: ubuntu-latest
    permissions:
      id-token: write # Needed for auth with Deno Deploy
      contents: read # Needed to clone the repository

    steps:
      - name: Clone repository
        uses: actions/checkout@v2

      - name: Build site
        shell: bash
        run: go run main.go build

      - name: Upload to Deno Deploy
        uses: denoland/deployctl@v1
        with:
          project: "croaky-blog"
          entrypoint: https://deno.land/std@0.131.0/http/file_server.ts
          root: public
```

Every day at midnight UTC, GitHub runs the workflow,
which executes a Go program to generate HTML,
and deploys the site to [Deno](https://deno.com).
Then, [Deno serves HTML as a static file server](https://deno.com/blog/deploy-static-files).

When there are blog articles whose scheduled date matches the new UTC date,
they are automatically published by this workflow.

Scheduled workflows are disabled automatically
after 60 days of repository inactivity.
