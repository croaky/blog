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
