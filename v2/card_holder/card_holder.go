package card_holder

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sorribas/shamir3pass"
	"github.com/whereswaldon/cryptage/v2/card"
	"math/big"
)

// Card represents a single card with methods to
// find the face of the card and each permutation
// of its encrypted state.
type Card interface {
	Face() (card.CardFace, error)
	Mine() (*big.Int, error)
	Theirs() (*big.Int, error)
	Both() (*big.Int, error)
	SetMine(*big.Int) error
	SetTheirKey(*shamir3pass.Key) error
	HasTheirKey() bool
	CanDecrypt() bool
	Validate() error
}

// NewHolder creates a CardHolder from a key and a slice of card faces
func NewHolder(key *shamir3pass.Key, faces []card.CardFace) (*CardHolder, error) {
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
	return &CardHolder{cards: cards, key: key}, nil
}

// HolderFromEncrypted creates a CardHolder from a kay and a slice of card
// faces that have already been encrypted.
func HolderFromEncrypted(key *shamir3pass.Key, theirEncrypted []*big.Int) (*CardHolder, error) {
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
	return &CardHolder{cards: cards, key: key}, nil
}

type CardHolder struct {
	cards []Card
	key   *shamir3pass.Key
}

// Get returns the card face at the given position within the deck,
// if possible
func (h *CardHolder) Get(index uint) (card.CardFace, error) {
	if can, err := h.CanGet(index); !can {
		return "", errors.Wrapf(err, "Unable to get card")
	}
	return "", nil
}

// CanGet determines whether it is currently possible to get the card
// face at the given position.
func (h *CardHolder) CanGet(index uint) (bool, error) {
	if index < uint(len(h.cards)) {
		return h.cards[index].CanDecrypt(), nil
	}
	return false, fmt.Errorf("Index %d out of bounds (%d cards)", index, len(h.cards))
}

// SetMine informs the card at the given index that its encrypted state
// with only the local player's key is the given big.Int. Knowing this
// value allows a card to be decrypted.
func (h *CardHolder) SetMine(index uint, mine *big.Int) error {
	return nil
}

// SetBothEncrypted erases the current cards within the deck and creates
// a new deck of cards with the given values as the values that have been
// encrypted by both players. This does not erase the encryption key stored
// within the collection (which is important, since this is needed to decrypt
// cards later)
func (h *CardHolder) SetBothEncrypted(encryptedFaces []*big.Int) error {
	if encryptedFaces == nil {
		return fmt.Errorf("Cannot set both encrypted to nil slice")
	} else if len(encryptedFaces) != len(h.cards) {
		return fmt.Errorf("Cannot set both encrypted to slice of length %d"+
			" when cardCardHolder has %d cards", len(encryptedFaces),
			len(h.cards))
	}
	var err error
	newCards := make([]Card, len(h.cards))
	for i, c := range encryptedFaces {
		newCards[i], err = card.CardFromBoth(c, h.key)
		if err != nil {
			return errors.Wrap(err, "Unable to create card from both encrypted face")
		}
	}
	h.cards = newCards
	return nil
}

// GetAllMine returns all known mine values for the cards. If all cards had known
// values, the second return value will be true.
func (h *CardHolder) GetAllMine() ([]*big.Int, bool, error) {
	cards := make([]*big.Int, len(h.cards))
	allDecryptable := true
	var err error
	for i, c := range h.cards {
		if c.CanDecrypt() {
			cards[i], err = c.Mine()
			if err != nil {
				return cards, false, err
			}
		} else {
			allDecryptable = false
		}
	}
	return cards, allDecryptable, nil
}
