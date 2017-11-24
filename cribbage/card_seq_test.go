package cribbage_test

import (
	. "github.com/whereswaldon/cryptage/cribbage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CardSeq", func() {
	var s *Sequence
	cards := Cards()
	BeforeEach(func() {
		s = NewSeq()
	})

	Describe("When you play a card", func() {
		Context("You should be able to get it back", func() {
			It("by Get-ing the card at Size()-1", func() {
				p1 := 1
				card := &Card{}
				card.UnmarshalText(cards[0])
				s.Play(p1, card)
				pOut, cardOut := s.Get(s.Size() - 1)
				Expect(pOut).To(BeEquivalentTo(p1))
				Expect(cardOut).To(BeEquivalentTo(card))
			})
		})
	})
})
