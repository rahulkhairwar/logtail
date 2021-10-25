package logtail

import (
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
)

func getRandomRequestID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", eris.Wrap(err, "new random uuid")
	}

	return id.String(), nil
}
