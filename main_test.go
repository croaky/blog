package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestBuild(t *testing.T) {
	// Set up temporary directories for the test
	tempOutputDir, err := ioutil.TempDir("", "site_test_output")
	if err != nil {
		t.Fatalf("Failed to create temp output directory: %v", err)
	}
	defer os.RemoveAll(tempOutputDir)

	// Run the build process with the temporary output directory
	build(tempOutputDir)

	// Check that the index page was created
	indexFilePath := filepath.Join(tempOutputDir, "index.html")
	if _, err := os.Stat(indexFilePath); os.IsNotExist(err) {
		t.Fatalf("Index page was not created")
	}

	// Parse the HTML and check for parsing errors
	doc := parseHTML(t, indexFilePath)

	// Check that the index page has the correct title
	checkHTMLTitle(t, doc, "Dan Croak")

	// Clean up the test environment
	os.RemoveAll(tempOutputDir)
}

func parseHTML(t *testing.T, filePath string) *html.Node {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	return doc
}

func checkHTMLTitle(t *testing.T, doc *html.Node, expectedTitle string) {
	var title string
	var inHead bool
	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.ElementNode {
			if node.Data == "head" {
				inHead = true
			} else if node.Data == "title" && inHead && node.FirstChild != nil {
				title = node.FirstChild.Data
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			traverse(child)
		}
		if node.Type == html.ElementNode && node.Data == "head" {
			inHead = false
		}
	}
	traverse(doc)

	if title != expectedTitle {
		t.Errorf("Expected title '%s', got '%s'", expectedTitle, title)
	}
}

func TestPreProcessEmbed(t *testing.T) {
	// Set the working directory to the directory containing the test files
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Path to the test article and code files
	articlePath := filepath.Join(wd, "articles", "postgres-tips.md")

	// Run the preProcess function on the test article
	_, body := preProcess(articlePath)

	// Check that the body contains the expected embedded code lines
	if !strings.Contains(string(body), "pg_dump") {
		t.Errorf("Expected body to contain embedded code line %q, got: %s", "pg_dump", body)
	}
	if !strings.Contains(string(body), "pg_restore") {
		t.Errorf("Expected body to contain embedded code line %q, got: %s", "pg_restore", body)
	}
}
