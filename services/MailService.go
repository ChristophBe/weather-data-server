package services

import "net/mail"

type MailService interface {
	SendeTxtMail(to mail.Address, subject, msg string) error
	SendHtmlMail(to mail.Address, subject, templateFile string, templateParams interface{}) error
}

func GetMailService() MailService {
	return mailServiceImpl{}
}
