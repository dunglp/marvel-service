package domain

// swagger:model CharacterDataWrapper
type CharacterDataWrapper struct {
	Code      interface{}            `json:"code"`
	Message   string                 `json:"message"`
	Status    string                 `json:"status"`
	Copyright string                 `json:"copyright"`
	Data      CharacterDataContainer `json:"data"`
}

// swagger:model CharacterDataContainer
type CharacterDataContainer struct {
	Offset  int         `json:"offset"`
	Limit   int         `json:"limit"`
	Total   int         `json:"total"`
	Count   int         `json:"count"`
	Results []Character `json:"results"`
}

// A Character is an actor in marvel universe
// It is used to describe the actor in the marvel universe
//
// swagger:model Character
type Character struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// A CharacterID parameter model.
//
// This is used for operations that want the ID of a character in the path
// swagger:parameters getCharacterByID
type CharacterID struct {
	// The ID of the character
	//
	// in: path
	// required: true
	CharacterID int `json:"characterId"`
}
