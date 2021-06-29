package http

import (
	"net/http"
	"time"
	"xendit-technical-assessment/pkg/domain"

	"go.opencensus.io/plugin/ochttp"
)

func NewMarvelServiceHTTPServer(
	address string,
	service domain.MarvelService,
) *http.Server {
	return &http.Server{
		Addr:         address,
		Handler:      &ochttp.Handler{Handler: marvelServiceRouter(service)},
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 200 * time.Second,
	}
}
