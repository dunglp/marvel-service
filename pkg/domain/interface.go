package domain

type MarvelService interface {
	GetCharacterDetails(int) (Character, error)
	GetCharacterIds() ([]int, error)
}
