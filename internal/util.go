package internal

import (
	"context"

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

// todo-rahul: tests.
func RunInBatches(ctx context.Context, batchSize, min, max int, fn func(ctx context.Context, start, end int) error) error {
	start := min
	end := start + batchSize

	for {
		if err := fn(ctx, start, end); err != nil {
			return eris.Wrap(err, "...")
		}

		start += batchSize
		end += batchSize

		if start > max {
			break
		}
	}

	return nil
}
