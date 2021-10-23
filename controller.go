package logtail

import (
	"github.com/gorilla/mux"
	"net/http"
)

type logsController struct {
	svc LogsService
}

// set up route -> handler mappings here.

func (l *logsController) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/logs", logsHandler).Name("logsHandler")
	r.HandleFunc("/logs", logsHandler).Name("logsHandler").Queries("pageSize", "{pageSize:[0-9]+}").Methods(http.MethodGet)
}

func NewLogsController(svc LogsService) *logsController {
	return &logsController{
		svc: svc,
	}
}

type homeController struct{}

func (l *homeController) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/", homeHandler).Methods(http.MethodGet)
}

func newHomeController() *homeController {
	return &homeController{}
}
