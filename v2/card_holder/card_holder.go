package card_holder

import (
	"github.com/sorribas/shamir3pass"
	"math/big"

	. "github.com/whereswaldon/cryptage/v2/types"
)

// NewHolder creates a CardHolder from a key and a slice of card faces
func NewHolder(key *shamir3pass.Key, faces []CardFace) (CardHolder, error) {
	return nil, nil
}

// HolderFromEncrypted creates a CardHolder from a kay and a slice of card
// faces that have already been encrypted.
func HolderFromEncrypted(key *shamir3pass.Key, theirEncrypted []*big.Int) (CardHolder, error) {
	return nil, nil
}

type holder struct {
	cards []Card
	key   *shamir3pass.Key
}

// ensure that holder implements CardHolder
var _ CardHolder = &holder{}

// Get returns the card face at the given position within the deck,
// if possible
func (h *holder) Get(index uint) (CardFace, error) {
	return "", nil
}

// SetBothEncrypted erases the current cards within the deck and creates
// a new deck of cards with the given values as the values that have been
// encrypted by both players. This does not erase the encryption key stored
// within the collection (which is important, since this is needed to decrypt
// cards later)
func (h *holder) SetBothEncrypted(encryptedFaces []*big.Int) error {
	return nil
}
