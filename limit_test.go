package limit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	h "github.com/bakins/test-helpers"
)

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

// This test just tests than limit actually is a valif http.handler.
// need a test to test limiting
func TestHandler(t *testing.T) {
	l := New(time.Second, 2)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, "Hello World\n")
	})

	w := httptest.NewRecorder()

	l.Handler(handler).ServeHTTP(w, newRequest("GET", "/foo"))

	h.Assert(t, w.Body.String() == "Hello World\n", "body does not match")

}
