package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/bryanvaz/go-lucide/test/pages"
)

func main() {

	http.Handle("/", templ.Handler(pages.Index()))

	log.Println("Server starting on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
