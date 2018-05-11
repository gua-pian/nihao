package main

import (
	"awesomeProject/helper"
	"log"
	"net/http"
)

var version string = "v1.0.0"

func main() {

	log.Printf("[INFO] Version %s Starting\n", version)
	mux := http.NewServeMux()

	mux.HandleFunc("/api_add/api", helper.AddApi)

	h := helper.Construct("/api_add/list")
	mux.Handle("/api_add/list", h)

	// Handle the 404 page.
	mux.HandleFunc("/", helper.HomeHandler)

	http.ListenAndServe(":8000", mux)

}
