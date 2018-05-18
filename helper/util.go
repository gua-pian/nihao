package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type (
	H           map[string]interface{}
	middlefunc  func(handler http.Handler) http.Handler
)

/* ApiManager consists of following fields.
	routers: indicates the registered api the upstream it should goes to. Example. login: 10.10.10.10:4000
	middle: all the middleware associated with api. Example. login: [Log, Auth, etc..]
    final: the final Api it should call after all the middleware.
*/
type ApiManager struct {
	routers map[string]string
	middle  map[string][]middlefunc
	final   map[string]http.HandlerFunc
	timeout map[string]int64
}

var am *ApiManager

func init() {
	am = &ApiManager{
		routers: make(map[string]string),
		middle:  make(map[string][]middlefunc),
		final:   make(map[string]http.HandlerFunc),
		timeout: make(map[string]int64),
	}
	am.routers["/user/list_material"] = "http://10.2.1.107:8085"
	am.middle["/api_add/list"] = []middlefunc{Log, Lol}
	am.final["/api_add/list"] = ShowHandler
}

func SetResponse(w http.ResponseWriter, h H) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(h)
}

func newHttpProxy(target *url.URL) http.Handler {

	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.URL.Path = target.Path
			req.URL.RawQuery = target.RawQuery
		},
		ModifyResponse: func(response *http.Response) error {
			body_bytes, err := httputil.DumpResponse(response, true)
			if err != nil {
				fmt.Println("error when dump response")
				fmt.Printf("%s\n", err.Error())
				return err
			}
			fmt.Printf("response body %s\n", string(body_bytes))
			return nil
		},
	}

}
