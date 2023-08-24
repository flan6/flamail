package flamail

type SmtpMailer struct {
	SmtpServerAuthEmail string // Optional
	SmtpServerAddress   string
	SmtpServerDomain    string
}

func NewSmtpMailer(options func(*SmtpMailer)) SmtpMailer {
	smtpMailer := SmtpMailer{}
	options(&smtpMailer)
	return smtpMailer
}

func (s SmtpMailer) Send(email Email, config string, attachments ...Attachment) error {
	return nil
}
