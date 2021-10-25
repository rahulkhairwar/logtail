package logtail

import (
	"context"
	"fmt"
	"log"
	"os"
)

var (
	logger = Logger{
		_log: log.Default(),
	}
)

const defaultFlags = log.Ltime | log.Lshortfile

type Logger struct {
	_log *log.Logger
}

func init() {
	logger._log.SetFlags(defaultFlags)
}

func (l *Logger) Print(ctx context.Context, msg string, args ...interface{}) {
	_ = logger._log.Output(2, prepareMsg(ctx, msg, args...))
}

func (l *Logger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	_ = logger._log.Output(2, prepareMsg(ctx, msg, args...))
	os.Exit(1)
}

func prepareMsg(ctx context.Context, msg string, args ...interface{}) string {
	msg = fmt.Sprintf(msg, args...)

	reqID := ctx.Value(requestIDKey)
	if reqID != nil {
		msg = fmt.Sprintf("[request_id = %v] %v", reqID, msg)
	}

	return msg
}
