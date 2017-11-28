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
				card.UnmarshalText(cards[1])
				s.Play(p1, card)
				pOut, cardOut = s.Get(s.Size() - 1)
				Expect(pOut).To(BeEquivalentTo(p1))
				Expect(cardOut).To(BeEquivalentTo(card))
			})
		})
	})
	Describe("When you ask for the total of a sequence", func() {
		Context("If no cards have been played", func() {
			It("Should return 0", func() {
				Expect(s.Total()).To(BeEquivalentTo(0))
			})
		})
		Context("If cards have been played", func() {
			It("Should return the total value of those cards", func() {
				card := &Card{}
				for i := 0; i < 6; i++ {
					card.UnmarshalText(cards[i])
					s.Play(1, card)
				}
				Expect(s.Total()).To(BeEquivalentTo(27))
			})
		})
		Context("If you try to play a card that would exceed a total of 31", func() {
			It("CanPlay should return false", func() {
				card := &Card{}
				for i := 0; i < 6; i++ {
					card.UnmarshalText(cards[i])
					s.Play(1, card)
				}
				card.UnmarshalText(cards[10])
				Expect(s.CanPlay(card)).To(BeEquivalentTo(false))
			})
		})
		Context("If you try to play a card that would not exceed a total of 31", func() {
			It("CanPlay should return false", func() {
				card := &Card{}
				for i := 0; i < 5; i++ {
					card.UnmarshalText(cards[i])
					s.Play(1, card)
				}
				card.UnmarshalText(cards[6])
				Expect(s.CanPlay(card)).To(BeEquivalentTo(true))
			})
		})
	})
})
