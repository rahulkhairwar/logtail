package logtail

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

func Serve(ctx context.Context, conf *Config, wg *sync.WaitGroup) {
	defer wg.Done()

	r := mux.NewRouter()

	r.Use(mux.CORSMethodMiddleware(r), RequestIDMiddleware, ResponseTimeMiddleware)
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)

	addr := fmt.Sprintf("127.0.0.1:%v", conf.Port)

	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	errs := make(chan error)

	logger.Print(ctx, "starting LogTail server on port {%v}", conf.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			errs <- err
		}
	}()

	for {
		select {
		case <-ctx.Done():
			if err := srv.Shutdown(ctx); err != nil {
				logger.Fatal(ctx, "failed to shutdown server, err: %+v", err)
			}

			logger.Fatal(ctx, "Application context completed, stopping server.")

		case err := <-errs:
			logger.Fatal(ctx, "Error while serving: %+v", err)
		}
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home Handler!!"))
	if err != nil {
		logger.Print(r.Context(), "error writing response: %v", err)
	}
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		startTime := time.Now()

		logger.Print(ctx, "incoming request [%v]{%v}", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		logger.Print(ctx, "request completed in {%v}", startTime.Sub(time.Now()))
	})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(requestIDKey)
		ctx := r.Context()

		if reqID == "" {
			reqID, err := getRandomRequestID()
			if err != nil {
				logger.Print(ctx, "failed to generate random UUID, err: %v", err)

				return
			}

			r = r.WithContext(context.WithValue(ctx, requestIDKey, reqID))
		}

		next.ServeHTTP(w, r)
	})
}
