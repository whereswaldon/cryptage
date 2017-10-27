package deck_test

import (
	mconn "github.com/jordwest/mock-conn"
	. "github.com/whereswaldon/cryptage/v2/protocol"
	"io"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockHandler struct {
	messages chan Message
}

func NewMockHandler() *mockHandler {
	return &mockHandler{messages: make(chan Message)}
}
func (m *mockHandler) HandleQuit() {
	m.messages <- Message{Type: QUIT}
}
func (m *mockHandler) HandleStartDeck(deck []*big.Int, prime *big.Int) {
	m.messages <- Message{Type: START_DECK}
}
func (m *mockHandler) HandleEndDeck(deck []*big.Int) {
	m.messages <- Message{Type: END_DECK}
}
func (m *mockHandler) HandleDecryptCard(index uint64) {
	m.messages <- Message{Type: DECRYPT_CARD}
}
func (m *mockHandler) HandleDecryptedCard(index uint64, card *big.Int) {
	m.messages <- Message{Type: ONE_CIPHER_CARD}
}

var _ = Describe("Protocol", func() {
	var (
		p1Conn  io.ReadWriteCloser
		p2Conn  io.ReadWriteCloser
		handler ProtocolHandler
		done    chan struct{}
	)
	BeforeEach(func() {
		connection := mconn.NewConn()
		p1Conn = connection.Client
		p2Conn = connection.Server
		handler = NewMockHandler()
	})
	Describe("Creating a Protocol instance", func() {
		Context("When the provided connection is nil", func() {
			It("Should return an error", func() {
				p, err := NewProtocol(nil, handler, done)
				Expect(err).ToNot(BeNil())
				Expect(p).To(BeNil())
			})
		})
	})
	Describe("Creating a Protocol instance", func() {
		Context("When the provided handler is nil", func() {
			It("Should return an error", func() {
				p, err := NewProtocol(p1Conn, nil, done)
				Expect(err).ToNot(BeNil())
				Expect(p).To(BeNil())
			})
		})
	})
	Describe("Creating a Protocol instance", func() {
		Context("When the provided channel is nil", func() {
			It("Should return an error", func() {
				p, err := NewProtocol(p1Conn, handler, nil)
				Expect(err).ToNot(BeNil())
				Expect(p).To(BeNil())
			})
		})
	})
})
