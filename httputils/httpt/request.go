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

func NewRequestPostJson(target string, v any) (*http.Request, error) {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(v); err != nil {
		return nil, err
	}
	r := NewRequestPost(target, body)
	r.Header.Set("Content-Type", "application/json")
	return r, nil
}

func NewRequestGet(target string) *http.Request {
	return httptest.NewRequest(http.MethodGet, target, nil)
}
