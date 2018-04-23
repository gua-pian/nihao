package main

import (
	"awesomeProject/helper"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api_add/api", helper.AddApi)

	h := helper.Construct("/api_add/list")
	mux.Handle("/api_add/list", h)

	// Handle the 404 page.
	mux.HandleFunc("/", helper.HomeHandler)

	http.ListenAndServe(":8000", mux)

}
