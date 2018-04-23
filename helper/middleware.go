package helper

import (
	"fmt"
	"net/http"
)

// All the middleware defines here.

func Log(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Log before call h.")
			h.ServeHTTP(w, r)
			fmt.Println("Log after call h.")
		})
}

func Lol(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Lol before call h.")
		h.ServeHTTP(w, r)
		fmt.Println("Lol after call h.")
	})
}

func Construct(name string) http.Handler {
	// Construct from the final method.
	finalMethod := am.final[name]
	handler := http.HandlerFunc(finalMethod)

	if len(am.middle[name]) == 0 {
		return handler
	}

	// Construct with the middle funcs.
	var h http.Handler = am.middle[name][0](handler)
	for i := 1; i < len(am.middle[name]); i++ {
		v := am.middle[name][i]
		h = v(h)
	}
	return h
}
