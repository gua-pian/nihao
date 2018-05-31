package main

import (
	"github.com/gua-pian/nihao/helper"
	"log"
	"net/http"
)

var version string = "v1.0.0"

func main() {

	log.Printf("[INFO] Version %s Starting\n", version)
	mux := http.NewServeMux()

	// Register apis to the gway.
	mux.HandleFunc("/api_add/api", helper.AddApi)

	// Show all the apis registerd to the gway.
	h := helper.Construct("/api_add/list")
	mux.Handle("/api_add/list", h)

	// Handle the 404 page.
	mux.HandleFunc("/", helper.HomeHandler)

	http.ListenAndServe(":8000", mux)

}
