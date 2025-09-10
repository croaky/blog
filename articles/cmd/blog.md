# cmd / blog

I edit this blog using a simple Go static site generator co-located
with the articles. It is open source at
[croaky/blog](https://github.com/croaky/blog).

## CLI

Install [Go](https://golang.org/doc/install), then run:

```
go install ./...
```

This installs a `blog` command-line program:

```
usage:
  blog serve
  blog build
```

It expects a file layout like this:

```
.
├── articles
│   └── example.md
└── ui
    ├── article.html
    ├── css
    │   └── site.css
    ├── font
    │   ├── et-book-bold-line-figures.woff
    │   ├── et-book-display-italic-old-style-figures.woff
    │   ├── et-book-roman-line-figures.woff
    │   ├── et-book-roman-old-style-figures.woff
    │   └── et-book-semi-bold-old-style-figures.woff
    ├── images
    │   └── favicon.ico
    └── index.html
```

## Write

Edit `articles/example.md`.
It is a [GitHub-Flavored Markdown](https://github.github.com/gfm/) file
with no front matter.

The first line of the file is the article title.
It must be an `<h1>` tag:

```md
# Example Article
```

Markdown headings automatically get IDs for deep linking.
Clicking any `<h2>` navigates to its anchor.

Preview at <http://localhost:2000> with:

```
blog serve
```

Articles are built on-demand when accessed during development.
Requests are logged with timing:

```
   32.1ms 200 GET /cmd/blog
    0.0ms 404 GET /.well-known/appspecific/com.chrome.devtools.json
```

Add images to the `images` directory.
Refer to them in articles:

```md
![alt text](/images/example.png)
```

## Modify UI

All `ui/public` files are copied to `public`.

The `ui/article.html` file is parsed as a [Go template](https://gowebexamples.com/templates/).
Syntax highlighting is generated at build time (no client-side JavaScript highlighting).
`ui/article.html` accepts a data structure like this:

```
{
  Article: {
    ID:            "example-article",
    Title:         "Example Article",
    LastUpdatedOn: "April 15, 2018",  // from git log
    Body:          "<p>Hello, world.</p>",
  },
  CSSPath: "/css/site-a1b2c3d4.css"  // fingerprinted in production
}
```

The `ui/index.html` template is pure HTML.
It is up to the author to decide how to lay out their index
and link to their articles.

CSS files are fingerprinted during production builds for cache-busting.

## Cloudflare Pages

Create a static site on [Cloudflare Pages](https://developers.cloudflare.com/pages/framework-guides/deploy-anything/):

- Repository: `https://github.com/croaky/blog`
- Production branch: `main`
- Build command: `git fetch --unshallow && go run main.go build`
- Build output directory: `public`

The build process:

- Cleans the output directory
- Builds articles concurrently
- Extracts last updated dates from git history

Use the [latest Cloudflare build environment](https://developers.cloudflare.com/pages/configuration/build-image/).

To deploy the site, commit and push to the `main` branch of the GitHub repo.

View deploy logs in the Cloudflare web interface.
