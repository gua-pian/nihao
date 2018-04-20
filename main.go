package main

import (
	"net/http"
	"awesomeProject/helper"
)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/api_add/api", helper.AddApi)

	h := helper.Construct("/api_add/list")
	mux.Handle("/api_add/list", h)

	http.ListenAndServe(":8000", mux)
}