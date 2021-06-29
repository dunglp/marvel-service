package marvel

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
	"xendit-technical-assessment/pkg/domain"

	"github.com/patrickmn/go-cache"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

var svc Service

const (
	id  = 1017100
	uri = "http://localhost:8080/characters/1017100"
)

func TestMain(m *testing.M) {

	// prepare
	svc = NewService("", domain.Credentials{}, cache.New(24*time.Hour, 24*time.Hour))

	// run tests
	code := m.Run()

	// end
	os.Exit(code)
}

func TestService_GetCharacterDetails(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/character_details.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()

		character, err := svc.GetCharacterDetails(id)

		assert.NoError(t, err)
		assert.Equal(t, id, character.ID)
	})
	t.Run("missing apikey", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/errors/missing_apikey.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()

		_, err := svc.GetCharacterDetails(id)

		assert.EqualError(t, err, "internal server error")
	})
	t.Run("missing hash", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/errors/missing_hash.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()

		_, err := svc.GetCharacterDetails(id)

		assert.EqualError(t, err, "internal server error")
	})
	t.Run("missing timestamp", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/errors/missing_timestamp.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()

		_, err := svc.GetCharacterDetails(id)

		assert.EqualError(t, err, "internal server error")
	})
	t.Run("invalid credentials", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/errors/invalid_credentials.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()

		_, err := svc.GetCharacterDetails(id)

		assert.EqualError(t, err, "internal server error")
	})
}

func TestService_GetCharacterIds(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/character_ids.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()
		ids, err := svc.GetCharacterIds()

		assert.NoError(t, err)
		assert.Len(t, ids, 5)
	})
}

func TestService_GetMarvelCharacters(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rawData, err := ioutil.ReadFile("../../fixtures/character_ids.json")
			require.NoError(t, err)
			rw.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(rw, string(rawData))
		}))
		// Close the server when test finishes
		defer server.Close()

		// prepare
		svc.Hostname = server.URL
		svc.HTTPClient = server.Client()
		limit := 5
		offset := 0
		cdc, err := svc.GetMarvelCharacters(limit, offset, time.Now())

		assert.NoError(t, err)
		assert.Len(t, cdc.Results, 5)
		assert.Equal(t, cdc.Total, 5)
	})
}

func TestService_NewRequest(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		params := url.Values{}
		req, err := svc.NewRequest(http.MethodGet, uri, params, nil)

		assert.NoError(t, err)
		assert.Equal(t, req.URL.Scheme, "http")
		assert.Equal(t, req.URL.Host, "localhost:8080")
		assert.Equal(t, req.URL.Path, "/characters/1017100")
		assert.Equal(t, req.URL.RawQuery, params.Encode())
	})
}

func TestService_Transport(t *testing.T) {
	t.Run("http do fail", func(t *testing.T) {
		svc.Hostname = ""
		svc.HTTPClient = &http.Client{
			Transport: RoundTripFunc(func(req *http.Request) *http.Response {
				return nil
			}),
		}
		params := url.Values{}
		req, _ := svc.NewRequest(http.MethodGet, uri, params, nil)
		_, err := svc.Transport(req)
		assert.EqualError(t, err, "marvel request failed: Get \"http://localhost:8080/characters/1017100\": foo")
	})
}

func TestService_TransformStatusCode(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		cdw := domain.CharacterDataWrapper{
			Code: http.StatusOK,
		}

		code, err := svc.TransformStatusCode(cdw)

		assert.Equal(t, http.StatusOK, code)
		assert.NoError(t, err)
	})
	t.Run("invalidcredentials", func(t *testing.T) {
		cdw := domain.CharacterDataWrapper{
			Code:    "InvalidCredentials",
			Message: "InvalidCredentials",
		}

		code, err := svc.TransformStatusCode(cdw)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.EqualError(t, err, "internal server error")
	})
	t.Run("missingParameter", func(t *testing.T) {
		cdw := domain.CharacterDataWrapper{
			Code:    "MissingParameter",
			Message: "MissingParameter",
		}

		code, err := svc.TransformStatusCode(cdw)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.EqualError(t, err, "internal server error")
	})
	t.Run("404", func(t *testing.T) {
		cdw := domain.CharacterDataWrapper{
			Code:   float64(404),
			Status: "Not found",
		}

		code, err := svc.TransformStatusCode(cdw)

		assert.Equal(t, http.StatusOK, code)
		assert.NoError(t, err)
	})
	t.Run("409", func(t *testing.T) {
		cdw := domain.CharacterDataWrapper{
			Code:   float64(409),
			Status: "",
		}

		code, err := svc.TransformStatusCode(cdw)

		assert.Equal(t, http.StatusInternalServerError, code)
		assert.EqualError(t, err, "internal server error")
	})
}

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), errors.New("foo")
}
