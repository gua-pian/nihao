package helper

import (
	"net/http"
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

	method := values.Get("method")
	handler := values.Get("handler")

	if _, ok := am.routers[method]; ok {
		SetResponse(w, H{"Status": -2, "message": "Method added already"})
		return
	}

	if method == "" || handler == "" {
		SetResponse(w, H{"Status": -3, "Info": "paramater error"})
		return
	}
	am.routers[method] = handler

	SetResponse(w, H{"Status": 0, "Info": "method added ok!"})
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	SetResponse(w, H{"Status": 0, "Data": am.routers})
	return
}

var HomeHandler = func(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		SetResponse(w, H{"Status": -1, "Info": "You are calling the wrong api."})
		return
	}
	SetResponse(w, H{"Status": 0, "Info": "This is home."})
}
