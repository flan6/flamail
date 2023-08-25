package flamail

import "github.com/flan6/flamail/entity"

type Mailer interface {
	Send(email entity.Email, attachments ...entity.Attachment) error
}
