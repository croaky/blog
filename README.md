# blog

Articles about making software at <https://dancroak.com>.

## genblog

A static blog generator featuring:

* Markdown files with no front matter
* Local preview server
* JSON feed
* Images
* Embedded code blocks
* Drafts
* Anonymous, single author, or multiple authors
* Tags
* "Last updated" timestamp
* Redirects
* `rel=canonical` tags

The theme features:

* Responsive design
* PageSpeed Insights performance score of 100
* Mozilla Observatory security grade of A+

The `./genblog` script requires Go.
On macOS, install Go with:

```
gover="1.14"
if ! go version | grep -Fq "$gover"; then
  sudo rm -rf /usr/local/go
  curl "https://dl.google.com/go/go$gover.darwin-amd64.tar.gz" | \
    sudo tar xz -C /usr/local
fi
```

## Write

Add an article:

```
./genblog add example-article
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
./genblog serve
```

See the [JSON feed](https://jsonfeed.org/) at <http://localhost:2000/feed.json>.

Add images to the `articles/images` directory.
Refer to them in articles via relative path:

```md
![alt text](images/example.png)
```

In addition to
[fenced code blocks](https://github.github.com/gfm/#fenced-code-blocks),
`./genblog` recognizes a special `embed`
[info string](https://github.github.com/gfm/#info-string).
This Markdown...

    Instantiate a client:

    ```embed
    example.rb instantiate
    ```

...embeds code from `articles/code/example.rb`
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

```json
{
  "blog": {
    "name": "Dan Croak",
    "url": "https://dancroak.com"
  },
  "articles": [
    {
      "id": "article-is-draft-if-published-is-future-date",
      "published": "2050-01-01"
    },
    {
      "id": "article-with-anonymous-author",
      "published": "2018-04-15"
    },
    {
      "author": "Alice",
      "id": "article-with-single-author",
      "published": "2018-04-01"
    },
    {
      "author": "Alice and Bob",
      "id": "article-with-multiple-authors",
      "published": "2018-03-15"
    },
    {
      "id": "article-with-tags",
      "published": "2018-03-01",
      "tags": [
        "go",
        "unix"
      ]
    },
    {
      "id": "article-with-updated-date",
      "published": "2018-02-15",
      "updated": "2018-02-20"
    },
    {
      "id": "article-with-redirects",
      "published": "2018-02-01",
      "redirects": [
        "/article-original-name",
        "/article-renamed-again",
        "/this-feature-works-only-on-netlify",
      ]
    },
    {
      "canonical": "https://seo.example.com/avoid-duplicate-content-penalty",
      "id": "article-with-rel-canonical",
      "published": "2018-01-15"
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
by `genblog`.

The `article.html` file accepts a data structure like this:

```
{
  Blog: {
    "name": "Dan Croak",
    "url": "https://dancroak.com"
  }
  Article: {
    Author:        "Alice",
    Body:          "<p>Hello, world.</p>",
    Canonical:     "https://seo.example.com/avoid-duplicate-content-penalty"
    ID:            "example-article",
    LastUpdated:   "2018-04-15",
    LastUpdatedIn: "2018 April",
    LastUpdatedOn: "April 15, 2018",
    Published:     "2018-04-10",
    Tags:          ["go", "unix"],
    Title:         "Example Article",
    Updated:       "2018-04-15",
  }
}
```

The `index.html` file accepts a data structure like this:

```
{
  Blog: {
    "name": "Dan Croak",
    "url": "https://dancroak.com"
  },
  Articles: [
    {
      Author:        "Alice",
      Body:          "<p>Hello, world.</p>",
      Canonical:     "https://seo.example.com/avoid-duplicate-content-penalty"
      ID:            "example-article",
      LastUpdated:   "2018-04-15",
      LastUpdatedIn: "2018 April",
      LastUpdatedOn: "April 15, 2018",
      Published:     "2018-04-10",
      Tags:          ["go", "unix"],
      Title:         "Example Article",
      Updated:       "2018-04-15",
    }
  ],
  Tags: ["go", "unix"],
}
```

## Publish

Configure [Netlify](https://www.netlify.com):

* Repository: `https://github.com/croaky/blog`
* Branch: `master`
* Build Cmd: `./genblog build`
* Public folder: `public`

To publish articles, commit and push to the GitHub repo.

View deploy logs in the Netlify web interface.
