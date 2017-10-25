package card_test

import (
	"github.com/sorribas/shamir3pass"
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
			It("Should return a card and no error", func() {
				key := shamir3pass.GenerateKey(1024)
				card, err := NewCard("test", &key)
				Expect(err).To(BeNil())
				Expect(card).ToNot(BeNil())
			})
		})
	})
})
