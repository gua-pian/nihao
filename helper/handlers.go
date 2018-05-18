package helper

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var AddApi = func(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SetResponse(w, H{"Status": -1, "Info": "Method not right."})
		return
	}

	// Handler content-type.
	values, err := ContentTypeHandler(r)
	if err != nil {
		SetResponse(w, H{"Status": -1, "Info": err})
		return
	}

	apiName := values.Get("apiName")
	apiUpstream := values.Get("apiUpstream")

	if _, ok := am.routers[apiName]; ok {
		SetResponse(w, H{"Status": -1, "message": "Method added already"})
		return
	}

	if apiName == "" || apiUpstream == "" {
		SetResponse(w, H{"Status": -1, "Info": "paramater error"})
		return
	}
	am.routers[apiName] = apiUpstream
	SetResponse(w, H{"Status": 0, "Info": "method added ok!"})
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	SetResponse(w, H{"Status": 0, "Data": am.routers})
	return
}

var HomeHandler = func(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("url: %s\n", r.URL.Path)
	path := r.URL.Path
	router, ok := am.routers[path]
	if !ok {
		if path != "/" {
			fmt.Printf("path: %s\n", path)
			SetResponse(w, H{"Status": -1, "Info": "You are calling the wrong api."})
			return
		}
		SetResponse(w, H{"Status": 0, "Info": "This is home."})
		return
	}

	// Init the request with uuid.

	// Log all the request parameters.
	err := dumpRequest(r)
	if err != nil {
		SetResponse(w, H{"Status": -1, "Info": err})
		return
	}

	u, err := url.Parse(router + path)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	h := newHttpProxy(u)
	begin := time.Now()
	h.ServeHTTP(w, r)
	duration := time.Since(begin).String()
	fmt.Println("time elapsed:" + duration)
}
