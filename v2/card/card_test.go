package card_test

import (
	"github.com/sorribas/shamir3pass"
	"math/big"

	. "github.com/whereswaldon/cryptage/v2/card"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Card", func() {
	Describe("Creating a Card from a string", func() {
		Context("Where the string is empty", func() {
			It("Should return no card and an error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := NewCard("", &key)
				Expect(err).ToNot(BeNil())
				Expect(card).To(BeNil())
			})
		})
		Context("Where the key is empty", func() {
			It("Should return no card and an error", func() {
				card, err := NewCard("test", nil)
				Expect(err).ToNot(BeNil())
				Expect(card).To(BeNil())
			})
		})
		Context("Where both arguments are valid", func() {
			It("Should return a card with a valid Face and Mine"+
				"value and no error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := NewCard("test", &key)
				Expect(err).To(BeNil())
				Expect(card).ToNot(BeNil())
				face, err := card.Face()
				Expect(err).To(BeNil())
				Expect(face).ToNot(BeEmpty())
				mine, err := card.Mine()
				Expect(err).To(BeNil())
				Expect(mine).ToNot(BeNil())
			})
		})
	})
	Describe("Creating a Card from the opponent's encrypted face", func() {
		Context("Where the encrypted face is nil", func() {
			It("Should return no card and an error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := CardFromTheirs(nil, &key)
				Expect(err).ToNot(BeNil())
				Expect(card).To(BeNil())
			})
		})
		Context("Where the key is empty", func() {
			It("Should return no card and an error", func() {
				card, err := CardFromTheirs(big.NewInt(0), nil)
				Expect(err).ToNot(BeNil())
				Expect(card).To(BeNil())
			})
		})
		Context("Where both arguments are valid", func() {
			It("Should return a card with a valid Theirs and Both"+
				"value and no error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := CardFromTheirs(big.NewInt(0), &key)
				Expect(err).To(BeNil())
				Expect(card).ToNot(BeNil())
				theirs, err := card.Theirs()
				Expect(err).To(BeNil())
				Expect(theirs).ToNot(BeNil())
				both, err := card.Both()
				Expect(err).To(BeNil())
				Expect(both).ToNot(BeNil())
			})
		})
	})
	Describe("Creating a Card from both players' encrypted face", func() {
		Context("Where the encrypted face is nil", func() {
			It("Should return no card and an error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := CardFromBoth(nil, &key)
				Expect(err).ToNot(BeNil())
				Expect(card).To(BeNil())
			})
		})
		Context("Where the key is empty", func() {
			It("Should return no card and an error", func() {
				card, err := CardFromBoth(big.NewInt(0), nil)
				Expect(err).ToNot(BeNil())
				Expect(card).To(BeNil())
			})
		})
		Context("Where both arguments are valid", func() {
			It("Should return a card with a valid Theirs and Both"+
				"value and no error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := CardFromBoth(big.NewInt(0), &key)
				Expect(err).To(BeNil())
				Expect(card).ToNot(BeNil())
				theirs, err := card.Theirs()
				Expect(err).To(BeNil())
				Expect(theirs).ToNot(BeNil())
				both, err := card.Both()
				Expect(err).To(BeNil())
				Expect(both).ToNot(BeNil())
			})
		})
	})
})
