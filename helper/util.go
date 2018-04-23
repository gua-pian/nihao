package helper

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type (
	H          map[string]interface{}
	middlefunc func(handler http.Handler) http.Handler
)

type ApiManager struct {
	routers map[string]string
	middle  map[string][]middlefunc
	final   map[string]http.HandlerFunc
}

var am *ApiManager

func init() {
	am = &ApiManager{
		routers: make(map[string]string),
		middle:  make(map[string][]middlefunc),
		final:   make(map[string]http.HandlerFunc),
	}
	am.middle["/api_add/list"] = []middlefunc{Log, Lol}
	am.final["/api_add/list"] = ShowHandler
}

func ContentTypeHandler(r *http.Request) (url.Values, error) {
	u := url.Values{}
	header := r.Header

	// Handle application/json.
	if header.Get("Content-Type") == "application/json" {
		defer r.Body.Close()
		m := make(map[string]string)
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&m)
		for k, v := range m {
			u.Set(k, v)
		}
		return u, nil
	}

	// Handle x-www-form-urlencoded.
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	form := r.PostForm
	return form, nil
}

func SetResponse(w http.ResponseWriter, h H) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(h)
}
