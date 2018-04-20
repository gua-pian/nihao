package helper


import (
	"net/http"
	"fmt"
)

var AddApi = func (w http.ResponseWriter, r *http.Request){
	if (r.Method != "POST"){
		SetResponse(w, H{"message", "Method not right"})
		return
	}

	// Handler content-type.
	values, _ := ContentTypeHandler(r)

	fmt.Println(values)

	method := values.Get("method")
	handler := values.Get("handler")

	if _, ok := am.routers[r.FormValue("method")]; ok {
		SetResponse(w, H{"message", "Method added already"})
		return
	}

	if method == "" || handler == "" {
		SetResponse(w, H{"message", "paramater error"})
		return
	}
	am.routers[method] = handler

	SetResponse(w, H{"message", "method added ok!"})
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("haha")
	SetResponse(w, H{"data", am.routers})
	return
}