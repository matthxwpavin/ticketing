package httptesting

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Run(t *testing.T, testing Testing) {
	testing.Run(t)
}

type Testing struct {
	Handler  http.Handler
	Specs    TestingSpecifications
	AfterRun AfterRunFunc
}

func (s Testing) After(afterRun AfterRunFunc) Testing {
	s.AfterRun = afterRun
	return s
}

func (s Testing) Run(t *testing.T) {
	a := s.Specs
	t.Run(a.Name, func(t *testing.T) {
		r := a.TestingRequest(t)
		if r == nil {
			t.Fatal("the testing request is nil")
		}

		w := httptest.NewRecorder()
		s.Handler.ServeHTTP(w, r)

		rs := w.Result()
		if a.StatusCodeFunc != nil {
			if !a.StatusCodeFunc(rs.StatusCode) {
				t.Fatalf("status code func returns false, got: %v", rs.StatusCode)
			}
		} else if a.StatusCode != 0 && rs.StatusCode != a.StatusCode {
			t.Fatalf(
				"status code is unexpected, expected: %v, received: %v",
				a.StatusCode,
				rs.StatusCode,
			)
		}
		if s.AfterRun != nil {
			s.AfterRun(t, rs)
		}
	})
}

type AfterRunFunc func(t *testing.T, r *http.Response)

type TestingSpecifications struct {
	Name           string
	TestingRequest func(*testing.T) *http.Request
	StatusCode     int
	StatusCodeFunc func(statusCode int) bool
}

type prepared struct {
	handler http.Handler
}

func Prepare(handler http.Handler) *prepared {
	return &prepared{handler: handler}
}

func (s *prepared) Testing(specs TestingSpecifications) Testing {
	return Testing{
		Handler: s.handler,
		Specs:   specs,
	}
}

type TestingList []Testing

func (l TestingList) Run(t *testing.T) {
	for _, c := range l {
		c.Run(t)
	}
}
