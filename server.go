package logtail

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"sync"
	"time"
)

func Serve(ctx context.Context, conf *Config, wg *sync.WaitGroup) {
	defer wg.Done()

	logsSvc := NewLogsService()
	// homeCont := newHomeController()
	logsCont := NewLogsController(logsSvc)
	router := mux.NewRouter()

	router.Use(
		handlers.RecoveryHandler(),
		mux.CORSMethodMiddleware(router),
		RequestIDMiddleware,
		ResponseTimeMiddleware,
	)

	// homeCont.SetupRoutes(router)
	logsCont.SetupRoutes(router)

	http.Handle("/", router)

	addr := fmt.Sprintf("127.0.0.1:%v", conf.Port)
	srv := &http.Server{
		Handler:      router,
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Home Handler!!"))
	if err != nil {
		logger.Print(r.Context(), "error writing response: %v", err)
	}
}

// todo: tests for file.
