# blog

Short articles about writing and operating software at <https://dancroak.com>.

## Setup

[Install Go](https://golang.org/doc/install). Then, run:

```
go install ./...
```

This installs a `blog` command-line program from [main.go](main.go):

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
├── code
│   └── example.rb
├── images
│   └── example.png
└-─ theme
    ├── public
    │   └── favicon.ico
    ├── article.html
    └── index.html
```

## Write

Add an article:

```
blog add example-article
```

Edit `articles/example-article.md` in a text editor.
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

Embed code blocks from external files into Markdown like this:

    Instantiate a client:

    ```embed
    code/example.rb instantiate
    ```

This embeds code from `code/example.rb`
between `begindoc` and `enddoc` magic comments
with an id `instantiate`:

```ruby
# begindoc: instantiate
require 'example-sdk'

client = Example::Client.new(
  credential: '...',
  name: 'example',
)
# enddoc: instantiate
```

This way, external files whose code is embedded in the Markdown prose
can be run, linted, or tested in CI.

Add images to the `images` directory.
Refer to them in articles:

```md
![alt text](/images/example.png)
```

## Modify theme

All `theme/public` files are copied to `public`.

The `theme/article.html` file is parsed as a [Go template](https://gowebexamples.com/templates/)
and accepts a data structure like this:

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

## Deploy

Create a static site on [Cloudflare Pages](https://developers.cloudflare.com/pages/framework-guides/deploy-anything/):

- Repository: `https://github.com/croaky/blog`
- Production branch: `main`
- Build command: `git fetch --unshallow && go run main.go build`
- Build output directory: `public`

To deploy the site, commit and push to the GitHub repo.

View deploy logs in the Cloudflare web interface.
