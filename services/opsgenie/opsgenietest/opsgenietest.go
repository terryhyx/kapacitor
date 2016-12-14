package opsgenietest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type Server struct {
	ts       *httptest.Server
	URL      string
	requests chan Request
	Requests <-chan Request
	closed   bool
}

func NewServer(count int) *Server {
	requests := make(chan Request, count)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		or := Request{
			URL: r.URL.String(),
		}
		dec := json.NewDecoder(r.Body)
		dec.Decode(&or.PostData)
		requests <- or
	}))
	return &Server{
		ts:       ts,
		URL:      ts.URL,
		requests: requests,
		Requests: requests,
	}
}

func (s *Server) Close() {
	if s.closed {
		return
	}
	s.closed = true
	s.ts.Close()
	close(s.requests)
}

type Request struct {
	URL      string
	PostData PostData
}

type PostData struct {
	ApiKey      string                 `json:"apiKey"`
	Message     string                 `json:"message"`
	Entity      string                 `json:"entity"`
	Alias       string                 `json:"alias"`
	Note        string                 `json:"note"`
	Details     map[string]interface{} `json:"details"`
	Description string                 `json:"description"`
	Teams       []string               `json:"teams"`
	Recipients  []string               `json:"recipients"`
}