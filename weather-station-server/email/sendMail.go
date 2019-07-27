package email

import (
	"../configs"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"
)
/*
func SendMail(recipient string, message string) {
	// Set up authentication information.

	auth := smtp.PlainAuth(
		"",
		configs.MAIL_USERNAME,
		configs.MAIL_PASSWORD,
		configs.MAIL_SERVER,
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		configs.MAIL_SERVER + ":" + configs.MAIL_SERVER_PORT,
		auth,
		configs.MAIL_MAIL_ADDRESS,
		[]string{recipient},
		[]byte(message),
	)
	if err != nil {
		log.Fatal(err)
	}
}
*/


func SendHtmlMail(recipient string,subj string,templateFileName string,themeParams interface{}) error{

	html , err := parsHtml("email/tmpl/" + templateFileName,themeParams)

	if err != nil {
		return err
	}


	return SendMail(recipient,subj,html)


}

func SendMail(recipient string, subj string, msg string)  error{

	from := mail.Address{"Wetter | christophb.de", configs.MAIL_MAIL_ADDRESS}
	to   := mail.Address{"", recipient}


	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Setup message
	message := ""
	for k,v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + msg

	// Connect to the SMTP Server
	servername := configs.MAIL_SERVER + ":" + configs.MAIL_SERVER_PORT

	host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("",configs.MAIL_USERNAME, configs.MAIL_PASSWORD, configs.MAIL_SERVER)

	// TLS config
	tlsconfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
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

func parsHtml( templateFileName string,data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "",err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "",err
	}
	return  buf.String(), nil
}

