package logtail

import (
	"context"
	"net/http"
	"sync"
	"testing"
)

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
