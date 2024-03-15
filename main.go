/*
Add an article:

	blog add <article-url-slug>

Serve site on localhost:

	blog serve

Build site (HTML, images, code) to `public/`:

	blog build
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
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alecthomas/chroma/v2"
	htmlfmt "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var blogURL = "https://dancroak.com"
var wd string

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	wd, err = os.Getwd()
	check(err)

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			usage()
		}
		id := os.Args[2]
		add(id)
		fmt.Println("Added ./articles/" + id + ".md")
	case "serve":
		fmt.Println("Serving at http://localhost:2000")
		serve(":2000")
	case "build":
		build()
		fmt.Println("Built at ./public")
	default:
		usage()
	}
}

func usage() {
	const s = `usage:
  blog add <article-url-slug>
  blog serve
  blog build
`
	fmt.Fprint(os.Stderr, s)
	os.Exit(2)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		_, file, no, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("%s#%d\n", file, no)
		}
		os.Exit(1)
	}
}

func exitWith(s string) {
	fmt.Println(s)
	os.Exit(1)
}

// Article contains data loaded from articles/*.md
type Article struct {
	ID        string
	Title     string
	UpdatedOn string
	Body      template.HTML
}

func add(id string) {
	noDashes := strings.ReplaceAll(id, "-", " ")
	noUnderscores := strings.ReplaceAll(noDashes, "_", " ")
	c := cases.Title(language.Und)
	title := c.String(noUnderscores)
	content := []byte("# " + title + "\n\n\n")
	check(os.WriteFile(wd+"/articles/"+id+".md", content, 0644))
}

func serve(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Trim trailing slash for consistency
		path := strings.TrimSuffix(r.URL.Path, "/")

		// Don't rebuild for favicon or images.
		if path != "" && path != "/favicon.ico" && !strings.HasPrefix(path, "/images") {
			buildArticle(path)
		} else {
			buildIndex()
		}

		fs := http.FileServer(http.Dir(wd + "/public"))
		fs.ServeHTTP(w, r)

		if path == "" {
			path = "/"
		}
		duration := time.Since(startTime)
		fmt.Printf("%7.1f ms %s %s\n", float64(duration)/float64(time.Millisecond), r.Method, path)
	})
	check(http.ListenAndServe(addr, nil))
}

func buildArticle(articleID string) {
	page := template.Must(template.ParseFiles(wd + "/theme/article.html"))

	// Load the specific article
	article, err := loadArticle(articleID)
	if err != nil {
		fmt.Println("Article not found:", articleID)
		return
	}

	// Build the article page
	check(os.MkdirAll(wd+"/public/"+article.ID, os.ModePerm))
	f, err := os.Create(wd + "/public/" + article.ID + "/index.html")
	check(err)
	data := struct {
		Article Article
	}{
		Article: article,
	}
	check(page.Execute(f, data))
}

func loadArticle(articleID string) (Article, error) {
	articlePath := wd + "/articles/" + articleID + ".md"
	content, err := ioutil.ReadFile(articlePath)
	if err != nil {
		return Article{}, err
	}

	// Split the content into title and body
	parts := strings.SplitN(string(content), "\n", 2)
	if len(parts) < 2 {
		return Article{}, fmt.Errorf("error: article must have a title and body")
	}
	title, body := parts[0], parts[1]

	// Preprocess the article
	title, htmlBody := preProcess(title, body)

	// Get the last updated date using Git
	cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", articlePath)
	updatedOn, err := cmd.Output()
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:        articleID,
		Title:     title,
		UpdatedOn: strings.TrimSpace(string(updatedOn)),
		Body:      htmlBody,
	}, nil
}

func buildIndex() {
	// Load the index page template
	page := template.Must(template.ParseFiles(wd + "/theme/index.html"))

	// Create the index page
	check(os.MkdirAll(wd+"/public", os.ModePerm))
	f, err := os.Create(wd + "/public/index.html")
	check(err)

	// You can pass data to the template if needed, for example, a list of articles
	// For simplicity, we'll pass an empty struct here
	data := struct{}{}

	check(page.Execute(f, data))
}

func build() {
	// clean public dir
	dirEntries, err := os.ReadDir(wd + "/public")
	check(err)
	for _, d := range dirEntries {
		os.RemoveAll(path.Join("public", d.Name()))
	}

	// build article pages
	page := template.Must(template.ParseFiles(wd + "/theme/article.html"))
	articles := load()

	var wg sync.WaitGroup
	for _, a := range articles {
		wg.Add(1)
		go func(a Article) {
			defer wg.Done()
			check(os.MkdirAll(wd+"/public/"+a.ID, os.ModePerm))
			f, err := os.Create(wd + "/public/" + a.ID + "/index.html")
			check(err)
			data := struct {
				Article Article
			}{
				Article: a,
			}
			check(page.Execute(f, data))
		}(a)
	}
	wg.Wait()

	// copy index page
	copyFile(wd+"/theme/index.html", wd+"/public/index.html")

	// copy static assets
	copyDir(wd+"/images", wd+"/public/images")
	copyDir(wd+"/theme/public", wd+"/public")
}

func copyFile(src, dst string) {
	source, err := os.Open(src)
	check(err)
	defer source.Close()

	destination, err := os.Create(dst)
	check(err)
	defer destination.Close()

	_, err = io.Copy(destination, source)
	check(err)
}

func copyDir(src, dst string) {
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		check(err)
		relPath := strings.TrimPrefix(path, src)
		targetPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}
		copyFile(path, targetPath)
		return nil
	})
	check(err)
}

func load() []Article {
	var articles []Article
	dir, err := ioutil.ReadDir(wd + "/articles")
	check(err)

	for _, f := range dir {
		articlePath := wd + "/articles/" + f.Name()
		content, err := ioutil.ReadFile(articlePath)
		check(err)

		// Split the content into title and body
		parts := strings.SplitN(string(content), "\n", 2)
		if len(parts) < 2 {
			exitWith("error: article must have a title and body")
		}
		title, body := parts[0], parts[1]

		// Preprocess the article
		title, htmlBody := preProcess(title, body)

		// Get the last updated date using Git
		cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", articlePath)
		updatedOn, err := cmd.Output()
		check(err)

		a := Article{
			ID:        strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
			Title:     title,
			UpdatedOn: strings.TrimSpace(string(updatedOn)),
			Body:      htmlBody,
		}
		articles = append(articles, a)
	}

	return articles
}

/*
preProcess scans the Markdown document at filepath line-by-line,
extracting article title and "pre-processing" the article body
which can then be passed to a Markdown compiler at the call site.

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
*/
func preProcess(title, body string) (string, template.HTML) {
	// Remove the "# " prefix from the title
	if strings.HasPrefix(title, "# ") {
		title = strings.TrimPrefix(title, "# ")
	}

	// Markdown to HTML conversion
	ext := parser.CommonExtensions | parser.AutoHeadingIDs
	htmlBody := markdown.ToHTML(
		[]byte(body),
		parser.NewWithExtensions(ext),
		html.NewRenderer(html.RendererOptions{
			AbsolutePrefix: blogURL,
			RenderNodeHook: func(w io.Writer, node ast.Node, _entering bool) (ast.WalkStatus, bool) {
				codeBlock, ok := node.(*ast.CodeBlock)
				if !ok {
					return ast.GoToNext, false
				}
				lang := string(codeBlock.Info)
				syntaxHighlight(w, string(codeBlock.Literal), lang)
				return ast.GoToNext, true
			},
		}),
	)

	return title, template.HTML(htmlBody)
}

func syntaxHighlight(w io.Writer, source, lang string) error {
	// lexer
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// formatter
	f := htmlfmt.New(htmlfmt.Standalone(false), htmlfmt.WithClasses(true))

	// style
	s := styles.Fallback

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return f.Format(w, s, it)
}
