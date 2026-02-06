package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rahulkhairwar/logtail/logger"
	"github.com/rotisserie/eris"
)

const (
	pageSizeKey = "pageSize"
)

type logsController struct {
	svc LogsService
}

func (l *logsController) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", l.ServeUI).Methods(http.MethodGet)
	r.Handle("/logs", genericHandler(l.Get)).Name("logsHandler")
	r.Handle("/logs", genericHandler(l.Get)).Name("logsHandler").Queries(pageSizeKey, "{pageSize:[0-9]+}").Methods(http.MethodGet)
	r.HandleFunc("/logs/stream", l.Stream).Methods(http.MethodGet)
}

func (g genericHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	status := http.StatusOK

	data, err := g(r)
	if err != nil {
		status = http.StatusInternalServerError
	}

	res, err := json.Marshal(data)
	if err != nil {
		errMsg := fmt.Sprintf("failed to parse response body: %+v", err)
		http.Error(w, errMsg, status)

		return
	}

	ctx := r.Context()

	if _, err = w.Write(res); err != nil {
		logger.Print(ctx, "failed to write response to writer")
		return
	}
}

func (l *logsController) Get(r *http.Request) (interface{}, error) {
	ctx := r.Context()
	pageSizeStr := r.FormValue(pageSizeKey)

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return nil, eris.Wrap(err, "invalid page size")
	}

	res, err := l.svc.GetLogs(ctx, pageSize)
	if err != nil {
		return nil, eris.Wrapf(err, "call service")
	}

	return res, nil
}

func (l *logsController) Stream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()

	err := l.svc.StreamLogs(ctx, func(line string) error {
		_, writeErr := fmt.Fprintf(w, "data: %s\n\n", line)
		if writeErr != nil {
			return writeErr
		}

		flusher.Flush()
		return nil
	})

	if err != nil {
		logger.Print(ctx, "stream error: %v", err)
	}
}

func (l *logsController) ServeUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

func NewLogsController(svc LogsService) *logsController {
	return &logsController{
		svc: svc,
	}
}

type genericHandler func(*http.Request) (interface{}, error)
