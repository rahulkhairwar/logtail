package logtail

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

var _testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	logger.Print(r.Context(), "[_testHandler]")
})

func TestRequestIDMiddleware(t *testing.T) {
	log.SetFlags(0)
	want := "[request_id = asdf-ghjk] [_testHandler]\n"

	t.Run("existing requestID", func(t *testing.T) {
		b := bytes.Buffer{}
		log.SetOutput(&b)

		handler := RequestIDMiddleware(_testHandler)
		req := newRequest(http.MethodGet, "/temp")

		req.Header.Set(requestIDKey, _reqID)
		handler.ServeHTTP(httptest.NewRecorder(), req)
		assert.Equal(t, want, b.String())
	})
	t.Run("non-existing requestID", func(t *testing.T) {
		b := bytes.Buffer{}
		log.SetOutput(&b)

		handler := RequestIDMiddleware(_testHandler)
		handler.ServeHTTP(httptest.NewRecorder(), newRequest(http.MethodGet, "/temp"))

		s := regexp.MustCompile(`\[request_id = [\w\-]+]`).ReplaceAllString(b.String(), "[request_id = asdf-ghjk]")

		assert.Equal(t, want, s)
	})

	log.SetOutput(os.Stderr)
	log.SetFlags(defaultFlags)
}

func TestResponseTimeMiddleware(t *testing.T) {
	t.Run("check positive response time", func(t *testing.T) {
		b := bytes.Buffer{}
		log.SetOutput(&b)

		handler := ResponseTimeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Print(r.Context(), "hello")
		}))

		req := newRequest(http.MethodGet, "/temp")
		handler.ServeHTTP(httptest.NewRecorder(), req)

		resp := b.String()
		firstLog := strings.Index(resp, "incoming request [GET]{/temp}")
		secondLog := strings.Index(resp, "request completed in")

		assert.Positive(t, firstLog)
		assert.Positive(t, secondLog)
		assert.True(t, firstLog < secondLog)
	})
}
