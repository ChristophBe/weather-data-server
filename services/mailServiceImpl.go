package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/ChristophBe/weather-data-server/config"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
)

type mailServiceImpl struct{}

func (m mailServiceImpl) SendeTxtMail(to mail.Address, subject, messageContent string) error {
	headers, err := m.createHeader(to, subject)
	if err != nil {
		return err
	}
	return m.sendMail(to, headers, messageContent)
}

func (m mailServiceImpl) SendHtmlMail(to mail.Address, subject, templateFile string, templateParams interface{}) error {
	html, err := m.handleMailTemplate(templateFile, templateParams)
	if err != nil {
		return err
	}

	headers, err := m.createHeader(to, subject)
	if err != nil {
		return err
	}
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	return m.sendMail(to, headers, html)
}

func (m mailServiceImpl) createHeader(to mail.Address, subject string) (headers map[string]string, err error) {

	conf, err := config.GetConfigManager().GetConfig()

	if err != nil {
		return
	}

	from := mail.Address{
		Name:    conf.Mail.MailName,
		Address: conf.Mail.MailAddress,
	}
	// Setup headers
	headers = make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	return
}
func (m mailServiceImpl) sendMail(to mail.Address, headers map[string]string, mailContent string) error {

	conf, err := config.GetConfigManager().GetConfig()

	if err != nil {
		panic(err)
	}
	mailConfiguration := conf.Mail

	data := m.buildMailData(headers, mailContent)

	// Connect to the SMTP Server
	servername := fmt.Sprintf("%s:%d", mailConfiguration.Server, mailConfiguration.ServerPort)

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", mailConfiguration.Username, mailConfiguration.Password, mailConfiguration.Server)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	smtpClient, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = smtpClient.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = smtpClient.Mail(mailConfiguration.MailAddress); err != nil {
		return err
	}

	if err = smtpClient.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	messageWriter, err := smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = messageWriter.Write([]byte(data))
	if err != nil {
		return err
	}

	err = messageWriter.Close()
	if err != nil {
		return err
	}
	err = smtpClient.Quit()
	return err
}

func (m mailServiceImpl) buildMailData(headers map[string]string, mailContent string) string {
	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + mailContent
	return message
}

func (m mailServiceImpl) handleMailTemplate(templateFileName string, templateParameters interface{}) (htmlString string, err error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, templateParameters)
	if err != nil {
		return
	}
	htmlString = buf.String()
	return
}
