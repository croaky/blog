# go / static site

This blog is written with a small, custom Go static site generator. It is
auto-deployed to Cloudflare Pages when Markdown articles are merged into the
[croaky/blog](https://github.com/croaky/blog) Git repo's `main` branch.

## Go CLI

Install [Go](https://golang.org/doc/install), then run:

```
go install ./...
```

This installs a `blog` command-line program:

```
usage:
  blog add <article-url-slug>
  blog serve
  blog build
```

It expects a file layout like this:

```
.
├── articles
│   └── example.md
├── images
│   └── example.png
└-─ theme
    ├── public
    │   └── favicon.ico
    ├── article.html
    └── index.html
```

## Write

Edit `articles/example.md` in a text editor.
It is a [GitHub-Flavored Markdown](https://github.github.com/gfm/) file
with no front matter.

The first line of the file is the article title.
It must be an `<h1>` tag:

```md
# Example Article
```

Preview at <http://localhost:2000> with:

```
blog serve
```

Add images to the `images` directory.
Refer to them in articles:

```md
![alt text](/images/example.png)
```

## Modify theme

All `theme/public` files are copied to `public`.

The `theme/article.html` file is parsed as a [Go template](https://gowebexamples.com/templates/).
Syntax highlighting is generated at build time (no client-side JavaScript highlighting).
`theme/article.html` accepts a data structure like this:

```
{
  Article: {
    ID:            "example-article",
    Title:         "Example Article",
    LastUpdatedOn: "April 15, 2018",
    Body:          "<p>Hello, world.</p>",
  }
}
```

The `theme/index.html` template is pure HTML.
It is up to the author to decide how to lay out their index
and link to their articles.

## Cloudflare Pages

Create a static site on [Cloudflare Pages](https://developers.cloudflare.com/pages/framework-guides/deploy-anything/):

- Repository: `https://github.com/croaky/blog`
- Production branch: `main`
- Build command: `git fetch --unshallow && go run main.go build`
- Build output directory: `public`

Use [the Cloudflare v2 build
environment](https://developers.cloudflare.com/pages/platform/language-support-and-tools/)
for the latest version of Go.

To deploy the site, commit and push to the GitHub repo.

View deploy logs in the Cloudflare web interface.
