package logtail

import (
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func Serve(ctx context.Context, conf *Config) {
	logsSvc, err := NewLogsService(conf.FileToTail)
	if err != nil {
		logger.Fatal(ctx, "can't set up new logs service, err: %+v", err)
	}

	defer func(logsSvc LogsService, ctx context.Context) {
		if err = logsSvc.Shutdown(ctx); err != nil {
			logger.Fatal(ctx, "failed to shutdown logs service, err: %+v", err)
		}
	}(logsSvc, ctx)

	logger.Print(ctx, "new logs service set up successfully")

	logsCont := NewLogsController(logsSvc)
	router := mux.NewRouter()

	router.Use(
		handlers.RecoveryHandler(),
		mux.CORSMethodMiddleware(router),
		RequestIDMiddleware,
		ResponseTimeMiddleware,
	)

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
