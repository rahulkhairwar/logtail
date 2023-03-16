package internal

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewLogsController(t *testing.T) {
	type args struct {
		svc LogsService
	}
	tests := []struct {
		name string
		args args
		want *logsController
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewLogsController(tt.args.svc), "NewLogsController(%v)", tt.args.svc)
		})
	}
}

func Test_genericHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		g    genericHandler
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.g.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}

func Test_logsController_Get(t *testing.T) {
	type fields struct {
		svc LogsService
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logsController{
				svc: tt.fields.svc,
			}
			got, err := l.Get(tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.r)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.r)
		})
	}
}

func Test_logsController_SetupRoutes(t *testing.T) {
	type fields struct {
		svc LogsService
	}
	type args struct {
		r *mux.Router
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &logsController{
				svc: tt.fields.svc,
			}
			l.SetupRoutes(tt.args.r)
		})
	}
}
