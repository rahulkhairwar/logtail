package logtail

import (
	"context"
	"github.com/rahulkhairwar/logtail/internal"
	"sync"
	"testing"
)

func TestServe(t *testing.T) {
	type args struct {
		ctx  context.Context
		conf *internal.Config
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
