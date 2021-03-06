package deck

import (
	"fmt"
	"github.com/sorribas/shamir3pass"
	"github.com/whereswaldon/cryptage/card"
	"github.com/whereswaldon/cryptage/card_holder"
	p "github.com/whereswaldon/cryptage/protocol"
	"io"
	"log"
	"math/big"
	"time"
)

// CardHolder is an ordered collection of cards.
type CardHolder interface {
	Size() uint
	CanGet(index uint) (bool, error)
	CanGetTheirs(index uint) (bool, error)
	Get(index uint) (card.CardFace, error)
	GetAllBoth() ([]*big.Int, bool, error)
	GetAllMine() ([]*big.Int, bool, error)
	GetTheirs(index uint) (*big.Int, error)
	SetBothEncrypted(encryptedFaces []*big.Int) error
	SetMine(index uint, mine *big.Int) error
	SetTheirKey(key *shamir3pass.Key) error
	ValidateAll() error
}

// RulesEngine
type RulesEngine interface {
	OpponentCanDrawCard(index uint64) bool
}

// request is a function that needs to affect the state of the
// deck. Requests are submitted to the deck's requests channel
// to be processed serially (thereby preventing gross race
// conditions)
type request func()

type Deck struct {
	keys         shamir3pass.Key
	cards        CardHolder
	protocol     *p.Protocol
	done         chan struct{}
	ready        bool
	faceRequests []chan card.CardFace
	requests     chan request
	messages     chan []byte
}

// NewDeck creates a Deck of cards and assumes that the given
// io.ReadWriteCloser is a connection of some sort to another
// Deck.
func NewDeck(DeckConnection io.ReadWriteCloser) (*Deck, error) {
	d := &Deck{}
	done := make(chan struct{})
	p, err := p.NewProtocol(DeckConnection, d, done)
	if err != nil {
		return nil, err
	}
	d.done = done
	d.protocol = p
	d.requests = make(chan request)
	d.messages = make(chan []byte)
	d.ready = false
	go d.handleRequests()

	return d, nil
}

// ready returns whether the deck can be used for network
// operations yet.
func (d *Deck) isReady() bool {
	return d.ready
}

func (d *Deck) handleRequests() {
	for r := range d.requests {
		r()
	}
}

// Start runs the game, and initiates the first hand
func (d *Deck) Start(faces []card.CardFace) error {
	e := make(chan error)
	defer close(e)
	d.requests <- func() {
		d.faceRequests = make([]chan card.CardFace, len(faces))
		prime := shamir3pass.Random1024BitPrime()
		d.keys = shamir3pass.GenerateKeyFromPrime(prime)
		cards, err := card_holder.NewHolder(&d.keys, faces)
		if err != nil {
			e <- err
			return
		}
		d.cards = cards
		enc, _, err := d.cards.GetAllMine()
		if err != nil {
			e <- err
			return
		}
		e <- d.protocol.SendStartDeck(prime, enc)
	}
	err := <-e
	// don't return until we time out or the deck is fully initialized
	deadline := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-deadline.C:
			return fmt.Errorf("Connecting timed out")
		default:
			if d.isReady() {
				return err
			}
		}
	}
}

// Draw draws a single card from the Deck
func (d *Deck) Draw(index uint) (card.CardFace, error) {
	if !d.isReady() {
		return nil, fmt.Errorf("Deck not fully initialized, cannot draw card")
	} else if index >= d.Size() {
		return nil, fmt.Errorf("Deck has size %d, tried to draw card %d", d.Size(), index)
	}
	faces := make(chan card.CardFace)
	defer close(faces)
	d.requests <- func() {
		d.protocol.RequestDecryptCard(uint64(index))
		log.Printf("Requesting decryption of card:\n%v\n", index)
		d.faceRequests[index] = faces
	}
	return <-faces, nil
}

func (d *Deck) Size() uint {
	resp := make(chan uint)
	d.requests <- func() {
		resp <- d.cards.Size()
	}
	return <-resp
}

func (d *Deck) Quit() {
	close(d.done)
	close(d.requests)
}

// Higher-layer communication:
func (d *Deck) Send(message []byte) error {
	e := make(chan error)
	d.requests <- func() {
		e <- d.protocol.SendApplicationMessage(message)
	}
	return <-e
}

func (d *Deck) Recieve() <-chan []byte {
	return d.messages
}

// handler implementations

func (d *Deck) HandleQuit() {
	d.requests <- func() {
		log.Println("QUIT")
		d.Quit()
	}
}
func (d *Deck) HandleStartDeck(deck []*big.Int, prime *big.Int) {
	d.requests <- func() {
		log.Println("START_DECK")
		d.keys = shamir3pass.GenerateKeyFromPrime(prime)
		d.cards, _ = card_holder.HolderFromEncrypted(&d.keys, deck)
		d.faceRequests = make([]chan card.CardFace, len(deck))
		d.ready = true
		both, _, _ := d.cards.GetAllBoth()
		d.protocol.SendEndDeck(both)
	}
}
func (d *Deck) HandleEndDeck(deck []*big.Int) {
	d.requests <- func() {
		log.Println("END_DECK")
		d.cards.SetBothEncrypted(deck)
		d.ready = true
		log.Println(d.cards)
	}
}
func (d *Deck) HandleDecryptCard(index uint64) {
	d.requests <- func() {
		log.Println("DECRYPT_CARD")
		theirs, _ := d.cards.GetTheirs(uint(index))
		log.Println("decrypted: ", theirs)
		d.protocol.SendDecryptedCard(index, theirs)
	}
}
func (d *Deck) HandleDecryptedCard(index uint64, card *big.Int) {
	d.requests <- func() {
		log.Println("ONE_CIPHER_CARD")
		d.cards.SetMine(uint(index), card)
		if d.faceRequests[index] != nil {
			face, _ := d.cards.Get(uint(index))
			d.faceRequests[index] <- face
			d.faceRequests[index] = nil
		}
	}
}

func (d *Deck) HandleAppMessage(data []byte) {
	d.requests <- func() {
		log.Println("APP_MESSAGE")
		d.messages <- data
	}
}
