package types

import (
	"github.com/sorribas/shamir3pass"
	"math/big"
)

type CardFace string

// Card represents a single card with methods to
// find the face of the card and each permutation
// of its encrypted state.
type Card interface {
	Face() (CardFace, error)
	Mine() (*big.Int, error)
	Theirs() (*big.Int, error)
	Both() (*big.Int, error)
	SetMine(*big.Int) error
	SetTheirKey(*shamir3pass.Key) error
	HasTheirKey() bool
	CanDecrypt() bool
	Validate() error
}

// CardHolder is an ordered collection of cards.
type CardHolder interface {
	Get(uint) (CardFace, error)
	SetBothEncrypted([]*big.Int) error
	GetAllMine() ([]*big.Int, bool, error)
}
