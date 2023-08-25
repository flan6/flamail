package flamail

import "net/textproto"

type Email struct {
	From            string
	To              string
	Subject         string
	Body            string
	UnsubscribeLink string
}

type Attachment struct {
	ContentType      string
	Content          []byte
	AddCustomHeaders func(*textproto.MIMEHeader)
}

type Mailer interface {
	Send(email Email, attachments ...Attachment) error
}
