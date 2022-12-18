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
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"
	"time"

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
var showScheduled = true

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
		showScheduled = false
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

// Article contains data loaded from config.json and parsed Markdown
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
	articles := load()

	noDashes := strings.Replace(id, "-", " ", -1)
	noUnderscores := strings.Replace(noDashes, "_", " ", -1)
	title := strings.Title(noUnderscores)
	content := []byte("# " + title + "\n\n\n")
	check(ioutil.WriteFile(wd+"/articles/"+id+".md", content, 0644))

	a := Article{
		ID:      id,
		Updated: time.Now().Format("2006-01-02"),
	}

	articles = append([]Article{a}, articles...)
	config, err := json.MarshalIndent(articles, "", "  ")
	check(err)
	check(ioutil.WriteFile(wd+"/config.json", config, 0644))
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
	articles := load()

	// public directories
	dir, err := ioutil.ReadDir(wd + "/public")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"public", d.Name()}...))
	}
	check(os.MkdirAll(wd+"/public/images", os.ModePerm))

	// index page
	indexPage := template.Must(template.ParseFiles(wd + "/theme/index.html"))
	f, err := os.Create("public/index.html")
	check(err)
	indexData := struct {
		Articles []Article
	}{
		Articles: articles,
	}
	check(indexPage.Execute(f, indexData))

	// article pages
	articlePage := template.Must(template.ParseFiles(wd + "/theme/article.html"))
	for _, a := range articles {
		check(os.Mkdir(wd+"/public/"+a.ID, os.ModePerm))
		f, err := os.Create(wd + "/public/" + a.ID + "/index.html")
		check(err)
		articleData := struct {
			Article Article
		}{
			Article: a,
		}
		check(articlePage.Execute(f, articleData))
	}

	// images
	cmd := exec.Command("cp", "-a", wd+"/images/.", wd+"/public/images")
	cmd.Run()

	// favicon.ico, and additional files from theme
	cmd = exec.Command("cp", "-a", wd+"/theme/public/.", wd+"/public")
	cmd.Run()
}

func load() []Article {
	config, err := ioutil.ReadFile(wd + "/config.json")
	check(err)
	var articles []Article
	check(json.Unmarshal(config, &articles))

	for i, a := range articles {
		t, err := time.Parse("2006-01-02", a.Updated)
		check(err)

		now := time.Now()
		if showScheduled == false && t.After(now) {
			continue
		}

		title, body := preProcess("articles/" + a.ID + ".md")
		ext := parser.CommonExtensions | parser.AutoHeadingIDs
		html := markdown.ToHTML(
			[]byte(body),
			parser.NewWithExtensions(ext),
			html.NewRenderer(html.RendererOptions{
				AbsolutePrefix: blogURL,
				RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
					codeBlock, ok := node.(*ast.CodeBlock)
					if !ok {
						return ast.GoToNext, false
					}
					lang := string(codeBlock.Info)
					syntaxHighlight(w, string(codeBlock.Literal), lang)
					return ast.GoToNext, true
				},
			}), // RenderNodeHook
		)

		a := Article{
			Body:        template.HTML(html),
			Canonical:   a.Canonical,
			Description: a.Description,
			ID:          a.ID,
			Updated:     a.Updated,
			UpdatedOn:   t.Format("January 2, 2006"),
			Title:       title,
		}
		articles[i] = a
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
				begindoc := strings.Index(string(srcCode), sep)
				if begindoc == -1 {
					exitWith("error: embed separator not found " + sep + " in " + filename)
				}
				// end of comment line
				begindoc += len(sep)

				sep = "enddoc: " + id
				enddoc := strings.Index(string(srcCode), sep)
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
