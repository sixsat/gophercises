package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sixsat/gophercises/link"
)

func main() {
	htmlFile := flag.String("file", "examples/ex1.html", "an HTML file to be parsed")
	flag.Parse()

	f, err := os.Open(*htmlFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	links, err := link.Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", links)
}
