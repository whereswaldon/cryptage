package cribbage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCribbage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cribbage Suite")
}
