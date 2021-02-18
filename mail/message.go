package mail

import (
	"log"

	"github.com/emersion/go-imap"
)

// Message facade for IMAP message
type Message struct {
	// The message unique identifier. It must be greater than or equal to 1.
	UID uint32

	// Subject of the message
	Subject string

	// Body of the message
	Body string
}

// NewMessage instantiate msg
func NewMessage(imapMsg *imap.Message) *Message {
	msg := &Message{}
	msg.UID = imapMsg.Uid
	msg.Subject = imapMsg.Envelope.Subject

	for _, value := range imapMsg.Body {
		len := value.Len()
		buf := make([]byte, len)
		n, err := value.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		if n != len {
			log.Fatal("Didn't read correct length")
		}
		msg.Body = string(buf)
	}
	return msg
}
