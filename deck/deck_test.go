package deck_test

import (
	"bytes"
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
})
