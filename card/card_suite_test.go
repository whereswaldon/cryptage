package card_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Card Suite")
}
