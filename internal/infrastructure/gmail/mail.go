package gmail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type SMTPConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
}

type GmailMail struct {
	smtpConfig SMTPConfig
}

func NewGmailMail(smtpConfig SMTPConfig) *GmailMail {
	return &GmailMail{
		smtpConfig: smtpConfig,
	}
}

func (mail *GmailMail) SendToEmail(fromEmail string, toEmails []string, subject string, htmlBody string) error {
	client, err := createSMTPClient(&mail.smtpConfig)
	if err != nil {
		client.Close()
		return err
	}
	defer client.Close()

	if err = client.Mail(fromEmail); err != nil {
		return err
	}
	for _, email := range toEmails {
		if err = client.Rcpt(email); err != nil {
			return err
		}
	}

	message, err := buildEmailMessage(fromEmail, toEmails, subject, htmlBody)
	if err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(message)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func createSMTPClient(smtpConfig *SMTPConfig) (*smtp.Client, error) {
	auth := smtp.PlainAuth("", smtpConfig.SMTPUsername, smtpConfig.SMTPPassword, smtpConfig.SMTPHost)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Should be false in production
		ServerName:         smtpConfig.SMTPHost,
	}

	// Connect to the SMTP server
	conn, err := tls.Dial("tcp", smtpConfig.SMTPHost+":"+smtpConfig.SMTPPort, tlsConfig)
	if err != nil {
		return nil, err
	}

	client, err := smtp.NewClient(conn, smtpConfig.SMTPHost)
	if err != nil {
		return nil, err
	}

	// Authenticate
	if err = client.Auth(auth); err != nil {
		return nil, err
	}

	return client, nil
}

func buildEmailMessage(fromEmail string, toEmails []string, subject string, htmlBody string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", fromEmail))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", toEmails[0]))

	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: multipart/alternative; boundary=\"boundary123456\"\r\n\r\n")

	// Plaintext
	buf.WriteString("--boundary123456\r\n")
	buf.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	buf.WriteString("This is a plain text fallback.\r\n")

	// HTML
	buf.WriteString("--boundary123456\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")
	buf.WriteString("<!DOCTYPE html>\r\n")
	buf.WriteString("<html>\r\n")
	buf.WriteString("<head>\r\n")
	buf.WriteString("<meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\">\r\n")
	buf.WriteString("</head>\r\n")
	buf.WriteString(fmt.Sprintf(`<body>%s</body>`, htmlBody))
	buf.WriteString("</html>\r\n")

	buf.WriteString("\r\n--boundary123456--")

	return buf.Bytes(), nil
}
