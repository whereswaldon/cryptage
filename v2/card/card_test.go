package card_test

import (
	"fmt"
	"github.com/sorribas/shamir3pass"
	"math/big"

	. "github.com/whereswaldon/cryptage/v2/card"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getKeyPair() (*shamir3pass.Key, *shamir3pass.Key) {
	prime := shamir3pass.Random1024BitPrime()
	key1 := shamir3pass.GenerateKeyFromPrime(prime)
	key2 := shamir3pass.GenerateKeyFromPrime(prime)
	return &key1, &key2
}

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
				" value and no error", func() {
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
				" value and no error", func() {
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
				" value and no error", func() {
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
	Describe("Setting the mine value on a card", func() {
		Context("Where the mine value is nil", func() {
			It("Should return an error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				card, _ := CardFromTheirs(their, k1)
				err := card.SetMine(nil)
				Expect(err).ToNot(BeNil())
			})
		})
		Context("Where the mine value is valid", func() {
			It("Should return no error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				mine := EncryptString("test", k1)
				card, _ := CardFromTheirs(their, k1)
				err := card.SetMine(mine)
				Expect(err).To(BeNil())
				m2, err := card.Mine()
				Expect(m2).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})
	})
	Describe("Setting the opponent's key on a card", func() {
		Context("Where the key is nil", func() {
			It("Should return an error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				card, _ := CardFromTheirs(their, k1)
				err := card.SetTheirKey(nil)
				Expect(err).ToNot(BeNil())
				Expect(card.HasTheirKey()).ToNot(BeTrue())
			})
		})
		Context("Where the mine value is valid", func() {
			It("Should return no error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				card, _ := CardFromTheirs(their, k1)
				err := card.SetTheirKey(k2)
				Expect(err).To(BeNil())
				Expect(card.HasTheirKey()).To(BeTrue())
			})
		})
	})
	Describe("Validating a card", func() {
		Context("Where the opponent's key is nil", func() {
			It("Should return an error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				card, _ := CardFromTheirs(their, k1)
				err := card.Validate()
				Expect(err).ToNot(BeNil())
			})
		})
		Context("Where the both value is nil", func() {
			It("Should return an error", func() {
				k1, k2 := getKeyPair()
				card, _ := NewCard("test", k1)
				card.SetTheirKey(k2)
				err := card.Validate()
				Expect(err).ToNot(BeNil())
			})
		})
		Context("Where all required values are present and the card is invalid", func() {
			It("Should return an error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				fakeMine := EncryptString("test2", k1)
				card, _ := CardFromTheirs(their, k1)
				card.SetTheirKey(k2)
				card.SetMine(fakeMine)
				err := card.Validate()
				Expect(err).ToNot(BeNil())
			})
		})
		Context("Where all required values are present and the keys are"+
			" unrelated", func() {
			It("Should return an error", func() {
				k1, _ := getKeyPair()
				_, k2 := getKeyPair()
				their := EncryptString("test", k2)
				both := shamir3pass.Encrypt(their, *k1)
				mine := EncryptString("test", k1)
				card, _ := CardFromBoth(both, k1)
				card.SetTheirKey(k2)
				card.SetMine(mine)
				err := card.Validate()
				fmt.Fprintf(GinkgoWriter, "Err: %v", err)
				Expect(err).ToNot(BeNil())
			})
		})
		Context("Where all required values are present and the card is valid", func() {
			It("Should return no error", func() {
				k1, k2 := getKeyPair()
				their := EncryptString("test", k2)
				mine := EncryptString("test", k1)
				both := shamir3pass.Encrypt(their, *k1)
				card, _ := CardFromBoth(both, k1)
				card.SetTheirKey(k2)
				card.SetMine(mine)
				err := card.Validate()
				Expect(err).To(BeNil())
			})
		})
	})
})
