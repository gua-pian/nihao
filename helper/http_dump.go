package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

type RequestInfo struct {
	ClientAddress  string `json:"client_address"`
	RequestUri     string `json:"request_uri"`
	BodyParameters url.Values `json:"body_parameters"`
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func ContentTypeHandler(r *http.Request) (url.Values, error) {
	u := url.Values{}
	header := r.Header

	var err error
	save := r.Body
	save, r.Body, err = drainBody(r.Body)
	if err != nil {
		return nil, err
	}

	// Handle application/json.
	if header.Get("Content-Type") == "application/json" {
		defer r.Body.Close()
		m := make(map[string]string)
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&m)
		for k, v := range m {
			u.Set(k, v)
		}
		r.Body = save
		return u, nil
	}

	// Handle x-www-form-urlencoded.
	err = r.ParseForm()
	if err != nil {
		return nil, err
	}

	form := r.PostForm

	return form, nil
}

func RemoteIp(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		remoteAddr = ip
	} else if ip = r.Header.Get("X-Forwarded-For"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

func dumpRequest(r *http.Request) error {
	// Log the parameters in body.
	values, err := ContentTypeHandler(r)
	if err != nil {
		return err
	}

	// Log the parameters of RemoteAddr, url.RawQuery and parameters.
	requestBody := RequestInfo{RemoteIp(r), r.RequestURI, values}
	bytesRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("[INFO] request_body:%+q\n", bytesRequestBody)
	return nil
}
