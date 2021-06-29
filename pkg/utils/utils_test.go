package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMd5(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		s := "1481732840d168b3244a9f2a4efb3c21d2b1ce2e05d04d11a2b5d49b6aa05a49998b5f083"

		// when
		result, err := HashMd5(s)

		// then
		assert.NoError(t, err)
		assert.Equal(t, "21d7e05400dfc7ae57cfe3a3ed19caf4", result)
	})
}
