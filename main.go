package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/fs"
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
)

var blogURL, wd, cssPath string

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	var err error
	wd, err = os.Getwd()
	fatal(err, "Failed to get working directory")

	switch os.Args[1] {
	case "serve":
		blogURL = "http://localhost:2000"
		cssPath = "/css/site.css" // Use non-fingerprinted path for development
		fmt.Println("Serving at http://localhost:2000")
		serve(":2000")
	case "build":
		blogURL = "https://dancroak.com"
		build("public")
		fmt.Println("Built at ./public")
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:\n  blog serve\n  blog build")
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

type TemplateData struct {
	Article Article
	CSSPath string
}

func serve(addr string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Normalize the path
		path := r.URL.Path
		if path == "/" {
			http.ServeFile(w, r, filepath.Join(wd, "theme", "index.html"))
			fmt.Printf("%7.1fms %s %s\n", float64(time.Since(startTime))/float64(time.Millisecond), r.Method, path)
			return
		}

		// Serve static files
		if strings.HasPrefix(path, "/images/") {
			fs := http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(wd, "theme", "images"))))
			fs.ServeHTTP(w, r)
			fmt.Printf("%7.1fms %s %s\n", float64(time.Since(startTime))/float64(time.Millisecond), r.Method, path)
			return
		}

		// Serve CSS files with correct MIME type
		if strings.HasPrefix(path, "/css/") {
			if strings.HasSuffix(path, ".css") {
				w.Header().Set("Content-Type", "text/css")
			}
			fs := http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(wd, "theme", "css"))))
			fs.ServeHTTP(w, r)
			fmt.Printf("%7.1fms %s %s\n", float64(time.Since(startTime))/float64(time.Millisecond), r.Method, path)
			return
		}

		// Build and serve the article for non-root paths
		articleID := strings.TrimPrefix(path, "/")
		buildArticle(articleID)

		articleFilePath := filepath.Join(wd, "public", articleID, "index.html")
		if _, err := os.Stat(articleFilePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			fmt.Printf("%7.1fms %s %s (not found)\n", float64(time.Since(startTime))/float64(time.Millisecond), r.Method, path)
			return
		}

		http.ServeFile(w, r, articleFilePath)
		fmt.Printf("%7.1fms %s %s\n", float64(time.Since(startTime))/float64(time.Millisecond), r.Method, path)
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

	// Copy and fingerprint CSS files
	cssPath = fingerprintCSS(outputDir)

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
			fatal(page.Execute(f, TemplateData{Article: a, CSSPath: cssPath}), "Failed to execute article template")
		}(a)
	}
	wg.Wait()
}

func copyFile(srcPath, dstPath string) error {
	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	source, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
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
	articlesDir := filepath.Join(wd, "articles")
	err := filepath.WalkDir(articlesDir, func(path string, d fs.DirEntry, err error) error {
		fatal(err, "Failed to walk articles directory")
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			// Generate article ID relative to articles dir, without extension
			relPath, err := filepath.Rel(articlesDir, path)
			fatal(err, "Failed to get relative path")
			id := strings.TrimSuffix(relPath, filepath.Ext(relPath))
			title, body := preProcess(path)

			// Get last updated date
			cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", path)
			updatedOn, err := cmd.Output()
			fatal(err, "Failed to get last updated date")

			articles = append(articles, Article{
				ID:        filepath.ToSlash(id),
				Title:     title,
				UpdatedOn: strings.TrimSpace(string(updatedOn)),
				Body:      body,
			})
		}
		return nil // Continue walking
	})
	fatal(err, "Failed to walk articles directory")

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
	fatal(page.Execute(f, TemplateData{Article: article, CSSPath: cssPath}), "Failed to execute article template")
}

func loadArticle(articleID string) (Article, error) {
	articlePath := filepath.Join(wd, "articles", articleID+".md")
	if _, err := os.Stat(articlePath); os.IsNotExist(err) {
		return Article{}, fmt.Errorf("article not found")
	}
	title, body := preProcess(articlePath)

	// Get last updated date
	cmd := exec.Command("git", "log", "-1", "--format=%cd", "--date=format:%B %d, %Y", "--", articlePath)
	updatedOn, err := cmd.Output()
	if err != nil {
		return Article{}, err
	}

	return Article{
		ID:        filepath.ToSlash(articleID),
		Title:     title,
		UpdatedOn: strings.TrimSpace(string(updatedOn)),
		Body:      body,
	}, nil
}

func preProcess(filePath string) (string, template.HTML) {
	f, err := os.Open(filePath)
	if err != nil {
		fatal(err, "Failed to open file")
	}
	defer f.Close()

	var (
		scanner = bufio.NewScanner(f)
		isFirst = true
		title   string
		lines   []string // Collect lines in a slice
	)

	for scanner.Scan() {
		line := scanner.Text()

		if isFirst {
			if !strings.HasPrefix(line, "# ") {
				fatal(fmt.Errorf("first line must be an h1 like: # Intro"), "Invalid first line")
			}

			title = line[2:]
			isFirst = false
			continue
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fatal(err, "Error reading file")
	}

	body := strings.Join(lines, "\n")

	// Render the markdown to HTML with syntax highlighting
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

// fingerprintCSS copies CSS files with MD5 fingerprints and returns the fingerprinted path
func fingerprintCSS(outputDir string) string {
	srcPath := filepath.Join(wd, "theme", "css", "site.css")

	content, err := os.ReadFile(srcPath)
	fatal(err, "Failed to read CSS file")

	hash := fmt.Sprintf("%x", md5.Sum(content))

	cssDir := filepath.Join(outputDir, "css")
	fatal(os.MkdirAll(cssDir, os.ModePerm), "Failed to create CSS directory")

	copyDir(filepath.Join(wd, "theme", "css"), cssDir)

	fpName := fmt.Sprintf("site-%s.css", hash[:8])
	fpPath := filepath.Join(cssDir, fpName)

	fatal(os.WriteFile(fpPath, content, 0644), "Failed to write fingerprinted CSS")

	return "/css/" + fpName
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
