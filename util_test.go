package logtail

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getRandomRequestID(t *testing.T) {
	set := make(map[string]bool)
	tests := 100

	t.Run(fmt.Sprintf("check unique UUIDs for %v runs", tests), func(t *testing.T) {
		for i := 0; i < tests; i++ {
			uuid, err := getRandomRequestID()

			assert.NoError(t, err)
			assert.NotContains(t, set, uuid)
			set[uuid] = true
		}
	})
}
