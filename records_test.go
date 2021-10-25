package logtail

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

const (
	_recordsFile = "records.log"
	_logText     = `1) Line 1
2) Line 2
3) Line 3
4) Line 4
5) Line 5
6) Line 6
7) Line 7
8) Line 8
9) Line 9
10) Line 10`
)

func Test_newRecords(t *testing.T) {
	t.Run("tail a non-existing file", func(t *testing.T) {
		got, err := newRecords("some_random_file")

		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("tail systemLogFile", func(t *testing.T) {
		got, err := newRecords("/var/log/system.log")

		assert.NoError(t, err)
		assert.NotNil(t, got)
	})
}

func Test_records_Next(t *testing.T) {
	f := openOrCreateFile(_recordsFile)
	defer deleteFile(f, _recordsFile)

	assert.NoError(t, writeTempDataToFile(f))

	expectedLogs := strings.Split(_logText, "\n")

	t.Run("file with 10 records", func(t *testing.T) {
		var logs []string

		r, err := newRecords(_recordsFile)
		assert.NoError(t, err)

		for i := 0; i < 10; i++ {
			got, err := r.Next()

			logs = append(logs, got)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		}

		assert.Equal(t, expectedLogs, logs)
	})

	t.Run("read timeout", func(t *testing.T) {
		var logs []string

		r, err := newRecords(_recordsFile)
		assert.NoError(t, err)

		for i := 0; i < 10; i++ {
			got, err := r.Next()

			logs = append(logs, got)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		}

		// trying to tail an 11th log line, which doesn't exist in the file.
		got, err := r.Next()

		assert.EqualError(t, err, ErrNoRecords.Error())
		assert.Empty(t, got)
		assert.Equal(t, expectedLogs, logs)
	})

	t.Run("dead tail", func(t *testing.T) {
		var logs []string

		r, err := newRecords(_recordsFile)
		assert.NoError(t, err)

		for i := 0; i < 5; i++ {
			got, err := r.Next()

			logs = append(logs, got)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
		}

		r.t.Kill(nil)

		// trying to tail a 6th log line, but the tail has been killed off.
		got, err := r.Next()

		assert.EqualError(t, err, ErrClosed.Error())
		assert.Empty(t, got)
		assert.Equal(t, expectedLogs[:5], logs)
	})
}

func Test_records_Close(t *testing.T) {
	f := openOrCreateFile(_recordsFile)
	defer deleteFile(f, _recordsFile)

	assert.NoError(t, writeTempDataToFile(f))

	t.Run("check", func(t *testing.T) {
		r, err := newRecords(_recordsFile)
		assert.NoError(t, err)

		assert.NoError(t, r.Close())
	})
}

func writeTempDataToFile(f *os.File) error {
	_, err := f.Write([]byte(_logText))
	return err
}
