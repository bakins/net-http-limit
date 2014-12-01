package limit

import (
	"net/http"
	"time"
)

// Limiter is used to create limiting handlers. All handlers share the created limit
type Limiter struct {
	timeout time.Duration
	max     uint
	channel chan struct{}
}

// LimiterHandler is used to wrap an http.Handler with a request limit
type LimiterHandler struct {
	handler http.Handler
	limiter *Limiter
}

// New creates a new Limit Handler creator.  When max is reached, new requests may
// wait up until timeout. a 503 is returned if the timeout is reachd
func New(timeout time.Duration, max uint) *Limiter {
	l := &Limiter{
		timeout: timeout,
		max:     max,
		channel: make(chan struct{}, max),
	}

	for i := 0; i < max; i++ {
		l.channel <- struct{}{}
	}

	return l
}

// Handler wraps an http.Handler with a limiting handler
func (l *Limiter) Handler(h http.Handler) http.Handler {
	return &LimiterHandler{
		handler: h,
		limiter: l,
	}
}

// ServeHTTP will attempt to "reserve" a slot and wait up to timeout to do so.
// Will return a 503 error.
func (h *LimiterHandler) ServeHTTP(w ResponseWriter, r *Request) {
	select {
	case <-h.limiter.channel:
		defer func(l *Limiter) {
			l.channel <- struct{}{}
		}(h.limiter)
		h.handler(w, r)
	case <-time.After(h.limiter.timeout):
		http.Error(w, "max concurreny", 503)
	}
}
