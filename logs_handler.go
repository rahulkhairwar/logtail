package logtail

import "net/http"

func logsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pageSize := r.FormValue("pageSize")

	if pageSize != "" {
		logger.Print(ctx, "pageSize: %v", pageSize)
	}

	_, err := w.Write([]byte("Logs Handler!!"))
	if err != nil {
		logger.Print(ctx, "error writing response: %v", err)
	}
}
