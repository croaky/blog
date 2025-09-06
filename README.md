# blog

Articles at <https://dancroak.com>

## Setup

[Install Go](https://golang.org/doc/install). Then, run:

```
go install ./...
```

This installs a `blog` command-line program from [main.go](main.go):

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

The `theme/index.html` template is pure HTML.
It is up to the author to decide how to lay out their index
and link to their articles.

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

The CSS is a variant of [Tufte CSS](https://edwardtufte.github.io/tufte-css/) for
typography and layout. Edit `theme/css/site.css` and format it:

```bash
prettier -w theme/css/site.css
```

Instead of footnotes, use sidenotes that appear in the margin on desktop
and toggle on mobile:

```html
Text with a sidenote.
<label for="sn-example" class="margin-toggle sidenote-number"></label>
<input type="checkbox" id="sn-example" class="margin-toggle" />
<span class="sidenote">This is a sidenote that appears in the margin.</span>
```

For notes without numbers, use marginal notes:

```html
Text with a marginal note.
<label for="mn-example" class="margin-toggle">&#8853;</label>
<input type="checkbox" id="mn-example" class="margin-toggle" />
<span class="marginnote">This is a marginal note with a symbol.</span>
```

Make images, code blocks, or other content span the full page width:

```html
<div class="fullwidth">
  <p>This content spans the full width of the page.</p>
</div>
```

Add the HTML directly to Markdown articles.

## Deploy

Create a static site on [Cloudflare Pages](https://developers.cloudflare.com/pages/framework-guides/deploy-anything/):

- Repository: `https://github.com/croaky/blog`
- Production branch: `main`
- Build command: `git fetch --unshallow && go run main.go build`
- Build output directory: `public`

To deploy the site, commit and push to the GitHub repo.

View deploy logs in the Cloudflare web interface.
