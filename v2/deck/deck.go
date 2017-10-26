package deck

import (
	"github.com/whereswaldon/cryptage/v2/card"
	"math/big"
)

// CardHolder is an ordered collection of cards.
type CardHolder interface {
	Get(uint) (card.CardFace, error)
	SetBothEncrypted([]*big.Int) error
	GetAllMine() ([]*big.Int, bool, error)
}
