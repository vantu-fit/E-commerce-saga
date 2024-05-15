package command

import (
	"context"

	"github.com/vantu-fit/saga-pattern/internal/mail/sender"
)

type Email struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

type SendRegisterEmailHandler CommandHandler[Email]

type sendRegisterEmailHandler struct {
	sender sender.MailSender
}

func NewEmailSenderHandler(sender sender.MailSender) SendRegisterEmailHandler {
	return &sendRegisterEmailHandler{
		sender: sender,
	}
}

func (h *sendRegisterEmailHandler) Handle(ctx context.Context, cmd Email) error {
	return h.sender.SendRegisterEmail(sender.Message{
		From: 	  cmd.From,
		FromName: cmd.FromName,
		To:       cmd.To,
		Subject:  cmd.Subject,
		Attachments: cmd.Attachments,
		Data:     cmd.Data,
		DataMap:  cmd.DataMap,
	})
}
