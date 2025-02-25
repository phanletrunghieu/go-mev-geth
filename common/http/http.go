package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Param struct {
	Query  map[string]string
	Header map[string]string
}

func Get(url string, res interface{}, param ...Param) (err error) {
	var body []byte
	if len(param) > 0 {
		body, err = get(url, param[0])
	} else {
		body, err = get(url)
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		_ = errors.New("decode resp to struct err")
	}
	return
}

func Post(url string, payload interface{}, res interface{}, client *http.Client, header ...map[string]string) (err error) {
	var body []byte
	if len(header) > 0 {
		body, err = post(url, payload, client, header[0])
	} else {
		body, err = post(url, payload, client)
	}
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		_ = errors.New("decode resp to struct err")
	}
	return
}

func get(url string, param ...Param) (res []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// deal with headers and query param if passed
	if len(param) > 0 {
		for k, v := range param[0].Query {
			q := req.URL.Query()
			q.Add(k, v)
			req.URL.RawQuery = q.Encode()
		}
		for k, v := range param[0].Header {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status code %d", resp.StatusCode)
		return
	}
	return ioutil.ReadAll(resp.Body)
}

func post(url string, payload interface{}, client *http.Client, header ...map[string]string) (res []byte, err error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	if client == nil {
		client = &http.Client{}
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	if len(header) > 0 {
		for k, v := range header[0] {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(jsonBytes)))
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status code %d", resp.StatusCode)
		return
	}
	return ioutil.ReadAll(resp.Body)
}
