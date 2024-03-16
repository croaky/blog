/*
Add an article:

	blog add <article-url-slug>

Serve site on localhost:

	blog serve

Build site (HTML, images, code) to `public/`:

	blog build

When building an article, the program scans line-by-line.
It first extracts the title from the first line...

```
# Article Title
```

When the scanner encounters an "embed" code fence like this...

```embed
code/example.rb id
```

...it loads the file and finds "magic comments" in the source like this...

# begindoc: id
puts "here"
# enddoc: id

The lines between magic comments
are embedded back in the original code fence.

To embed the entire file, don't include an ID or magic comments...

```embed
code/example.rb
```

Bad input in the Markdown document or source code file
will stop the program with a non-zero exit code and error text.

When the scanner is done, it passes its processed result to a standard Markdown
compiler to finish compilation to HTML.
*/

package main

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/chroma/v2"
	chromaHTML "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	markdownHTML "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	blogURL = "https://dancroak.com"
	wd      string
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	wd, err = os.Getwd()
	fatal(err, "Failed to get working directory")

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			usage()
		}
		add(os.Args[2])
		fmt.Println("Added ./articles/" + os.Args[2] + ".md")
	case "serve":
		fmt.Println("Serving at http://localhost:2000")
		serve(":2000")
	case "build":
		build("public")
		fmt.Println("Built at ./public")
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:\n  blog add <article-url-slug>\n  blog serve\n  blog build")
	os.Exit(2)
}

func fatal(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err)
		os.Exit(1)
	}
}

type Article struct {
	ID        string
	Title     string
	UpdatedOn string
	Body      template.HTML
}

func add(id string) {
	title := cases.Title(language.Und).String(strings.ReplaceAll(strings.ReplaceAll(id, "-", " "), "_", " "))
	content := []byte("# " + title + "\n\n\n")
	fatal(os.WriteFile(filepath.Join(wd, "articles", id+".md"), content, 0644), "Failed to add article")
}

func serve(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Normalize the path
		path := strings.TrimSuffix(r.URL.Path, "/")
		if path == "" {
			path = "/"
		}

		// Serve the index page for the root path
		if path == "/" {
			http.ServeFile(w, r, filepath.Join(wd, "theme", "index.html"))
		} else if strings.HasPrefix(path, "/images") {
			// Serve static files
			fs := http.StripPrefix("/images", http.FileServer(http.Dir(filepath.Join(wd, "theme", "images"))))
			fs.ServeHTTP(w, r)
		} else {
			// Build and serve the article for non-root paths
			buildArticle(strings.TrimPrefix(path, "/"))
			http.ServeFile(w, r, filepath.Join(wd, "public", path, "index.html"))
		}

		fmt.Printf("%7.1f ms %s %s\n", float64(time.Since(startTime))/float64(time.Millisecond), r.Method, path)
	})
	fatal(http.ListenAndServe(addr, nil), "Failed to serve")
}

func build(outputDir string) {
	// Ensure the output directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	fatal(err, "Failed to create output directory")

	// Clean the output directory
	dirEntries, err := os.ReadDir(outputDir)
	fatal(err, "Failed to read output directory")
	for _, d := range dirEntries {
		fatal(os.RemoveAll(filepath.Join(outputDir, d.Name())), "Failed to remove file in output directory")
	}

	// Copy theme static files
	copyDir(filepath.Join(wd, "theme", "index.html"), filepath.Join(outputDir, "index.html"))
	copyDir(filepath.Join(wd, "theme", "images"), filepath.Join(outputDir, "images"))

	// Build article pages
	page := template.Must(template.ParseFiles(filepath.Join(wd, "theme", "article.html")))
	articles := load()

	var wg sync.WaitGroup
	for _, a := range articles {
		wg.Add(1)
		go func(a Article) {
			defer wg.Done()
			articleDir := filepath.Join(outputDir, a.ID)
			fatal(os.MkdirAll(articleDir, os.ModePerm), "Failed to create article directory")
			f, err := os.Create(filepath.Join(articleDir, "index.html"))
			fatal(err, "Failed to create article index.html")
			fatal(page.Execute(f, struct{ Article Article }{a}), "Failed to execute article template")
		}(a)
	}
	wg.Wait()
}

func copyFile(src, dst string) {
	source, err := os.Open(src)
	fatal(err, "Failed to open source file")
	defer source.Close()

	destination, err := os.Create(dst)
	fatal(err, "Failed to create destination file")
	defer destination.Close()

	_, err = io.Copy(destination, source)
	fatal(err, "Failed to copy file")
}

func copyDir(src, dst string) {
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		fatal(err, "Failed to walk directory")
		targetPath := filepath.Join(dst, strings.TrimPrefix(path, src))
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		copyFile(path, targetPath)
		return nil
	})
	fatal(err, "Failed to copy directory")
}

func load() []Article {
	var articles []Article
	dir, err := ioutil.ReadDir(filepath.Join(wd, "articles"))
	fatal(err, "Failed to read articles directory")

	for _, f := range dir {
		articlePath := filepath.Join(wd, "articles", f.Name())
		content, err := ioutil.ReadFile(articlePath)
		fatal(err, "Failed to read article file")

		parts := strings.SplitN(string(content), "\n", 2)
		if len(parts) < 2 {
			fatal(fmt.Errorf("article must have a title and body"), "Invalid article format")
		}

		title, body := preProcess(parts[0], parts[1])

		cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", articlePath)
		updatedOn, err := cmd.Output()
		fatal(err, "Failed to get last updated date")

		articles = append(articles, Article{
			ID:        strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
			Title:     title,
			UpdatedOn: strings.TrimSpace(string(updatedOn)),
			Body:      body,
		})
	}

	return articles
}

func buildArticle(articleID string) {
	article, err := loadArticle(articleID)
	if err != nil {
		fmt.Printf("Article not found: %s\n", articleID)
		return
	}

	articleDir := filepath.Join(wd, "public", article.ID)
	fatal(os.MkdirAll(articleDir, os.ModePerm), "Failed to create article directory")

	f, err := os.Create(filepath.Join(articleDir, "index.html"))
	fatal(err, "Failed to create article index.html")

	page := template.Must(template.ParseFiles(filepath.Join(wd, "theme", "article.html")))
	fatal(page.Execute(f, struct{ Article Article }{article}), "Failed to execute article template")
}

func loadArticle(articleID string) (Article, error) {
	articlePath := filepath.Join(wd, "articles", articleID+".md")
	content, err := ioutil.ReadFile(articlePath)
	if err != nil {
		return Article{}, err
	}

	parts := strings.SplitN(string(content), "\n", 2)
	if len(parts) < 2 {
		return Article{}, fmt.Errorf("article must have a title and body")
	}

	title, body := preProcess(parts[0], parts[1])

	cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", articlePath)
	updatedOn, err := cmd.Output()
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:        articleID,
		Title:     title,
		UpdatedOn: strings.TrimSpace(string(updatedOn)),
		Body:      body,
	}, nil
}

func preProcess(title, body string) (string, template.HTML) {
	if strings.HasPrefix(title, "# ") {
		title = strings.TrimPrefix(title, "# ")
	}

	ext := parser.CommonExtensions | parser.AutoHeadingIDs
	htmlBody := markdown.ToHTML(
		[]byte(body),
		parser.NewWithExtensions(ext),
		markdownHTML.NewRenderer(markdownHTML.RendererOptions{
			AbsolutePrefix: blogURL,
			RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
				if codeBlock, ok := node.(*ast.CodeBlock); ok {
					syntaxHighlight(w, string(codeBlock.Literal), string(codeBlock.Info))
					return ast.GoToNext, true
				}
				return ast.GoToNext, false
			},
		}),
	)

	return title, template.HTML(htmlBody)
}

func syntaxHighlight(w io.Writer, source, lang string) {
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Analyse(source)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	formatter := chromaHTML.New(chromaHTML.WithClasses(true))
	style := styles.Fallback

	iterator, err := lexer.Tokenise(nil, source)
	fatal(err, "Failed to tokenise source for syntax highlighting")

	err = formatter.Format(w, style, iterator)
	fatal(err, "Failed to format syntax highlighting")
}
