package logtail

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func getRandomRequestID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "new random uuid")
	}

	return id.String(), nil
}
