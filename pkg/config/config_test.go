package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	// run tests
	code := m.Run()

	// end
	os.Exit(code)
}

func TestLoad(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// prepare
		os.Setenv("APP_HOST_NAME", "http://gateway.marvel.com")
		os.Setenv("APP_PUBLIC_KEY", "public_key")
		os.Setenv("APP_PRIVATE_KEY", "private_key")
		os.Setenv("APP_TS", "1")

		cfg, err := Load()

		assert.NoError(t, err)
		assert.Equal(t, cfg, &Config{
			"info",
			"severity",
			":8080",
			60 * time.Second,
			"http://gateway.marvel.com",
			"private_key",
			"public_key",
			"1",
		})
	})
}
