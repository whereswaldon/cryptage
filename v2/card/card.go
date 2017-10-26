package card

import (
	"fmt"
	"github.com/sorribas/shamir3pass"
	"math/big"
)

type CardFace string

const EMPTY_CARD = ""

func EncryptString(s CardFace, k *shamir3pass.Key) *big.Int {
	return shamir3pass.Encrypt(big.NewInt(0).SetBytes([]byte(s)), *k)
}

func DecryptString(i *big.Int, k *shamir3pass.Key) CardFace {
	return CardFace(shamir3pass.Decrypt(i, *k).Bytes())
}

// NewCard creates an entirely new Card from the given face
// value and key. After this operation, the Card will have
// both a Face() and a Mine() value, but Both() and Theirs()
// will results in errors because the Card has not been
// encrypted by another party.
func NewCard(face CardFace, myKey *shamir3pass.Key) (*Card, error) {
	if face == "" {
		return nil, fmt.Errorf("Unable to create Card with empty string as face")
	} else if myKey == nil {
		return nil, fmt.Errorf("Unable to create Card with nil key pointer")
	}
	return &Card{myKey: myKey, face: face}, nil
}

// CardFrom creates a Card from the given big integer. This
// assumes that the provided integer is the encrypted value
// of the Card provided by another player.
func CardFromTheirs(theirs *big.Int, myKey *shamir3pass.Key) (*Card, error) {
	if theirs == nil {
		return nil, fmt.Errorf("Unable to create Card from nil encrypted value")
	} else if myKey == nil {
		return nil, fmt.Errorf("Unable to create Card with nil key pointer")
	}
	return &Card{myKey: myKey, theirs: theirs}, nil
}

// CardFromBoth creats a Card from the given big integer. This
// assumes that the provided integer is the encrypted value
// after both players have encrypted the Card.
func CardFromBoth(both *big.Int, myKey *shamir3pass.Key) (*Card, error) {
	if both == nil {
		return nil, fmt.Errorf("Unable to create Card from nil encrypted value")
	} else if myKey == nil {
		return nil, fmt.Errorf("Unable to create Card with nil key pointer")
	}
	return &Card{myKey: myKey, both: both}, nil
}

type Card struct {
	myKey, theirKey    *shamir3pass.Key
	face               CardFace
	mine, theirs, both *big.Int
}

// Face returns the face of the Card if it is known or can be computed locally.
// If neither the face nor mine fields are populated, the opponent must consent
// to decrypt the Card, which is handled elsewhere.
func (c *Card) Face() (CardFace, error) {
	if c.face != "" {
		return c.face, nil
	} else if c.mine != nil {
		c.face = DecryptString(c.mine, c.myKey)
		return c.face, nil
	}
	return "", fmt.Errorf("Unable to view Card face, need other player to decrypt Card: %v", c)
}

// Mine returns the Card's face encrypted solely by the local player's key,
// if possible.
func (c *Card) Mine() (*big.Int, error) {
	if c.mine != nil {
		return c.mine, nil
	} else if c.face != "" {
		c.mine = EncryptString(c.face, c.myKey)
		return c.mine, nil
	}
	return nil, fmt.Errorf("Unable to get Card solely encrypted by local player: %v", c)
}

// Theirs returns the Card's face encrypted solely by the opponent's key,
// if possible.
func (c *Card) Theirs() (*big.Int, error) {
	if c.theirs != nil {
		return c.theirs, nil
	} else if c.both != nil {
		c.theirs = shamir3pass.Decrypt(c.both, *c.myKey)
		return c.theirs, nil
	}
	return nil, fmt.Errorf("Unable to get Card solely encrypted with opponent's key: %v", c)
}

// Both returns the Card's face encrypted with the keys of both players.
func (c *Card) Both() (*big.Int, error) {
	if c.both != nil {
		return c.both, nil
	} else if c.theirs != nil {
		c.both = shamir3pass.Encrypt(c.theirs, *c.myKey)
		return c.both, nil
	}
	return nil, fmt.Errorf("Unable to get Card encrypted with both player's keys: %v", c)
}

// SetMine gives the Card the value of a Card face encrypted solely with
// the current player's key. This value is set and trusted implicitly,
// but can be validated if you use SetTheirKey and Validate on the Card
// at a later time.
func (c *Card) SetMine(mine *big.Int) error {
	if mine == nil {
		return fmt.Errorf("Cannot set nil as mine value")
	}
	c.mine = mine
	return nil
}

// SetTheirKey gives the Card the key that the opponent used to perform
// their side of the encryption. It returns and error if the key is invalid.
func (c *Card) SetTheirKey(theirKey *shamir3pass.Key) error {
	if theirKey == nil {
		return fmt.Errorf("Cannot set theirKey as nil")
	}
	c.theirKey = theirKey
	return nil
}

func (c *Card) HasTheirKey() bool {
	return c.theirKey != nil
}

// CanDecrypt returns whether the Card is able to be decrypted in its
// current state.
func (c *Card) CanDecrypt() bool {
	return c.face != "" || c.mine != nil ||
		(c.HasTheirKey() && (c.both != nil || c.theirs != nil))
}

// Validate checks the Card's internal consistency. In order to be called,
// mykey, theirKey, and both need to be set. It will not return an error
// if the Card is internally consistent.
func (c *Card) Validate() error {
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

func (c *Card) String() string {
	return fmt.Sprintf("mine: %v\ntheirs:%v\nboth:%v\nface:%s", c.mine, c.theirs, c.both, c.face)
}
