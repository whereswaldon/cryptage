package deck_test

import (
	mconn "github.com/jordwest/mock-conn"
	. "github.com/whereswaldon/cryptage/v2/protocol"
	"io"
	"math/big"
	"time"

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

var _ ProtocolHandler = &mockHandler{}

var _ = Describe("Protocol", func() {
	var (
		p1Conn  io.ReadWriteCloser
		p2Conn  io.ReadWriteCloser
		handler *mockHandler
		done    chan struct{}
	)
	BeforeEach(func() {
		connection := mconn.NewConn()
		p1Conn = connection.Client
		p2Conn = connection.Server
		handler = NewMockHandler()
		done = make(chan struct{})
	})
	Describe("Creating a Protocol instance", func() {
		Context("When the provided connection is nil", func() {
			It("Should return an error", func() {
				p, err := NewProtocol(nil, handler, done)
				Expect(err).ToNot(BeNil())
				Expect(p).To(BeNil())
			})
		})
		Context("When the provided handler is nil", func() {
			It("Should return an error", func() {
				p, err := NewProtocol(p1Conn, nil, done)
				Expect(err).ToNot(BeNil())
				Expect(p).To(BeNil())
			})
		})
		Context("When the provided channel is nil", func() {
			It("Should return an error", func() {
				p, err := NewProtocol(p1Conn, handler, nil)
				Expect(err).ToNot(BeNil())
				Expect(p).To(BeNil())
			})
		})
		Context("When the parameters are valid", func() {
			It("Should return a Protocol instance", func() {
				p1, err := NewProtocol(p1Conn, handler, done)
				Expect(err).To(BeNil())
				Expect(p1).ToNot(BeNil())
			})
		})
	})
	Describe("Connecting two Protocol instances", func() {
		awaitMsg := func(wait time.Duration, read chan Message) Message {
			select {
			case <-time.Tick(wait):
				Fail("handler not invoked soon enough, timed out")
				return Message{}
			case msg := <-read:
				return msg
			}
		}
		Context("When one sends the QUIT message", func() {
			It("should make the other one invoke its QuitHandler", func() {
				p1, _ := NewProtocol(p1Conn, handler, done)
				p2, _ := NewProtocol(p2Conn, handler, done)
				var msg Message

				//p1 asks p2 to quit
				Expect(p1.SendQuit()).To(BeNil())
				msg = awaitMsg(time.Second, handler.messages)
				Expect(msg.Type).To(BeEquivalentTo(QUIT))

				//p2 asks p1 to quit
				Expect(p2.SendQuit()).To(BeNil())
				msg = awaitMsg(time.Second, handler.messages)
				Expect(msg.Type).To(BeEquivalentTo(QUIT))
			})
		})
	})
})
