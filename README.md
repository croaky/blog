# blog

Short articles about software at <https://dancroak.com>.
Static site generator deployed to [Render](https://render.com/docs/static-sites).

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
├── theme
│   ├── public
│   │   └── favicon.ico
│   ├── article.html
│   └── index.html
└── config.json
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

Configure articles in `config.json`:

```
[
  {
    "description": "Draft scheduled for future date.",
    "id": "article-draft-scheduled",
    "updated": "2050-01-01"
  },
  {
    "canonical": "https://seo.example.com/avoid-duplicate-content-penalty",
    "description": "Canonical article is on a separate site.",
    "id": "article-with-rel-canonical",
    "updated": "2018-01-15"
  }
]
```

The `description` is used for the article page's `meta` description.

The `id` must match a Markdown file `articles/id.md`.
It is also used for the article's URL slug.

The `updated` date can be in the future.
A [GitHub Action is scheduled daily](https://dancroak.com/schedule-deno-builds-with-github-actions)
to auto-publish.

## Modify theme

All `theme/public` files are copied to `public`.

The `theme/*.html` files
are parsed as [Go templates](https://gowebexamples.com/templates/).
The `theme/article.html` template accepts a data structure like this:

```
{
  Article: {
    Body:          "<p>Hello, world.</p>",
    Canonical:     "https://seo.example.com/avoid-duplicate-content-penalty"
    Description:   "Hello, world.",
    ID:            "example-article",
    LastUpdatedOn: "April 15, 2018",
    Title:         "Example Article",
  }
}
```

The `theme/index.html` template accepts a data structure like this:

```
{
  Articles: [
    {
      Body:          "<p>Hello, world.</p>",
      Description:   "Hello, world.",
      ID:            "example-article",
      LastUpdatedIn: "2018 April",
      Title:         "Example Article",
    }
  ],
}
```

## Publish

Configure [Deno Deploy](https://deno.com/deploy):

- Repository: `https://github.com/croaky/blog`
- Production branch: `main`
- Build command: `go run main.go build`
- Public folder: `public`

To publish articles, commit and push to the GitHub repo.

View deploy logs in the Deno Deploy web interface.
