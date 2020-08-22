package services

type MailService interface {
	SendeMail(recipientAddress string, recipientName string, msg string)
	SendHtmlMail(recipientAddress string, recipientName string, templateFile string, templateParams interface{})
}

func GetMailService()  MailService{
	return mailServiceImpl{}
}