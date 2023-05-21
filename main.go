/*
Add an article:

	blog add <article-url-slug>

Serve site on localhost:

	blog serve

Build site (HTML, images, code, Netlify files) to `public/`:

	blog build
*/
package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/alecthomas/chroma/v2"
	htmlfmt "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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
	Canonical   string `json:"canonical,omitempty"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Updated     string `json:"updated"`

	Body      template.HTML `json:"-"`
	Title     string        `json:"-"`
	UpdatedOn string        `json:"-"`
}

func add(id string) {
	noDashes := strings.Replace(id, "-", " ", -1)
	noUnderscores := strings.Replace(noDashes, "_", " ", -1)
	title := strings.Title(noUnderscores)
	content := []byte("# " + title + "\n\n\n")
	check(ioutil.WriteFile(wd+"/articles/"+id+".md", content, 0644))
}

func serve(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// log every request except favicon.ico noise
		if r.URL.Path != "/favicon.ico" {
			fmt.Println(r.Method + " " + r.URL.Path)
		}

		// don't rebuild for images or favicon
		if !strings.HasPrefix(r.URL.Path, "/images/") && !strings.HasPrefix(r.URL.Path, "/favicon.ico") {
			build()
		}

		fs := http.FileServer(http.Dir(wd + "/public"))
		fs.ServeHTTP(w, r)
	})
	check(http.ListenAndServe(addr, nil))
}

func build() {
	// clean public dir
	check(os.MkdirAll(wd+"/public/", os.ModePerm))
	dir, err := ioutil.ReadDir(wd + "/public")
	check(err)
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"public", d.Name()}...))
	}

	// build article pages
	page := template.Must(template.ParseFiles(wd + "/theme/article.html"))
	articles := load()

	for _, a := range articles {
		check(os.MkdirAll(wd+"/public/"+a.ID, os.ModePerm))
		f, err := os.Create(wd + "/public/" + a.ID + "/index.html")
		check(err)
		data := struct {
			Article Article
		}{
			Article: a,
		}
		check(page.Execute(f, data))
	}

	// copy index page
	exec.Command("cp", "-a", wd+"/theme/index.html", wd+"/public/").Run()

	// copy static assets
	check(os.MkdirAll(wd+"/public/images", os.ModePerm))
	exec.Command("cp", "-a", wd+"/images/.", wd+"/public/images").Run()
	exec.Command("cp", "-a", wd+"/theme/public/.", wd+"/public").Run()
}

func load() []Article {
	var articles []Article
	dir, err := ioutil.ReadDir(wd + "/articles")
	check(err)

	for _, f := range dir {
		title, body := preProcess("articles/" + f.Name())
		ext := parser.CommonExtensions | parser.AutoHeadingIDs
		html := markdown.ToHTML(
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

		// use the accurate Git date for CI / CD
		cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", "articles/"+f.Name())
		updatedOn, err := cmd.Output()
		check(err)

		a := Article{
			Body:      template.HTML(html),
			ID:        strings.TrimSuffix(f.Name(), filepath.Ext(f.Name())),
			UpdatedOn: string(updatedOn),
			Title:     title,
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
func preProcess(filepath string) (title, body string) {
	f, err := os.Open(filepath)
	check(err)
	defer f.Close()

	var (
		scanner = bufio.NewScanner(f)
		isFirst = true
		isEmbed = false
	)

	for scanner.Scan() {
		line := scanner.Text()

		if isFirst {
			if strings.Index(line, "# ") == -1 {
				exitWith("error: first line must be an h1 like: # Intro")
			}

			title = line[2:]
			isFirst = false
			continue
		}

		if line == "```embed" {
			isEmbed = true
			continue
		}

		if isEmbed {
			parts := strings.Split(line, " ")
			if len(parts) != 1 && len(parts) != 2 {
				exitWith("error: embed line must be filepath id (code/test.rb id) or filepath (code/test.rb)")
			}

			filename := parts[0]
			srcCode, err := ioutil.ReadFile(wd + "/" + filename)
			check(err)

			begindoc := 0
			enddoc := len(srcCode) - 1

			if len(parts) == 2 {
				id := parts[1]
				sep := "begindoc: " + id + "\n"
				begindoc = strings.Index(string(srcCode), sep)
				if begindoc == -1 {
					exitWith("error: embed separator not found " + sep + " in " + filename)
				}
				// end of comment line
				begindoc += len(sep)

				sep = "enddoc: " + id
				enddoc = strings.Index(string(srcCode), sep)
				if enddoc == -1 {
					exitWith("error: embed separator not found " + sep + " in " + filename)
				}
				// backtrack to last newline to cut out comment character(s)
				enddoc = strings.LastIndex(string(srcCode[0:enddoc]), "\n")
			}

			rawLines := strings.Split(string(srcCode[begindoc:enddoc]), "\n")

			leadingWhitespace := regexp.MustCompile("(?m)(^[ \t]*)(?:[^ \t])")
			var margin string
			var lines []string

			for i, l := range rawLines {
				if i == 0 {
					margin = leadingWhitespace.FindAllStringSubmatch(l, -1)[0][1]
				}
				dedented := regexp.MustCompile("(?m)^"+margin).ReplaceAllString(l, "")
				lines = append(lines, dedented)
			}

			ext := strings.Trim(path.Ext(filename), ".")
			body += "```" + ext + "\n" + strings.Join(lines, "\n")

			isEmbed = false
			continue
		}

		body += "\n" + line
	}

	return title, body
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
