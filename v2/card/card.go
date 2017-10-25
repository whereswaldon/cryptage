package card

import (
	//"fmt"
	//"github.com/pkg/errors"
	"github.com/sorribas/shamir3pass"
	"math/big"
)

// Card represents a single card with methods to
// find the face of the card and each permutation
// of its encrypted state.
type Card interface {
	Face() (string, error)
	Mine() (*big.Int, error)
	Theirs() (*big.Int, error)
	Both() (*big.Int, error)
}

// NewCard creates an entirely new card from the given face
// value and key. After this operation, the card will have
// both a Face() and a Mine() value, but Both() and Theirs()
// will results in errors because the card has not been
// encrypted by another party.
func NewCard(face string, myKey shamir3pass.Key) (Card, error) {
	return nil, nil
}

// CardFrom creates a card from the given big integer. This
// assumes that the provided integer is the encrypted value
// of the card provided by another player.
func CardFromTheirs(theirs *big.Int, myKey shamir3pass.Key) (Card, error) {
	return nil, nil
}

// CardFromBoth creats a card from the given big integer. This
// assumes that the provided integer is the encrypted value
// after both players have encrypted the card.
func CardFromBoth(both *big.Int, myKey shamir3pass.Key) (Card, error) {
	return nil, nil
}
