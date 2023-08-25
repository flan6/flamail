package smtpsender

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/flan6/flamail"
)

type SmtpMailer struct {
	SmtpServerAuthEmail    string // Optional
	SmtpServerAuthPassword string // Optional

	SmtpServerAddress string
	SmtpServerDomain  string
}

func NewSmtpMailer(options func(*SmtpMailer)) SmtpMailer {
	smtpMailer := SmtpMailer{}

	if options != nil {
		options(&smtpMailer)
	}
	return smtpMailer
}

func WithGmail(login, pass string) func(*SmtpMailer) {
	return func(mailer *SmtpMailer) {
		mailer.SmtpServerAuthEmail = login
		mailer.SmtpServerAuthPassword = pass
		mailer.SmtpServerAddress = "smtp.gmail.com:587"
		mailer.SmtpServerDomain = "smtp.gmail.com"
	}
}

func (s SmtpMailer) Send(content flamail.Email, attachments ...flamail.Attachment) error {
	fromAddress, err := mail.ParseAddress(content.From)
	if err != nil {
		return err
	}

	toAddress, err := mail.ParseAddress(content.To)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	email := make(textproto.MIMEHeader)
	email.Set("From", fromAddress.String())
	email.Set("To", toAddress.String())
	email.Set("Return-Path", fromAddress.String())
	email.Set("List-Unsubscribe", content.UnsubscribeLink)
	email.Set("Subject", mime.QEncoding.Encode("utf-8", content.Subject))
	email.Set("MIME-Version", "1.0")
	email.Set("Content-Type", "multipart/mixed; boundary=\""+writer.Boundary()+"\"")

	_, err = writer.CreatePart(email)
	if err != nil {
		return err
	}

	// html body
	email = make(textproto.MIMEHeader)
	email.Set("Content-Type", "text/html; charset=\"utf-8\"")
	email.Set("Content-Transfer-Encoding", "8bit")
	part, err := writer.CreatePart(email)
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(content.Body))
	if err != nil {
		return err
	}

	for _, attachment := range attachments {
		header := textproto.MIMEHeader{
			"Content-Type": {attachment.ContentType},
		}

		if attachment.AddCustomHeaders != nil {
			attachment.AddCustomHeaders(&header)
		}

		part, err = writer.CreatePart(header)
		if err != nil {
			return err
		}

		_, err = part.Write(attachment.Content)
		if err != nil {
			return err
		}
	}

	// Remove boundary antes do header
	message := buf.String()
	if strings.Count(message, "\n") < 2 {
		return fmt.Errorf("invalid e-mail content")
	}
	message = strings.SplitN(message, "\n", 2)[1]

	emailClient := &smtp.Client{}

	conn, err := smtp.Dial(s.SmtpServerAddress)
	if err != nil {
		return fmt.Errorf("fail to send Hello message; %w", err)
	}

	err = conn.StartTLS(
		&tls.Config{
			ServerName:         s.SmtpServerDomain,
			InsecureSkipVerify: false,
		},
	)
	if err != nil {
		return fmt.Errorf("could not start TLS connection; %w", err)
	}

	err = conn.Auth(smtp.PlainAuth("", s.SmtpServerAuthEmail, s.SmtpServerAuthPassword, s.SmtpServerDomain))
	if err != nil {
		return fmt.Errorf("could not auth to SMTP; %w", err)
	}

	defer emailClient.Close()
	if err = emailClient.Mail(fromAddress.Address); err != nil {
		return fmt.Errorf("fail to set from address; %w", err)
	}

	if err = emailClient.Rcpt(toAddress.Address); err != nil {
		return fmt.Errorf("fail to send Rcpt command; %w", err)
	}

	// Send the email body.
	smtpWriter, err := emailClient.Data()
	if err != nil {
		return fmt.Errorf("fail to send Data command; %w", err)
	}

	if _, err = smtpWriter.Write([]byte(message)); err != nil {
		return fmt.Errorf("could not write the message; %w", err)
	}

	if err = smtpWriter.Close(); err != nil {
		return fmt.Errorf("fail to close message; %w", err)
	}

	// Send the QUIT command and close the connection.
	if err := emailClient.Quit(); err != nil {
		return fmt.Errorf("fail to quit and close connection; %w", err)
	}

	return nil
}
