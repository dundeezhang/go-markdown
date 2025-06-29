package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	read()
	content, err := os.ReadFile("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(content))
}
