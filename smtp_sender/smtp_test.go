package smtpsender_test

import (
	"testing"

	"github.com/flan6/flamail"
	smtpsender "github.com/flan6/flamail/smtp_sender"
)

func Test_SMTPSend(t *testing.T) {
	t.Run("fail - to connect, invalid password", func(t *testing.T) {
		mailer := smtpsender.NewSmtpMailer(
			smtpsender.WithGmail("example@gmail.com", "invalid"),
		)

		if err := mailer.Send(flamail.Email{
			From: "example@example.com",
			To:   "example@example.com",
		}); err == nil {
			t.Fail()
		}
	})
}
