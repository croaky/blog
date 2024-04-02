// begindoc: all
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	url := os.Getenv("SLACK_WEBHOOK")
	if url == "" {
		log.Fatalln("no webhook provided")
	}

	reqBody, err := json.Marshal(map[string]string{
		"text": "Hello, world!",
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(respBody))
}

// enddoc: all
