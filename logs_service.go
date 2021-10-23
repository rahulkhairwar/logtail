package logtail

import "context"

type LogsService interface {
	GetLogs(context.Context, int) ([]string, error)
	Shutdown(context.Context) error
}

type logsService struct {
	records records
}

// GetLogs returns the new logs available. If pageSize is provided, pageSize logs are returned, else defaultPageSize
// logs are returned. If the total available logs are less than the determined pageSize, all those logs are returned.
func (l logsService) GetLogs(ctx context.Context, pageSize int) ([]string, error) {
	if pageSize == 0 {
		pageSize = defaultPageSize
	}

	logger.Print(ctx, "get logs, pageSize {%v}", pageSize)

	var logs []string

	for i := 0; i < pageSize; i++ {
		lg, err := l.records.Next()
		if err != nil {
			if err == ErrNoRecords {
				return logs, nil
			}

			logger.Print(ctx, "error while fetching next log record: %v", err)

			return nil, err
		}

		logs = append(logs, lg)
	}

	return logs, nil
}

// Shutdown closes all resources being used by the service.
// Returns any error occurred during shutdown.
func (l logsService) Shutdown(context.Context) error {
	return l.records.Close()
}

// NewLogsService returns an instance of LogsService. Should ideally call this with deferred LogsService.Shutdown to
// properly close resources.
func NewLogsService() LogsService {
	return &logsService{
		records: records{
			ch: make(chan string),
		},
	}
}
