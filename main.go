package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tracker/handlers"
)

func main() {

args := os.Args
if len(args) > 1 {
	fmt.Println("Error: Usage: go run .")
	return
}

	http.HandleFunc("/", handlers.HomepageHandler)
	http.HandleFunc("/artist", handlers.ArtistHandler)
	http.HandleFunc("/dates", handlers.DateHandler)
	http.HandleFunc("/locations", handlers.LocationHandler)
	// serve the static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Print("Starting server at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
