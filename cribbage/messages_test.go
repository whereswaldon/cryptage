package cribbage_test

import (
	. "github.com/whereswaldon/cryptage/cribbage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Messages", func() {
	Context("Encoding and decoding a message", func() {
		It("Should come out the same as it went in", func() {
			m := &Message{Type: TO_CRIB, Val: 77}
			enc, err := Encode(m)
			Expect(err).To(BeNil())
			Expect(enc).ToNot(BeNil())
			dec, err := Decode(enc)
			Expect(err).To(BeNil())
			Expect(dec).To(BeEquivalentTo(m))
		})
	})
})
