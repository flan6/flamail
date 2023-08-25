package entity

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
