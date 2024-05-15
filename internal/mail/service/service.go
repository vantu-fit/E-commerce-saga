package service

import (
	"github.com/vantu-fit/saga-pattern/internal/mail/sender"
	"github.com/vantu-fit/saga-pattern/internal/mail/service/command"
)

type Command struct {
	SendRegisterEmailHandler command.SendRegisterEmailHandler
}

type Query struct {
	sender sender.MailSender
}

type Service struct {
	Command
	Query
}

func NewService(
	sender sender.MailSender,
) Service {
	return Service{
		Command: Command{
			SendRegisterEmailHandler: command.NewEmailSenderHandler(sender),
		},
		Query: Query{},
	}
}
