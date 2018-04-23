package helper

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
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
		SetResponse(w, H{"Status": -1, "message": "Method added already"})
		return
	}

	if method == "" || handler == "" {
		SetResponse(w, H{"Status": -1, "Info": "paramater error"})
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
	path := r.URL.Path
	router, ok := am.routers[path]
	if !ok {
		if path != "/" {
			SetResponse(w, H{"Status": -1, "Info": "You are calling the wrong api."})
			return
		}
		SetResponse(w, H{"Status": 0, "Info": "This is home."})
		return
	}

	// Log all the request parameters.
	// Handler content-type.
	values, err := ContentTypeHandler(r)
	if err != nil {
		SetResponse(w, H{"Status": -1, "Info": err})
		return
	}
	fmt.Printf("%+v\n", values)

	// Forward the request to backend server.
	newRequest := ConstructHttpRequest(r, router+path)
	fmt.Println(r.URL)
	fmt.Println(r.Host)

	client := http.Client{}
	res, err := client.Do(newRequest)
	if err != nil {
		SetResponse(w, H{"Status": -1, "Info": err})
		return
	}
	defer res.Body.Close()

	resBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		SetResponse(w, H{"Status": -1, "Info": err})
		return
	}

	resReader4Log := bytes.NewReader(resBytes)
	resReader4Return := bytes.NewReader(resBytes)

	// Log the response.
	res4Log := ioutil.NopCloser(resReader4Log)
	loginfo, _ := ioutil.ReadAll(res4Log)
	fmt.Printf("%s\n", loginfo)

	// Send the response to downstream.
	res4Return := ioutil.NopCloser(resReader4Return)
	s, _ := ioutil.ReadAll(res4Return)
	w.Header().Set("Content-Type","application/json")
	w.Write([]byte(s))
}
