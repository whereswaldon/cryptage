package deck

import (
	"fmt"
	shamir "github.com/sorribas/shamir3pass"
	"io"
	"math/big"
)

type Deck interface {
	Draw() (string, error)
	Start() error
	Play()
	String() string
}

// ensure that *deck fulfills Deck interface
var _ Deck = &deck{}

type deck struct {
	keys     shamir.Key
	cards    []card
	protocol *Protocol
	done     chan struct{}
}

// Draw draws a single card from the deck
func (d *deck) Draw() (string, error) {
	d.protocol.RequestDecryptCard(0)
	fmt.Printf("Requesting decryption of card:\n%v\n", d.cards[0].BothCipher)
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
	go d.handleMessages()
}

// initEncryptCards  encrypts each card from the
// plaintext to a player1 ciphertext
func (d *deck) initEncryptCards() {
	for i := range d.cards {
		d.cards[i].MyCipher = shamir.Encrypt(big.NewInt(0).SetBytes([]byte(Cards[i])), d.keys)
		fmt.Printf("Layer1: %d: %s -> %v -> %s\n", i, Cards[i],
			d.cards[i].MyCipher, string(shamir.Decrypt(d.cards[i].MyCipher, d.keys).Bytes()))
	}
	//	fmt.Println(d.cards)
}

// encryptCards takes a deck with player1 ciphertext populated and
// encrypts that ciphertext again to arrive at both players'
// ciphertext.
func (d *deck) encryptCards() {
	for i, c := range d.cards {
		d.cards[i].BothCipher = shamir.Encrypt(c.MyCipher, d.keys)
		fmt.Printf("Layer2: %d: %v -> %v -> %v\n", i, c.MyCipher,
			d.cards[i].BothCipher, shamir.Decrypt(d.cards[i].BothCipher, d.keys))
	}
	//	fmt.Println(d.cards)
}

// setMyCiphers sets the ciphertext of the deck to the provided
// array
func (d *deck) setMyCiphers(ciphers []*big.Int) {
	for i, c := range ciphers {
		d.cards[i].MyCipher = c
	}
}

// clearMyCiphers erases old ciphertext
func (d *deck) clearDeck() {
	for i := range d.cards {
		d.cards[i].MyCipher = nil
		d.cards[i].TheirCipher = nil
		d.cards[i].BothCipher = nil
		d.cards[i].plain = ""
	}
}

// decryptCard removes the current player's encryption, leaving
// only the other player's encryption layer.
func (d *deck) decryptCard(index uint64) {
	if d.cards[index].TheirCipher != nil {
		return
	}
	d.cards[index].TheirCipher = shamir.Decrypt(d.cards[index].BothCipher, d.keys)
}

// revealCard removes the current player's encryption from a card,
// revealing the plaintext face of the card.
func (d *deck) revealCard(index uint64) {
	if d.cards[index].plain != "" {
		return
	}

	//todo check whether the p2cipher is null
	plainBigInt := shamir.Decrypt(d.cards[index].MyCipher, d.keys)
	d.cards[index].plain = string(plainBigInt.Bytes())
	fmt.Println("Revealed: ", d.cards[index].plain, d.cards[index].MyCipher)
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
			d.keys = shamir.GenerateKeyFromPrime(msg.Value)
			d.setMyCiphers(msg.Deck)
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
		case DECRYPT_CARD:
			fmt.Println("DECRYPT_CARD")
			d.decryptCard(msg.Index)
			fmt.Println("decrypted", d.cards[msg.Index].TheirCipher)
			d.protocol.SendDecryptedCard(msg.Index, d.cards[msg.Index].TheirCipher)
		case ONE_CIPHER_CARD:
			fmt.Println("ONE_CIPHER_CARD")
			d.cards[msg.Index].MyCipher = msg.Value
			d.revealCard(msg.Index)
		default:
			fmt.Println("Unknown message: %v", msg)
		}
	}
}

// Start runs the game, and initiates the first hand
func (d *deck) Start() error {
	prime := shamir.Random1024BitPrime()
	d.keys = shamir.GenerateKeyFromPrime(prime)
	d.initEncryptCards()
	if err := d.protocol.SendStartDeck(prime, d.cards); err != nil {
		return err
	}
	//	if err := d.protocol.SendQuit(); err != nil {
	//		fmt.Println(err)
	//	}

	go d.handleMessages()
	return nil
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
	d.cards = make([]card, len(Cards))

	return d, nil
}
