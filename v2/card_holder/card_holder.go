package card_holder

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sorribas/shamir3pass"
	"github.com/whereswaldon/cryptage/v2/card"
	"math/big"

	. "github.com/whereswaldon/cryptage/v2/types"
)

// NewHolder creates a CardHolder from a key and a slice of card faces
func NewHolder(key *shamir3pass.Key, faces []CardFace) (CardHolder, error) {
	if key == nil {
		return nil, fmt.Errorf("Cannot construct CardHolder with nil key")
	} else if faces == nil {
		return nil, fmt.Errorf("Cannot construct CardHolder with nil card faces")
	} else if len(faces) < 1 {
		return nil, fmt.Errorf("Cannot construct CardHolder with empty card faces")
	}
	cards := make([]Card, len(faces))
	var err error
	for i, c := range faces {
		cards[i], err = card.NewCard(c, key)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating card")
		}
	}
	return &holder{cards: cards, key: key}, nil
}

// HolderFromEncrypted creates a CardHolder from a kay and a slice of card
// faces that have already been encrypted.
func HolderFromEncrypted(key *shamir3pass.Key, theirEncrypted []*big.Int) (CardHolder, error) {
	if key == nil {
		return nil, fmt.Errorf("Cannot construct CardHolder with nil key")
	} else if theirEncrypted == nil {
		return nil, fmt.Errorf("Cannot construct CardHolder with nil encrypted faces")
	} else if len(theirEncrypted) < 1 {
		return nil, fmt.Errorf("Cannot construct CardHolder with empty encrypted faces")
	}
	cards := make([]Card, len(theirEncrypted))
	var err error
	for i, c := range theirEncrypted {
		cards[i], err = card.CardFromTheirs(c, key)
		if err != nil {
			return nil, errors.Wrap(err, "Error creating card")
		}
	}
	return &holder{cards: cards, key: key}, nil
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
