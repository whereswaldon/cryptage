package types

import (
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
	SetMine(*big.Int) error
	SetTheirKey(*shamir3pass.Key) error
	Validate() error
}
