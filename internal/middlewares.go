package internal

import (
	"context"
	"github.com/rahulkhairwar/logtail/constants"
	"github.com/rahulkhairwar/logtail/logger"
	"net/http"
	"time"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(constants.RequestIDKey)
		ctx := r.Context()

		if reqID == "" {
			var err error
			reqID, err = getRandomRequestID()
			if err != nil {
				logger.Print(ctx, "failed to generate random UUID, err: %v", err)

				return
			}
		}

		r = r.WithContext(context.WithValue(ctx, constants.RequestIDKey, reqID))
		next.ServeHTTP(w, r)
	})
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		startTime := time.Now()

		logger.Print(ctx, "incoming request [%v]{%v}", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		logger.Print(ctx, "request completed in {%+v}", time.Now().Sub(startTime))
	})
}
