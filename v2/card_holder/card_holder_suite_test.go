package card_holder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCardHolder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CardHolder Suite")
}
