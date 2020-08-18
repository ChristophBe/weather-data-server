package email

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

func SendHtmlMail(recipient string, subj string, templateFileName string, themeParams interface{}) error {

	html, err := parsHtml("email/tmpl/"+templateFileName, themeParams)

	if err != nil {
		return err
	}

	return SendMail(recipient, subj, html)

}

func SendMail(recipient string, subj string, msg string) error {

	conf, err := config.GetConfigManager().GetConfig()

	if err != nil {
		return err
	}

	mailConf := conf.Mail
	from := mail.Address{"Wetter | christophb.de", mailConf.MailAddress}
	to := mail.Address{"", recipient}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + msg

	// Connect to the SMTP Server
	servername := fmt.Sprintf("%s:%d", mailConf.Server, mailConf.ServerPort)

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", mailConf.Username, mailConf.Password, mailConf.Server)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil

}

func parsHtml(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
