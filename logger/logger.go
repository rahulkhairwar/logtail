package logger

import (
	"context"
	"fmt"
	"github.com/rahulkhairwar/logtail/constants"
	"log"
	"os"
)

var (
	logger = _logger{
		l: log.Default(),
	}
)

const DefaultFlags = log.Ltime | log.Lshortfile

type _logger struct {
	l *log.Logger
}

func Initialize() {
	logger.l.SetFlags(DefaultFlags)
}

func Print(ctx context.Context, msg string, args ...interface{}) {
	_ = logger.l.Output(2, prepareMsg(ctx, msg, args...))
}

func Fatal(ctx context.Context, msg string, args ...interface{}) {
	_ = logger.l.Output(2, prepareMsg(ctx, msg, args...))
	os.Exit(1)
}

func prepareMsg(ctx context.Context, msg string, args ...interface{}) string {
	msg = fmt.Sprintf(msg, args...)

	reqID := ctx.Value(constants.RequestIDKey)
	if reqID != nil {
		msg = fmt.Sprintf("[request_id = %v] %v", reqID, msg)
	}

	return msg
}
