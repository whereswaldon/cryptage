package card

import (
	"fmt"
	"github.com/sorribas/shamir3pass"
	"math/big"

	. "github.com/whereswaldon/cryptage/v2/types"
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
	if theirs == nil {
		return nil, fmt.Errorf("Unable to create card from nil encrypted value")
	} else if myKey == nil {
		return nil, fmt.Errorf("Unable to create card with nil key pointer")
	}
	return &card{myKey: myKey, theirs: theirs}, nil
}

// CardFromBoth creats a card from the given big integer. This
// assumes that the provided integer is the encrypted value
// after both players have encrypted the card.
func CardFromBoth(both *big.Int, myKey *shamir3pass.Key) (Card, error) {
	if both == nil {
		return nil, fmt.Errorf("Unable to create card from nil encrypted value")
	} else if myKey == nil {
		return nil, fmt.Errorf("Unable to create card with nil key pointer")
	}
	return &card{myKey: myKey, both: both}, nil
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

// Mine returns the card's face encrypted solely by the local player's key,
// if possible.
func (c *card) Mine() (*big.Int, error) {
	if c.mine != nil {
		return c.mine, nil
	} else if c.face != "" {
		c.mine = shamir3pass.Encrypt(big.NewInt(0).SetBytes([]byte(c.face)), *c.myKey)
		return c.mine, nil
	}
	return nil, fmt.Errorf("Unable to get card solely encrypted by local player: %v", c)
}

// Theirs returns the card's face encrypted solely by the opponent's key,
// if possible.
func (c *card) Theirs() (*big.Int, error) {
	if c.theirs != nil {
		return c.theirs, nil
	} else if c.both != nil {
		c.theirs = shamir3pass.Decrypt(c.both, *c.myKey)
		return c.theirs, nil
	}
	return nil, fmt.Errorf("Unable to get card solely encrypted with opponent's key: %v", c)
}

// Both returns the card's face encrypted with the keys of both players.
func (c *card) Both() (*big.Int, error) {
	if c.both != nil {
		return c.both, nil
	} else if c.theirs != nil {
		c.both = shamir3pass.Encrypt(c.theirs, *c.myKey)
		return c.both, nil
	}
	return nil, fmt.Errorf("Unable to get card encrypted with both player's keys: %v", c)
}

func (c *card) String() string {
	return fmt.Sprintf("mine: %v\ntheirs:%v\nboth:%v\nface:%s", c.mine, c.theirs, c.both, c.face)
}
