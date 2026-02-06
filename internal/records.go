package internal

import (
	"context"
	"time"

	"github.com/nxadm/tail"
	"github.com/rotisserie/eris"
)

type records struct {
	t *tail.Tail
}

func newRecords(file string) (*records, error) {
	conf := tail.Config{
		MustExist: true,
		Poll:      true,
		Follow:    true,
	}

	t, err := tail.TailFile(file, conf)
	if err != nil {
		return nil, eris.Wrapf(err, "tail file {%v}", file)
	}

	return &records{
		t: t,
	}, nil
}

func (r *records) Next() (string, error) {
	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()

	select {
	case s := <-r.t.Lines:
		if s == nil {
			return "", ErrClosed
		}

		return s.Text, nil
	case <-r.t.Dying():
		return "", ErrClosed
	case <-timer.C:
		return "", ErrNoRecords
	}
}

func (r *records) NextBlocking(ctx context.Context) (string, error) {
	select {
	case s := <-r.t.Lines:
		if s == nil {
			return "", ErrClosed
		}

		return s.Text, nil
	case <-r.t.Dying():
		return "", ErrClosed
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (r *records) Close() error {
	if err := r.t.Stop(); err != nil {
		return eris.Wrapf(err, "stop tailing")
	}

	r.t.Cleanup()

	return nil
}
