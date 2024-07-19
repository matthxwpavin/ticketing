package httpt

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

func NewRequestPost(target string, body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodPost, target, body)
}

func NewRequestPut(target string, body io.Reader) *http.Request {
	return httptest.NewRequest(http.MethodPut, target, body)
}

func NewRequestDelete(target string) *http.Request {
	return httptest.NewRequest(http.MethodDelete, target, nil)
}

func NewRequestGet(target string) *http.Request {
	return httptest.NewRequest(http.MethodGet, target, nil)
}

func NewRequestPostJson(target string, v any) (*http.Request, error) {
	return jsonRequest(target, v, NewRequestPost)
}

func NewRequestPutJson(target string, v any) (*http.Request, error) {
	return jsonRequest(target, v, NewRequestPut)
}

func jsonRequest(target string, v any, createRequest func(string, io.Reader) *http.Request) (*http.Request, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(v); err != nil {
		return nil, err
	}
	r := createRequest(target, body)
	r.Header.Set("Content-Type", "application/json")
	return r, nil
}
