package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentials_ToMd5Hash(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		hashed, err := toMd5Hash(Credentials{
			PublicKey:  "public_key",
			PrivateKey: "private_key",
			Ts:         "1",
		})
		assert.NoError(t, err)
		assert.Equal(t, "f21c2ff8dbfa6ad2170916f0f32d375d", hashed)
	})
}

func TestNewCredentials(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cred, err := NewCredentials("public_key", "private_key", "1")
		assert.NoError(t, err)

		hashed, err := toMd5Hash(cred)
		require.NoError(t, err)
		assert.Equal(t, "public_key", cred.PublicKey)
		assert.Equal(t, "private_key", cred.PrivateKey)
		assert.Equal(t, "1", cred.Ts)
		assert.Equal(t, hashed, cred.Hashed)
	})
}

func TestCredentials_ToBaseMarvelRequestParams(t *testing.T) {
	cred, err := NewCredentials("public_key", "private_key", "1")
	require.NoError(t, err)

	hashed, err := toMd5Hash(cred)
	require.NoError(t, err)

	params, err := cred.ToBaseMarvelRequestParams()

	assert.NoError(t, err)
	assert.Equal(t, "public_key", params.Get(marvelParamApiKey))
	assert.Equal(t, "1", params.Get(marvelParamTs))
	assert.Equal(t, hashed, params.Get(marvelParamHash))
}

func TestCredentials_ToCharacterIDsReqParams(t *testing.T) {
	cred, err := NewCredentials("public_key", "private_key", "1")
	require.NoError(t, err)

	hashed, err := toMd5Hash(cred)
	require.NoError(t, err)

	since := time.Now()
	params, err := cred.ToCharacterIDsReqParams(100, 0, since)

	assert.NoError(t, err)
	assert.Equal(t, "public_key", params.Get(marvelParamApiKey))
	assert.Equal(t, "1", params.Get(marvelParamTs))
	assert.Equal(t, hashed, params.Get(marvelParamHash))
	assert.Equal(t, "100", params.Get(marvelParamLimit))
	assert.Equal(t, "0", params.Get(marvelParamOffset))
	assert.Equal(t, since.Format(marvelParamModifiedSinceFormat), params.Get(marvelParamModifiedSince))
}
