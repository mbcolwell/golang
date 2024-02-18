package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mbcolwell/golang/internal/cyoa"
)

func main() {
	filename := flag.String("filename", "gopher.json", "The file containing the story to explore")
	port := flag.Int("port", 3000, "The port to access the CYOA web app on")
	flag.Parse()

	story, err := cyoa.ParseJson(*filename)
	if err != nil {
		panic(err)
	}

	h := cyoa.StoryHandler(story)
	fmt.Printf("Starting the server on port: %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
