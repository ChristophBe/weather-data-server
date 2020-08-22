package services

type mailServiceImpl struct {}

func (m mailServiceImpl) SendeEMail(recipientAddress string, recipientName string, msg string) {
	panic("implement me")
}

func (m mailServiceImpl) SendHtmlEMail(recipientAddress string, recipientName string, templateFile string, templateParams interface{}) {
	panic("implement me")
}

