# Build Static HTML with GitHub Actions and Deploy to Deno

I use a custom static site generator to publish this blog.
It automatically deploys to Deno:

1. when I merge new articles into my Git repository's `main` branch and
2. every day at midnight UTC to support "scheduled article" functionality

```yaml
name: deno
on:
  push:
    branches:
      - main
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

GitHub runs the workflow,
which executes a Go program to generate HTML,
and deploys the site to [Deno](https://deno.com).
Then, [Deno serves HTML as a static file server](https://deno.com/blog/deploy-static-files).
