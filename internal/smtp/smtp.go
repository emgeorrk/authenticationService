package smtplib

import (
	"authenticationService/internal/app"
	"authenticationService/internal/config"
	"authenticationService/internal/logger"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/smtp"
)

func New(cfg config.SMTP) smtp.Auth {
	smtpHost := cfg.Host
	publicKey := cfg.PublicKey
	privateKey := cfg.PrivateKey

	auth := smtp.PlainAuth("", publicKey, privateKey, smtpHost)

	return auth
}

func SendEmail(a app.App, auth smtp.Auth, to, subject, body string) error {
	const op = "smtp.SendEmail"

	log := a.Logger.With(
		slog.String("op", op),
	)

	smtpHost := a.Config.SMTP.Host
	smtpPort := a.Config.SMTP.Port
	senderEmail := a.Config.SMTP.SenderEmail

	conn, err := smtp.Dial(smtpHost + ":" + smtpPort)
	if err != nil {
		return fmt.Errorf("%s: error connecting to SMTP server: %v", op, err)
	}
	defer func(conn *smtp.Client) {
		err := conn.Quit()
		if err != nil {
			log.Error("error closing connection", logger.Err(err))
		}
	}(conn)

	tlsConfig := &tls.Config{
		ServerName:         smtpHost,
		InsecureSkipVerify: true,
	}

	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("%s: error starting TLS: %v", op, err)
	}

	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("%s: error during authentication: %v", op, err)
	}

	if err = conn.Mail(senderEmail); err != nil {
		return fmt.Errorf("%s: error setting sender: %v", op, err)
	}

	if err = conn.Rcpt(to); err != nil {
		return fmt.Errorf("%s: error adding recipient: %v", op, err)
	}

	msg := []byte(subject + "\n" + body)

	w, err := conn.Data()
	if err != nil {
		return fmt.Errorf("%s: error sending data: %v", op, err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("%s: error writing message: %v", op, err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("%s: error closing connection: %v", op, err)
	}

	return nil
}
