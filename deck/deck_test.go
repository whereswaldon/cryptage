package deck_test

import (
	"bytes"
	mconn "github.com/jordwest/mock-conn"
	. "github.com/whereswaldon/cryptage/deck"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type ClosableBuffer struct {
	*bytes.Buffer
}

func (c *ClosableBuffer) Close() error {
	return nil
}

var _ = Describe("Deck", func() {
	Describe("Creating a deck", func() {
		Context("With a nil io.ReadWriteCloser", func() {
			It("Should return an error and a nil Deck", func() {
				deck, err := NewDeck(nil)
				Expect(err).ToNot(BeNil())
				Expect(deck).To(BeNil())
			})
		})
		Context("With a valid io.ReadWriteCloser", func() {
			It("Should return a Deck an no error", func() {
				buf := &ClosableBuffer{bytes.NewBuffer([]byte("testing"))}
				deck, err := NewDeck(buf)
				Expect(err).To(BeNil())
				Expect(deck).ToNot(BeNil())
			})
		})
	})
	Describe("Connecting two decks", func() {
		var (
			d1, d2 *Deck
		)
		BeforeEach(func() {
			c := mconn.NewConn()
			d1, _ = NewDeck(c.Client)
			d2, _ = NewDeck(c.Server)
		})
		Context("When one calls Draw before one has called Start", func() {
			It("All calls to draw should result in errors", func() {
				face, err := d1.Draw()
				Expect(err).ToNot(BeNil())
				Expect(face).To(BeEquivalentTo(""))
				face2, err := d2.Draw()
				Expect(err).ToNot(BeNil())
				Expect(face2).To(BeEquivalentTo(""))
				Expect(face).To(Equal(face2))
			})
		})
		Context("When one calls Draw after one has called Start", func() {
			It("All calls to draw should return card faces", func() {
				Expect(d1.Start()).To(BeNil())
				face, err := d1.Draw()
				Expect(err).To(BeNil())
				Expect(face).ToNot(BeEquivalentTo(""))
				face2, err := d2.Draw()
				Expect(err).To(BeNil())
				Expect(face2).ToNot(BeEquivalentTo(""))
				Expect(face).To(Equal(face2))
			})
		})
	})
})
