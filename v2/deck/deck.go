package deck

import (
	"math/big"
)

// CardHolder is an ordered collection of cards.
type CardHolder interface {
	Get(uint) (CardFace, error)
	SetBothEncrypted([]*big.Int) error
	GetAllMine() ([]*big.Int, bool, error)
}
