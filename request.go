package jump

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.59 Safari/537.36"

func NewRequest() *Request {
	jar, _ := cookiejar.New(nil)
	return &Request{
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		client: &http.Client{
			Timeout: time.Second * 30,
			Jar:     jar,
		},
		Jar: jar,
	}
}

type Request struct {
	client  *http.Client
	Headers map[string]string
	Jar     *cookiejar.Jar
	req     http.Request
}

func (r *Request) Do(method, uri string, headers map[string]string, body io.Reader) (*http.Response, []byte, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, nil, err
	}
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	return resp, respBody, err
}

func (r *Request) Get(uri string) (*http.Response, []byte, error) {
	return r.Do("GET", uri, nil, nil)
}

func (r *Request) Post(uri string, headers map[string]string, body io.Reader) (*http.Response, []byte, error) {
	return r.Do("POST", uri, headers, body)
}

func (r *Request) PostJSON(uri string, data map[string]interface{}) (*http.Response, []byte, error) {
	jsonValue, _ := json.Marshal(data)
	return r.Do("POST", uri, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}, bytes.NewBuffer(jsonValue))
}

func (r *Request) PostForm(uri string, data map[string]string) (*http.Response, []byte, error) {
	form := url.Values{}
	for key, value := range data {
		form.Add(key, value)
	}
	return r.Do("POST", uri, map[string]string{
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
	}, strings.NewReader(form.Encode()))
}
