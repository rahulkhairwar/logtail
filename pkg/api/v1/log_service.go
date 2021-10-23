package v1

import "context"

type LogService interface {
	Get(context.Context, *QueryParams) ([]*string, error)
}

type QueryParams struct {
	PageSize int
}
