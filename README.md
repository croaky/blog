# blog

Articles about making software at <https://dancroak.com>.
Bespoke static site generator
designed to be deployed to [Netlify](https://www.netlify.com/).

# Setup

[Install Go](https://dancroak.com/install-or-upgrade-go-on-macos).
Then, run:

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

`blog` is a static site generator featuring:

* Markdown files with no front matter
* Local preview server
* Images
* Embedded code blocks
* Drafts
* Tags
* "Last updated" timestamp
* Redirects
* `rel=canonical` tags
* Responsive design
* PageSpeed Insights performance score of 100
* Mozilla Observatory security grade of A+

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

Add images to the `articles/images` directory.
Refer to them in articles via relative path:

```md
![alt text](images/example.png)
```

In addition to
[fenced code blocks](https://github.github.com/gfm/#fenced-code-blocks),
`blog` recognizes a special `embed`
[info string](https://github.github.com/gfm/#info-string).
This Markdown...

    Instantiate a client:

    ```embed
    code/example.rb instantiate
    ```

...embeds code from `code/example.rb`
between `begindoc` and `enddoc` magic comments:

```ruby
# begindoc: instantiate
require 'example-sdk'

client = Example::Client.new(
  credential: '...',
  name: 'example',
)
# enddoc: instantiate
```

The magic comments demarcate code blocks by id.
In this example, the id is `instantiate`.

This allows you to run, lint, and test embedded code
separate from Markdown prose.

## Configure

Configure blog in `config.json`:

```
{
  "articles": [
    {
      "id": "article-is-draft-if-future-date",
      "last_updated": "2050-01-01",
      "tags": [
        "go",
        "unix"
      ]
    },
    {
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
      "canonical": "https://seo.example.com/avoid-duplicate-content-penalty",
      "id": "article-with-rel-canonical",
      "last_updated": "2018-01-15",
      "tags": [
        "go"
      ]
    }
  ]
}
```

## Modify theme

The `theme` directory's files can be modified
to customize the blog's HTTP headers, HTML, CSS, and JavaScript.

```
.
├── index.html
└── public
    ├── _headers
    └── favicon.ico
```

The `_headers` file is copied to `public/_headers` to be used as
[Netlify Headers](https://www.netlify.com/docs/headers-and-basic-auth/).

The `.html` files
are parsed as [Go templates](https://gowebexamples.com/templates/)
by `blog`.

The `article.html` template accepts a data structure like this:

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

The `index.html` template accepts a data structure like this:

```
{
  Articles: [
    {
      Body:          "<p>Hello, world.</p>",
      Canonical:     "https://seo.example.com/avoid-duplicate-content-penalty"
      Description:   "Hello, world.",
      ID:            "example-article",
      LastUpdated:   "2018-04-15",
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
