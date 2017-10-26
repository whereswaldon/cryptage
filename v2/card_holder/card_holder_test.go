package card_holder_test

import (
	"github.com/sorribas/shamir3pass"
	. "github.com/whereswaldon/cryptage/v2/card_holder"
	. "github.com/whereswaldon/cryptage/v2/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var Faces []CardFace = []CardFace{"ACE", "KING", "QUEEN"}

var _ = Describe("CardHolder", func() {
	Describe("Creating a CardHolder from scratch", func() {
		Context("When the key is nil", func() {
			It("Should return an error", func() {
				holder, err := NewHolder(nil, Faces)
				Expect(err).ToNot(BeNil())
				Expect(holder).To(BeNil())
			})
		})
		Context("When the faces are nil", func() {
			It("Should return an error", func() {
				key := shamir3pass.GenerateKey(1024)
				holder, err := NewHolder(&key, nil)
				Expect(err).ToNot(BeNil())
				Expect(holder).To(BeNil())
			})
		})
		Context("When the faces are empty", func() {
			It("Should return an error", func() {
				key := shamir3pass.GenerateKey(1024)
				faces := make([]CardFace, 0)
				holder, err := NewHolder(&key, faces)
				Expect(err).ToNot(BeNil())
				Expect(holder).To(BeNil())
			})
		})
		Context("When the arguments are valid", func() {
			It("Should return a CardHolder an a nil error", func() {
				key := shamir3pass.GenerateKey(1024)
				holder, err := NewHolder(&key, Faces)
				Expect(err).To(BeNil())
				Expect(holder).ToNot(BeNil())
			})
		})
	})
})
