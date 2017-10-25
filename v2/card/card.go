package card

import (
	"fmt"
	"github.com/sorribas/shamir3pass"
	"math/big"

	. "github.com/whereswaldon/cryptage/v2/types"
)

func EncryptString(s CardFace, k *shamir3pass.Key) *big.Int {
	return shamir3pass.Encrypt(big.NewInt(0).SetBytes([]byte(s)), *k)
}

func DecryptString(i *big.Int, k *shamir3pass.Key) CardFace {
	return CardFace(shamir3pass.Decrypt(i, *k).Bytes())
}

// NewCard creates an entirely new card from the given face
// value and key. After this operation, the card will have
// both a Face() and a Mine() value, but Both() and Theirs()
// will results in errors because the card has not been
// encrypted by another party.
func NewCard(face CardFace, myKey *shamir3pass.Key) (Card, error) {
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
	myKey, theirKey    *shamir3pass.Key
	face               CardFace
	mine, theirs, both *big.Int
}

// ensure that the card type always satisfies the Card interface
var _ Card = &card{}

// Face returns the face of the card if it is known or can be computed locally.
// If neither the face nor mine fields are populated, the opponent must consent
// to decrypt the card, which is handled elsewhere.
func (c *card) Face() (CardFace, error) {
	if c.face != "" {
		return c.face, nil
	} else if c.mine != nil {
		c.face = DecryptString(c.mine, c.myKey)
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
		c.mine = EncryptString(c.face, c.myKey)
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

// SetMine gives the card the value of a card face encrypted solely with
// the current player's key. This value is set and trusted implicitly,
// but can be validated if you use SetTheirKey and Validate on the card
// at a later time.
func (c *card) SetMine(mine *big.Int) error {
	if mine == nil {
		return fmt.Errorf("Cannot set nil as mine value")
	}
	c.mine = mine
	return nil
}

// SetTheirKey gives the card the key that the opponent used to perform
// their side of the encryption. It returns and error if the key is invalid.
func (c *card) SetTheirKey(theirKey *shamir3pass.Key) error {
	if theirKey == nil {
		return fmt.Errorf("Cannot set theirKey as nil")
	}
	c.theirKey = theirKey
	return nil
}

func (c *card) HasTheirKey() bool {
	return c.theirKey != nil
}

// Validate checks the card's internal consistency. In order to be called,
// mykey, theirKey, and both need to be set. It will not return an error
// if the card is internally consistent.
func (c *card) Validate() error {
	if c.both == nil {
		return fmt.Errorf("Missing required field both")
	} else if c.myKey == nil {
		return fmt.Errorf("Missing required field myKey")
	} else if c.theirKey == nil {
		return fmt.Errorf("Missing required field theirKey")
	}
	theirs := shamir3pass.Decrypt(c.both, *c.myKey)
	if c.theirs != nil && c.theirs.Cmp(theirs) != 0 {
		return fmt.Errorf("Decrypted value mismatch: stored theirs: %v, computed theirs: %v", c.theirs, theirs)
	}
	mine := shamir3pass.Decrypt(c.both, *c.theirKey)
	if c.mine != nil && c.mine.Cmp(mine) != 0 {
		return fmt.Errorf("Decrypted value mismatch: stored mine: %v, computed mine: %v", c.theirs, theirs)
	}
	face := DecryptString(mine, c.myKey)
	if c.face != "" && c.face != face {
		return fmt.Errorf("Decrypted faces do not match: stored face: %s, computed face: %s", c.face, face)
	}
	return nil
}

func (c *card) String() string {
	return fmt.Sprintf("mine: %v\ntheirs:%v\nboth:%v\nface:%s", c.mine, c.theirs, c.both, c.face)
}
