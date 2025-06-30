package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <markdown-file>")
		fmt.Println("Example: go run . file.txt")
		return
	}

	filename := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatal("File does not exist:", filename)
	}

	// Parse markdown file and convert to HTML
	html := ParseMarkdownFile(filename)
	fmt.Println(html)
}
