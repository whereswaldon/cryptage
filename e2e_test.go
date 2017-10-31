package main_test

import (
	"github.com/sorribas/shamir3pass"
	"github.com/whereswaldon/cryptage/card"
	"github.com/whereswaldon/cryptage/card_holder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var Faces []card.CardFace = []card.CardFace{
	card.CardFace("ACE"),
	card.CardFace("KING"),
	card.CardFace("QUEEN"),
	card.CardFace("JACK"),
}

var _ = Describe("E2e", func() {
	Describe("When two cardholders are created", func() {
		Context("With related keys", func() {
			It("Should be able to exchange cards"+
				"and fully validate both decks", func() {
				prime := shamir3pass.Random1024BitPrime()
				key1 := shamir3pass.GenerateKeyFromPrime(prime)
				key2 := shamir3pass.GenerateKeyFromPrime(prime)

				// create player1's card holder
				p1holder, err := card_holder.NewHolder(&key1, Faces)
				Expect(err).To(BeNil())
				p1encrypted, ok, err := p1holder.GetAllMine()
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())

				// create player2's card holder
				p2holder, err := card_holder.HolderFromEncrypted(&key2, p1encrypted)
				Expect(err).To(BeNil())
				bothEncrypted, ok, err := p2holder.GetAllBoth()
				Expect(err).To(BeNil())
				Expect(ok).To(BeTrue())

				// sync player1 to the new encrypted values
				err = p1holder.SetBothEncrypted(bothEncrypted)
				Expect(err).To(BeNil())

				// draw each card from player1's deck and give them
				// to player2
				for i := range Faces {
					encrypted, err := p1holder.GetTheirs(uint(i))
					Expect(err).To(BeNil())
					err = p2holder.SetMine(uint(i), encrypted)
					Expect(err).To(BeNil())
					face, err := p2holder.Get(uint(i))
					Expect(err).To(BeNil())
					Expect(Faces).To(ContainElement(face))

				}

				// validate all of player2's card integrity
				err = p2holder.SetTheirKey(&key1)
				Expect(err).To(BeNil())
				err = p2holder.ValidateAll()
				Expect(err).To(BeNil())
			})
		})
	})
})
