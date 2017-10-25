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
	Validate() error
}

type CardHolder interface {
	SetFaces([]CardFace)
	Get(uint) CardFace
	SetTheirEncrypted([]*big.Int)
	SetBothEncrypted([]*big.Int)
	Shuffle()
}
