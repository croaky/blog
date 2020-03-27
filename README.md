# blog

Articles about making software at <https://dancroak.com>.
Bespoke static site generator
designed to be deployed to [Netlify](https://www.netlify.com/).

## Setup

[Install Go](https://dancroak.com/install-or-upgrade-go-on-macos).
Then, run:

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
│   │   └── _headers
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
![alt text](images/example.png)
```

Configure articles in `config.json`:

```
[
  {
    "description": "Draft scheduled for future date.",
    "id": "article-draft-scheduled",
    "last_updated": "2050-01-01",
    "tags": [
      "go",
      "react"
    ]
  },
  {
    "description": "Redirect old URL slugs.",
    "id": "article-with-redirects",
    "last_updated": "2018-02-01",
    "tags": [
      "go"
    ]
    "redirects": [
      "/article-original-name",
      "/article-renamed-again",
      "/this-feature-works-only-on-netlify",
    ]
  },
  {
    "description": "Canonical article is on a separate site.",
    "canonical": "https://seo.example.com/avoid-duplicate-content-penalty",
    "id": "article-with-rel-canonical",
    "last_updated": "2018-01-15",
    "tags": [
      "go"
    ]
  }
]
```

The `id` must match a Markdown file `articles/id.md`.
It is also used for the article's URL slug.

The `description` is used for the article page's `meta` description.

## Modify theme

All `theme/public` files are copied to `public`.
`theme/public/_headers` are
[Netlify Headers](https://www.netlify.com/docs/headers-and-basic-auth/).

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
    Tags:          ["go", "unix"],
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
      Tags:          ["go", "unix"],
      Title:         "Example Article",
    }
  ],
  Tags: ["go", "unix"],
}
```

## Publish

Configure [Netlify](https://www.netlify.com):

* Repository: `https://github.com/croaky/blog`
* Branch: `master`
* Build Cmd: `go run main.go build`
* Public folder: `public`

To publish articles, commit and push to the GitHub repo.

View deploy logs in the Netlify web interface.
