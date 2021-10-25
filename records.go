package logtail

import (
	"github.com/nxadm/tail"
	"github.com/rotisserie/eris"
	"time"
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
	select {
	case s := <-r.t.Lines:
		return s.Text, nil
	case <-r.t.Dying():
		return "", ErrClosed
	case <-time.After(10 * time.Millisecond):
		return "", ErrNoRecords
	}
}

func (r *records) Close() error {
	if err := r.t.Stop(); err != nil {
		return eris.Wrapf(err, "stop tailing")
	}

	r.t.Cleanup()

	return nil
}
