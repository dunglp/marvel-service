package http

import (
	"net/http"
	"xendit-technical-assessment/pkg/domain"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func marvelServiceRouter(
	service domain.MarvelService,
) *mux.Router {
	router := mux.NewRouter()
	chain := alice.New()
	router.Handle("/characters/{characterId}", chain.Then(marvelServiceHandler{service, CharacterDetailsHandler})).Methods(http.MethodGet)
	router.Handle("/characters", chain.Then(marvelServiceHandler{service, CharactersHandler})).Methods(http.MethodGet)

	return router
}
