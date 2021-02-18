package mail

import (
	"log"
	"time"

	"github.com/emersion/go-imap"
	idle "github.com/emersion/go-imap-idle"
	imapClient "github.com/emersion/go-imap/client"
)

//VSBTODO: client lock

// Client is facade for IMAP client
type Client struct {
	Imap *imapClient.Client
}

//NewClient instantiates client
func NewClient(imapHost string, user string, password string) *Client {
	c := &Client{}
	c.connect(imapHost)
	c.login(user, password)
	c.selectInbox()
	return c
}

//Logout client
func (client *Client) Logout() {
	client.Imap.Logout()
}

//MarkMsgSeen marking message as seen on the server
func (client *Client) MarkMsgSeen(msg *Message) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(msg.UID)
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err := client.Imap.UidStore(seqSet, item, flags, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// WaitNewMsgs ...
func (client *Client) WaitNewMsgs(msgs chan<- *Message, pollInterval time.Duration) {
	idleClient := idle.NewClient(client.Imap)

	updates := make(chan imapClient.Update)
	client.Imap.Updates = updates

	done := make(chan error, 1)
	go func() {
		done <- idleClient.IdleWithFallback(nil, pollInterval)
	}()

	for {
		select {
		case update := <-updates:
			_, ok := update.(*imapClient.MailboxUpdate)
			if ok {
				log.Println("Got Mailbox update")
				for _, msg := range client.fetchUnseenMsgs() {
					msgs <- msg
				}
			}
		case err := <-done:
			if err != nil {
				log.Fatal(err)
			}
			log.Println("No idling anymore")
			return
		}
	}

}

func (client *Client) fetchUnseenMsgs() []*Message {
	uids := client.searchUnSseenMsgs()
	if len(uids) > 0 {
		return client.fetchMsgs(uids)
	}
	return []*Message{}
}

func (client *Client) fetchMsgs(seqNums []uint32) []*Message {

	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNums...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- client.Imap.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822Text}, messages)
	}()

	out := make([]*Message, 0)
	for msg := range messages {
		out = append(out, NewMessage(msg))
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
	return out
}

func (client *Client) searchUnSseenMsgs() []uint32 {
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	uids, err := client.Imap.Search(criteria)
	if err != nil {
		log.Println(err)
	}
	log.Println("Found unseen msgs:", uids)
	return uids
}

func (client *Client) connect(imapHost string) {
	c, err := imapClient.DialTLS(imapHost, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected")
	client.Imap = c
}

func (client *Client) login(user, password string) {
	if err := client.Imap.Login(user, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")
}

func (client *Client) selectInbox() {
	mbox, err := client.Imap.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("INBOX selected, flags:", mbox.Flags)
}
