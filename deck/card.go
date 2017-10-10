package deck

import (
	"math/big"
)

// Cards is the deck of cards for this game. It can be overriden to use any card deck
var Cards []string = []string{"ACE", "KING", "QUEEN"}

// card represents a single card in a deck. The members correspond with:
// p1cipher - face of card encrypted only with "player 1's" key
// p2cipher - face of card encrypted only with "player 2's" key
// bothCipher - face of card encrypted with both keys
// plain - face of card in plaintext
type card struct {
	P1cipher, P2cipher, BothCipher *big.Int
	plain                          string
}

func (c *card) Face() string {
	return c.plain
}
