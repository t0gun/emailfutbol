package mailer

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Mail struct {
	Server    string
	Port      int
	Email     string
	Password  string
	Recipient string
	Subject   string
	Body      string

	client *smtp.Client
}

// NewMailer returns a fully initialized Mail object
func NewMailer(server string, port int, email, password string) *Mail {
	return &Mail{
		Server:   server,
		Port:     port,
		Email:    email,
		Password: password,
	}
}

func (m *Mail) SendEmail() error {
	if err := m.dialTLS(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer m.client.Close()

	if err := m.login(); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	if err := m.client.Mail(m.Email, nil); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	if err := m.client.Rcpt(m.Recipient, nil); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}
	wc, err := m.client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}
	defer wc.Close()
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s\r\n", m.Email, m.Recipient, m.Subject, m.Body)

	if _, err = fmt.Fprintf(wc, msg); err != nil {
		return fmt.Errorf("writing message failed: %s", err)
	}

	return nil
}


func (m *Mail) dialTLS() error {
	tlsConfig := &tls.Config{
		ServerName: m.Server,
	}
	address := fmt.Sprintf("%s:%d", m.Server, m.Port)
	client, err := smtp.DialStartTLS(address, tlsConfig)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *Mail) login() error {
	auth := sasl.NewPlainClient("", m.Email, m.Password)
	if err := m.client.Auth(auth); err != nil {
		return err
	}
	return nil
}
