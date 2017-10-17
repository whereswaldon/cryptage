package deck

import (
	"fmt"
	shamir "github.com/sorribas/shamir3pass"
	"io"
	"math/big"
)

type Deck interface {
	Draw() (string, error)
	Start()
	Play()
	String() string
}

// ensure that *deck fulfills Deck interface
var _ Deck = &deck{}

type deck struct {
	keys              shamir.Key
	cards             []card
	topCard		int
	protocol          *Protocol
	playerNum         int8
	currentPendingReq chan *big.Int
	done              chan struct{}
}

// Draw draws a single card from the deck
func (d *deck) Draw() (string, error) {
    	d.currentPendingReq = make(chan *big.Int)
	if err := d.protocol.SendReveal(topCard); err != nil {
    		return "", err
	}
	return "", nil
}

func (d *deck) String() string {
	return "Deck"
}

func (d *deck) Quit() {
	close(d.done)
}

// Play runs the game, but does not initiate it.
func (d *deck) Play() {
	d.playerNum = 2
	d.handleMessages()
}

// initEncryptCards  encrypts each card from the
// plaintext to a player1 ciphertext
func (d *deck) initEncryptCards() {
	for i := range d.cards {
		d.cards[i].P1cipher = shamir.Encrypt(big.NewInt(0).SetBytes([]byte(Cards[i])), d.keys)
	}
	fmt.Println(d.cards)
}

// encryptCards takes a deck with player1 ciphertext populated and
// encrypts that ciphertext again to arrive at both players'
// ciphertext.
func (d *deck) encryptCards() {
	for i, c := range d.cards {
		d.cards[i].BothCipher = shamir.Encrypt(c.P1cipher, d.keys)
	}
	fmt.Println(d.cards)
}

// setPlayer1Ciphers sets the ciphertext of the deck to the provided
// array
func (d *deck) setPlayer1Ciphers(ciphers []*big.Int) {
	for i, c := range ciphers {
		d.cards[i].P1cipher = c
	}
}

// clearPlayer1Ciphers erases old ciphertext
func (d *deck) clearDeck() {
	for i := range d.cards {
		d.cards[i].P1cipher = nil
		d.cards[i].P2cipher = nil
		d.cards[i].BothCipher = nil
		d.cards[i].plain = ""
	}
}

// setBothCiphers sets the ciphertext of the deck to the provided
// array
func (d *deck) setBothCiphers(ciphers []*big.Int) {
	for i, c := range ciphers {
		d.cards[i].BothCipher = c
	}
}
func (d *deck) handleMessages() {
	defer d.Quit()
	for msg := range d.protocol.Listen() {
		switch msg.Type {
		case QUIT:
			fmt.Println("QUIT")
			return
		case START_DECK:
			fmt.Println("START_DECK")
			d.setPlayer1Ciphers(msg.Deck)
			d.encryptCards()
			d.protocol.SendEndDeck(d.cards)
		case END_DECK:
			fmt.Println("END_DECK")
			// since the ciphertext has been shuffled, we no longer
			// know which card is which. All of our old data is now
			// irrelevant
			d.clearDeck()
			d.setBothCiphers(msg.Deck)
			fmt.Println(d.cards)
		case REVEAL_CARD:
    			fmt.Println("REVEAL_CARD: ", msg.Index)

		default:
			fmt.Println("Unknown message: %v", msg)
		}
	}
}

// Start runs the game, and initiates the first hand
func (d *deck) Start() {
	d.initEncryptCards()
	d.playerNum = 1
	if err := d.protocol.SendStartDeck(d.cards); err != nil {
		fmt.Println(err)
	}
	if err := d.protocol.SendQuit(); err != nil {
		fmt.Println(err)
	}

	d.handleMessages()
}

// NewDeck creates a deck of cards and assumes that the given
// io.ReadWriteCloser is a connection of some sort to another
// deck.
func NewDeck(deckConnection io.ReadWriteCloser) (Deck, error) {
	d := &deck{}
	done := make(chan struct{})
	p, err := NewProtocol(deckConnection, done)
	if err != nil {
		return nil, err
	}
	d.done = done
	d.protocol = p
	d.keys = shamir.GenerateKey(128)
	d.cards = make([]card, len(Cards))
	d.topCard = 0

	return d, nil
}
