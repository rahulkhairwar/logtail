package logtail

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rotisserie/eris"
	"net/http"
	"strconv"
)

const (
	pageSizeKey = "pageSize"
)

type logsController struct {
	svc LogsService
}

func (l *logsController) SetupRoutes(r *mux.Router) {
	r.Handle("/logs", genericHandler(l.Get)).Name("logsHandler")
	r.Handle("/logs", genericHandler(l.Get)).Name("logsHandler").Queries(pageSizeKey, "{pageSize:[0-9]+}").Methods(http.MethodGet)
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
		logger.Print(ctx, "failed to write ")
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

func NewLogsController(svc LogsService) *logsController {
	return &logsController{
		svc: svc,
	}
}

type genericHandler func(*http.Request) (interface{}, error)
