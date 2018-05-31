package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
	"fmt"
)

type (
	RequestInfo struct {
		TimeStamp       string     `json:"time"`
		ClientAddress   string     `json:"client_address"`
		UpstreamAddress string     `json:"upstream_address"`
		UpstreamUri     string     `json:"upstream_uri"`
		BodyParameters  url.Values `json:"body_parameters"`
	}

	ResponseInfo struct {
		TimeStamp     string          `json:"time"`
		StatusCode    int             `json:"status_code"`
		ContentLength int64           `json:"content_length"`
		ResponseBody  json.RawMessage `json:"response_body"`
	}
)

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
	header := r.Header
	var err error
	var save io.ReadCloser
	save, r.Body, err = drainBody(r.Body)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	// Handle application/json.
	if header.Get("Content-Type") == "application/json" {

		u := url.Values{}
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
	r.Body = save
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

func dumpRequest(r *http.Request, router string) error {
	// Log the parameters in body.
	values, err := ContentTypeHandler(r)
	// fmt.Println(values)
	if err != nil {
		return err
	}

	// Log the parameters of RemoteAddr, url.RawQuery and parameters.
	requestBody := RequestInfo{TimeStamp: time.Now().Format(timeFormat), ClientAddress: RemoteIp(r), UpstreamAddress: router, UpstreamUri: r.RequestURI, BodyParameters: values}
	bytesRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		dumpLogFile.WriteString("Error When Marshal Resquest Body: " + err.Error())
		fmt.Fprintln(dumpLogFile)
		return err
	}
	dumpLogFile.WriteString("request_body:")
	dumpLogFile.Write(bytesRequestBody)
	fmt.Fprintln(dumpLogFile)
	// fmt.Fprintln(dumpLogFile)
	return nil
}

func dumpReponse(res *http.Response) error {
	var err error
	var save io.ReadCloser
	save, res.Body, err = drainBody(res.Body)
	if err != nil {
		return err
	}
	// dump the res.Body and write to log.
	responseBody := ResponseInfo{TimeStamp: time.Now().Format(timeFormat), StatusCode: res.StatusCode, ContentLength: res.ContentLength}
	bytesBody, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		dumpLogFile.WriteString("Error When Read Response Body:" + err.Error())
		fmt.Fprintln(dumpLogFile)
		return err
	}
	responseBody.ResponseBody = json.RawMessage(bytesBody)
	// fmt.Printf("%+v\n",responseBody)
	bytesResponseBody, err := json.Marshal(responseBody)
	if err != nil {
		dumpLogFile.WriteString("Error When Marshal Response Body:" + err.Error())
		fmt.Fprintln(dumpLogFile)
		return err
	}
	dumpLogFile.WriteString("response_body:")
	dumpLogFile.Write(bytesResponseBody)
	fmt.Fprintln(dumpLogFile)
	// fmt.Fprintln(dumpLogFile)
	// Restore res.Body.
	res.Body = save
	return nil
}
