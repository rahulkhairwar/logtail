package logtail

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var (
	_emptyCtx = context.Background()
	_reqID    = "asdf-ghjk"
	_reqIDCtx = context.WithValue(_emptyCtx, requestIDKey, _reqID)
	_msg      = "This is a sample message."
	_argsMsg  = _msg + " {%v}."
)

type args struct {
	ctx  context.Context
	msg  string
	args []interface{}
}

type test struct {
	name string
	args args
	want string
}

func getTests() []test {
	return []test{
		{
			name: "empty ctx, simple msg",
			args: args{
				ctx: _emptyCtx,
				msg: _msg,
			},
			want: _msg + "\n",
		},
		{
			name: "request ID ctx, simple msg",
			args: args{
				ctx: _reqIDCtx,
				msg: _msg,
			},
			want: fmt.Sprintf("[request_id = %v] %v\n", _reqID, _msg),
		},
		{
			name: "request ID ctx, args msg",
			args: args{
				ctx:  _reqIDCtx,
				msg:  _argsMsg,
				args: []interface{}{123},
			},
			want: fmt.Sprintf("[request_id = %v] %v {%v}.\n", _reqID, _msg, 123),
		},
	}
}

func TestLogger_Print(t *testing.T) {
	log.SetFlags(0)

	for _, tt := range getTests() {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.Buffer{}

			log.SetOutput(&b)
			logger.Print(tt.args.ctx, tt.args.msg, tt.args.args...)
			assert.Equal(t, tt.want, b.String())
		})
	}

	log.SetOutput(os.Stderr)
}

func TestLogger_Fatal(t *testing.T) {
	log.SetFlags(0)

	tests := getTests()
	for i, tt := range tests {
		testName := fmt.Sprintf("FATAL_TEST_%d", i)
		logFile := fmt.Sprintf("FATAL_%d.log", i)

		f := openOrCreateFile(logFile)

		logger._log.SetOutput(f)

		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv(testName) == "1" {
				logger.Fatal(tt.args.ctx, tt.args.msg, tt.args.args...)
				return
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestLogger_Fatal")
			cmd.Env = append(os.Environ(), testName+"=1")

			err := cmd.Run()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}

			t.Errorf("process ran with err %v, want exit status 1", err)
		})

		assert.Equal(t, tt.want, readLogs(logFile))
		deleteFile(f, logFile)
	}

	log.SetOutput(os.Stderr)
}

func openOrCreateFile(file string) *os.File {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logger._log.Fatalf("error opening file: %v", err)
	}

	return f
}

func deleteFile(f *os.File, file string) {
	_ = f.Close()

	if err := os.Remove(file); err != nil {
		panic("failed to delete file : " + err.Error())
	}
}

func readLogs(logFile string) string {
	bts, err := ioutil.ReadFile(logFile)
	if err != nil {
		panic("failed to read file : " + err.Error())
	}

	return string(bts)
}

func Test_prepareMsg(t *testing.T) {
	remNewLine := func(s string) string {
		return strings.Trim(s, "\n")
	}

	tests := getTests()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, remNewLine(tt.want), prepareMsg(tt.args.ctx, remNewLine(tt.args.msg), tt.args.args...))
		})
	}
}
