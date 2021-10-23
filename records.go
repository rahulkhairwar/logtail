package logtail

type records struct {
	ch chan string
}

func (r *records) Next() (string, error) {
	select {
	case s := <-r.ch:
		return s, nil
	default:
		return "", ErrNoRecords
	}
}

func (r *records) Close() error {
	close(r.ch)

	// currently, no possibility of error.
	return nil
}
