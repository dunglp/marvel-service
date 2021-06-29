package marvel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"xendit-technical-assessment/pkg/domain"
	cache2 "xendit-technical-assessment/pkg/marvel/cache"

	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

const (
	charactersEndpoint       = "/v1/public/characters"
	characterDetailsEndpoint = "/v1/public/characters/%d"
	charactersReqLimit       = 100
	marvelErrorInvalidCredentials = "InvalidCredentials"
	marvelErrorMissingParameter = "MissingParameter"
)

type Service struct {
	HTTPClient *http.Client
	Hostname   string
	Cred       domain.Credentials
	Cache      cache2.CachedMarvelService
}

func NewService(hostname string, cred domain.Credentials, cache *cache.Cache) Service {
	return Service{
		HTTPClient: &http.Client{},
		Hostname:   hostname,
		Cred:       cred,
		Cache:      cache2.NewCachedMarvelService(cache),
	}
}

func (s Service) GetCharacterDetails(id int) (character domain.Character, err error) {
	params, err := s.Cred.ToBaseMarvelRequestParams()
	if err != nil {
		return
	}
	apiURL := s.Hostname + fmt.Sprintf(characterDetailsEndpoint, id)
	req, err := s.NewRequest(http.MethodGet, apiURL, params, nil)
	if err != nil {
		return
	}

	resp, err := s.Transport(req)
	if err != nil {
		return
	}

	code, err := s.TransformStatusCode(resp)
	if err != nil {
		return
	}

	if code != http.StatusNotFound && len(resp.Data.Results) > 0 {
		character = resp.Data.Results[0]
	}
	return
}

func (s Service) GetCharacterIds() (ids []int, err error) {
	log.Info().Msg("Lookup from cache")
	found, cacheData := s.Cache.CharacterIds()
	var since time.Time
	total := 0

	if !found {
		log.Info().Msg("Not found in cache. Get from Marvel")
		ids, total, err = s.GetCharacterIdsFrom(since)
		if err != nil {
			return
		}
		s.Cache.CacheCharacterIds(ids, total, time.Now())
		return
	}

	since = cacheData.Since
	log.Info().Msg("Found in cache. Get modified characters from Marvel")
	modifiedIds, _, err := s.GetCharacterIdsFrom(since)
	if err != nil {
		return
	}

	ids = s.Cache.UpdateCache(modifiedIds, time.Now())
	return
}

func (s Service) GetCharacterIdsFrom(since time.Time) (ids []int, total int, err error) {
	offset := 0

	for {
		cdc, err := s.GetMarvelCharacters(charactersReqLimit, offset, since)
		if err != nil {
			return ids, total, err
		}
		total = cdc.Total

		for _, c := range cdc.Results {
			ids = append(ids, c.ID)
		}
		if offset >= total-charactersReqLimit {
			break
		}
		offset += charactersReqLimit
	}

	return
}

func (s Service) GetMarvelCharacters(limit int, offset int, since time.Time) (cdc domain.CharacterDataContainer, err error) {
	params, err := s.Cred.ToCharacterIDsReqParams(limit, offset, since)
	if err != nil {
		return
	}

	apiURL := s.Hostname + charactersEndpoint

	req, err := s.NewRequest(http.MethodGet, apiURL, params, nil)
	if err != nil {
		return
	}

	resp, err := s.Transport(req)
	if err != nil {
		return
	}

	code, err := s.TransformStatusCode(resp)
	if err != nil {
		return
	}

	if code != http.StatusNotFound {
		cdc = resp.Data
	}
	return
}

func (s Service) NewRequest(method string, url string, params url.Values, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)

	if err != nil {
		return nil, fmt.Errorf("new request failed: %w", err)
	}

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	return
}

func (s Service) Transport(req *http.Request) (cdw domain.CharacterDataWrapper, err error) {

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return cdw, fmt.Errorf("marvel request failed: %w", err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	err = json.NewDecoder(resp.Body).Decode(&cdw)
	if err != nil {
		return cdw, fmt.Errorf("read body: %v", err)
	}

	return
}

func (s Service) TransformStatusCode(c domain.CharacterDataWrapper) (int, error) {
	switch c.Code {
	case marvelErrorInvalidCredentials, marvelErrorMissingParameter:
		log.Error().Msg(fmt.Sprintf("%s: %s", c.Code, c.Message))
		return http.StatusInternalServerError, errors.New("internal server error")
	case float64(404):
		code := c.Code.(float64)
		log.Error().Msg(fmt.Sprintf("%d: %s", int(code), c.Status))
		return http.StatusOK, nil
	case float64(409):
		code := c.Code.(float64)
		log.Error().Msg(fmt.Sprintf("%d: %s", int(code), c.Status))
		return http.StatusInternalServerError, errors.New("internal server error")
	}
	return http.StatusOK, nil
}
