package http

import (
	"net/http"
	"strconv"
	"xendit-technical-assessment/pkg/domain"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type marvelServiceHandler struct {
	service domain.MarvelService
	h       func(domain.MarvelService, http.ResponseWriter, *http.Request)
}

func (msh marvelServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3003")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	msh.h(msh.service, w, r)
}

// CharacterDetailsHandler swagger:route GET /characters/{characterId} characters getCharacterByID
//
// Get Character Details
// Responses:
// 		default: response
func CharacterDetailsHandler(service domain.MarvelService, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["characterId"])

	if err != nil {
		domain.JSON(w, http.StatusInternalServerError, domain.Response{
			Status:  "Bad Request",
			Code:    http.StatusBadRequest,
			Message: "Invalid Character ID ",
		})
		return
	}

	log.Info().Msgf("Start getting details for character %d", id)
	character, err := service.GetCharacterDetails(id)
	if err != nil {
		log.Err(err).Msgf("Failed getting character details %d", id)
		domain.JSON(w, http.StatusInternalServerError, domain.Response{
			Status:  "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	log.Info().Msgf("Character details : %+v", character)
	if character.ID == 0 {
		domain.JSON(w, http.StatusNotFound, domain.Response{
			Status: "Character not found",
			Code:   http.StatusNotFound,
		})
		return
	}

	domain.JSON(w, http.StatusOK, domain.Response{
		Status: "OK",
		Code:   http.StatusOK,
		Data:   character,
	})
}

// CharactersHandler swagger:route GET /characters characters getCharacterIDs
//
// Get Character IDs
// Responses:
// 		default: response
func CharactersHandler(service domain.MarvelService, w http.ResponseWriter, r *http.Request) {

	log.Info().Msgf("Start getting character IDs list")
	ids, err := service.GetCharacterIds()
	if err != nil {
		log.Err(err).Msg("Error getting character IDs list")
		domain.JSON(w, http.StatusInternalServerError, domain.Response{
			Status:  "Internal Server Error",
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	log.Info().Msgf("Character IDs : %v", ids)

	if len(ids) == 0 {
		domain.JSON(w, http.StatusNotFound, domain.Response{
			Status: "Character IDs Not found",
			Code:   http.StatusNotFound,
		})
		return
	}
	domain.JSON(w, http.StatusOK, domain.Response{
		Status: "OK",
		Code:   http.StatusOK,
		Data:   ids,
	})
}
