package domain

import (
	"net/url"
	"strconv"
	"time"
	"xendit-technical-assessment/pkg/utils"
)

const (
	marvelParamTs     = "ts"
	marvelParamApiKey = "apikey"
	marvelParamHash   = "hash"
	marvelParamLimit  = "limit"
	marvelParamOffset   = "offset"
	marvelParamModifiedSince   = "modifiedSince"
	marvelParamModifiedSinceFormat = "2006-01-02T15:04:05-0700"
)

type Credentials struct {
	PublicKey  string
	PrivateKey string
	Ts         string
	Hashed     string
}

func (c Credentials) ToBaseMarvelRequestParams() (params url.Values, err error) {
	params = url.Values{}
	params.Add(marvelParamTs, c.Ts)
	params.Add(marvelParamApiKey, c.PublicKey)
	params.Add(marvelParamHash, c.Hashed)

	return
}

func (c Credentials) ToCharacterIDsReqParams(limit int, offset int, since time.Time) (params url.Values, err error) {
	params, err = c.ToBaseMarvelRequestParams()
	if err != nil {
		return
	}
	params.Add(marvelParamLimit, strconv.Itoa(limit))
	params.Add(marvelParamOffset, strconv.Itoa(offset))
	if !since.IsZero() {
		params.Add(marvelParamModifiedSince, since.Format(marvelParamModifiedSinceFormat))
	}
	return
}

func NewCredentials(publicKey, privateKey, ts string) (cred Credentials, err error) {
	cred = Credentials{
		publicKey, privateKey, ts, "",
	}
	hashed, err := toMd5Hash(cred)
	if err != nil {
		return
	}
	cred.Hashed = hashed
	return
}

func toMd5Hash(c Credentials) (string, error) {
	return utils.HashMd5(c.Ts + c.PrivateKey + c.PublicKey)
}
