package logtail

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrNoRecords = Error("no_records")
)
