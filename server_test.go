package logtail

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"sync"
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
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := ResponseTimeMiddleware(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResponseTimeMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServe(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf *Config
		wg   *sync.WaitGroup
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}
