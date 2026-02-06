package logtail

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rahulkhairwar/logtail/internal"
	"github.com/rahulkhairwar/logtail/logger"
)

func Serve(ctx context.Context, conf *internal.Config) {
	logsSvc, err := internal.NewLogsService(conf.FileToTail)
	if err != nil {
		logger.Fatal(ctx, "can't set up new logs service, err: %+v", err)
	}

	defer func(logsSvc internal.LogsService, ctx context.Context) {
		if err = logsSvc.Shutdown(ctx); err != nil {
			logger.Fatal(ctx, "failed to shutdown logs service, err: %+v", err)
		}
	}(logsSvc, ctx)

	logger.Print(ctx, "new logs service set up successfully")

	logsCont := internal.NewLogsController(logsSvc)
	router := mux.NewRouter()

	router.Use(
		handlers.RecoveryHandler(),
		mux.CORSMethodMiddleware(router),
		internal.RequestIDMiddleware,
		internal.ResponseTimeMiddleware,
	)

	logsCont.SetupRoutes(router)

	http.Handle("/", router)

	addr := fmt.Sprintf("127.0.0.1:%v", conf.Port)
	srv := &http.Server{
		Handler:     router,
		Addr:        addr,
		ReadTimeout: 15 * time.Second,
		// No WriteTimeout — SSE streams need long-lived connections.
	}

	// Buffered so the goroutine can always send without blocking,
	// even if shutdown has already started via ctx.Done().
	errs := make(chan error, 1)

	logger.Print(ctx, "starting LogTail server on port {%v}", conf.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errs <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Print(ctx, "failed to shutdown server, err: %+v", err)
		}

		logger.Print(ctx, "server stopped")
	case err := <-errs:
		logger.Fatal(ctx, "Error while serving: %+v", err)
	}
}
