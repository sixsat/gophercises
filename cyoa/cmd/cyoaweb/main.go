package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sixsat/gophercises/cyoa"
)

func main() {
	port := flag.Int("port", 3000, "the port to start the CYOA web application on")
	jsonFile := flag.String("json", "gopher.json", "the json file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *jsonFile)

	// Open the JSON file and parse the story in it.
	file, err := os.Open(*jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	story, err := cyoa.JsonToStory(file)
	if err != nil {
		log.Fatal(err)
	}

	// Create custom CYOA story handler.
	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
