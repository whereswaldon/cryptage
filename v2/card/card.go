package card

import (
	"fmt"
	//"github.com/pkg/errors"
	"github.com/sorribas/shamir3pass"
	. "github.com/whereswaldon/cryptage/v2/types"
	"math/big"
)

// NewCard creates an entirely new card from the given face
// value and key. After this operation, the card will have
// both a Face() and a Mine() value, but Both() and Theirs()
// will results in errors because the card has not been
// encrypted by another party.
func NewCard(face string, myKey *shamir3pass.Key) (Card, error) {
	if face == "" {
		return nil, fmt.Errorf("Unable to create card with empty string as face")
	} else if myKey == nil {
		return nil, fmt.Errorf("Unable to create card with nil key pointer")
	}
	return &card{myKey: myKey, face: face}, nil
}

// CardFrom creates a card from the given big integer. This
// assumes that the provided integer is the encrypted value
// of the card provided by another player.
func CardFromTheirs(theirs *big.Int, myKey *shamir3pass.Key) (Card, error) {
	return nil, nil
}

// CardFromBoth creats a card from the given big integer. This
// assumes that the provided integer is the encrypted value
// after both players have encrypted the card.
func CardFromBoth(both *big.Int, myKey *shamir3pass.Key) (Card, error) {
	return nil, nil
}

type card struct {
	myKey              *shamir3pass.Key
	face               string
	mine, theirs, both *big.Int
}

// ensure that the card type always satisfies the Card interface
var _ Card = &card{}

// Face returns the face of the card if it is known or can be computed locally.
// If neither the face nor mine fields are populated, the opponent must consent
// to decrypt the card, which is handled elsewhere.
func (c *card) Face() (string, error) {
	if c.face != "" {
		return c.face, nil
	} else if c.mine != nil {
		c.face = string(shamir3pass.Decrypt(c.mine, *c.myKey).Bytes())
		return c.face, nil
	}
	return "", fmt.Errorf("Unable to view card face, need other player to decrypt card: %v", c)
}
func (c *card) Mine() (*big.Int, error) {
	if c.mine != nil {
		return c.mine, nil
	}
	return nil, fmt.Errorf("Unable to get card solely encrypted by local player: %v", c)
}
func (c *card) Theirs() (*big.Int, error) {
	return nil, nil
}
func (c *card) Both() (*big.Int, error) {
	return nil, nil
}
