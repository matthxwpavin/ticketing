package httpt

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type HandlerTesting struct {
	h        http.Handler
	attrs    *TestingAttributes
	afterRun AfterFunc
}

func NewHandlerTesting(h http.Handler, attrs *TestingAttributes) *HandlerTesting {
	return &HandlerTesting{
		h:     h,
		attrs: attrs,
	}
}

type AfterFunc func(t *testing.T, r *http.Response)

type TestingAttributes struct {
	Name           string
	TestingRequest func(*testing.T) *http.Request
	StatusCode     int
	StatusCodeFunc func(statusCode int) bool
}

type HandlerPrepared struct {
	h http.Handler
}

func Prepare(h http.Handler) *HandlerPrepared {
	return &HandlerPrepared{h: h}
}

func (s *HandlerPrepared) Testing(attrs *TestingAttributes) *HandlerTesting {
	return &HandlerTesting{
		h:     s.h,
		attrs: attrs,
	}
}

func (s *HandlerTesting) After(afterRun AfterFunc) *HandlerTesting {
	s.afterRun = afterRun
	return s
}

func (s *HandlerTesting) Run(t *testing.T) {
	a := s.attrs
	t.Run(a.Name, func(t *testing.T) {
		r := a.TestingRequest(t)
		if r == nil {
			t.Fatal("the testing request is nil")
		}

		w := httptest.NewRecorder()
		s.h.ServeHTTP(w, r)

		rs := w.Result()
		if a.StatusCodeFunc != nil {
			if !a.StatusCodeFunc(rs.StatusCode) {
				t.Fatalf("status code func returns false, expected: %v", rs.StatusCode)
			}
		} else if rs.StatusCode != 0 && rs.StatusCode != a.StatusCode {
			t.Fatalf(
				"status code is unexpected, expected: %v, received: %v",
				a.StatusCode,
				rs.StatusCode,
			)
		}
		if s.afterRun != nil {
			s.afterRun(t, rs)
		}
	})
}

type HandlerTestingList []*HandlerTesting

func (l HandlerTestingList) Run(t *testing.T) {
	for _, c := range l {
		c.Run(t)
	}
}
