package smtp

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"dillmann.com.br/nginx-ignition/core/notification"
)

type mailSettings struct {
	host        string
	username    string
	password    string
	fromAddress string
	port        int
	useTLS      bool
	useStartTLS bool
}

type mailSender interface {
	Send(
		ctx context.Context,
		settings mailSettings,
		recipients []string,
		message []byte,
	) error
}

type defaultMailSender struct{}

func (defaultMailSender) Send(
	_ context.Context,
	settings mailSettings,
	recipients []string,
	message []byte,
) error {
	address := net.JoinHostPort(settings.host, strconv.Itoa(settings.port))

	if settings.useTLS {
		return sendWithTLS(address, settings, recipients, message)
	}

	return sendWithPlainOrStartTLS(address, settings, recipients, message)
}

func sendWithTLS(
	address string,
	settings mailSettings,
	recipients []string,
	message []byte,
) error {
	tlsConfig := &tls.Config{ServerName: settings.host}
	connection, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		return err
	}
	defer connection.Close()

	client, err := smtp.NewClient(connection, settings.host)
	if err != nil {
		return err
	}
	defer client.Close()

	return deliver(client, settings, recipients, message)
}

func sendWithPlainOrStartTLS(
	address string,
	settings mailSettings,
	recipients []string,
	message []byte,
) error {
	client, err := smtp.Dial(address)
	if err != nil {
		return err
	}
	defer client.Close()

	if settings.useStartTLS {
		tlsConfig := &tls.Config{ServerName: settings.host}
		if err := client.StartTLS(tlsConfig); err != nil {
			return err
		}
	}

	return deliver(client, settings, recipients, message)
}

func deliver(
	client *smtp.Client,
	settings mailSettings,
	recipients []string,
	message []byte,
) error {
	if settings.username != "" || settings.password != "" {
		auth := smtp.PlainAuth("", settings.username, settings.password, settings.host)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(settings.fromAddress); err != nil {
		return err
	}

	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	if _, err := writer.Write(message); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return client.Quit()
}

func parseMailSettings(parameters map[string]any) (mailSettings, []string) {
	host, _ := parameters[hostFieldID].(string)
	portText, _ := parameters[portFieldID].(string)
	port, _ := strconv.Atoi(strings.TrimSpace(portText))
	fromAddress, _ := parameters[fromFieldID].(string)
	recipientsText, _ := parameters[toFieldID].(string)
	useTLS, _ := parameters[useTLSFieldID].(bool)
	useStartTLS, _ := parameters[useStartTLSFieldID].(bool)

	settings := mailSettings{
		host:        strings.TrimSpace(host),
		port:        port,
		fromAddress: strings.TrimSpace(fromAddress),
		useTLS:      useTLS,
		useStartTLS: useStartTLS,
	}

	if username, casted := parameters[usernameFieldID].(string); casted {
		settings.username = strings.TrimSpace(username)
	}

	if password, casted := parameters[passwordFieldID].(string); casted {
		settings.password = password
	}

	return settings, parseRecipients(recipientsText)
}

func parseRecipients(recipientsText string) []string {
	parts := strings.Split(recipientsText, ",")
	recipients := make([]string, 0, len(parts))

	for _, part := range parts {
		address := strings.TrimSpace(part)
		if address == "" {
			continue
		}

		recipients = append(recipients, address)
	}

	return recipients
}

func buildMessage(deliverable notification.Deliverable, fromAddress string) []byte {
	plainBody := formatPlainBody(deliverable)
	htmlBody := formatHTMLBody(deliverable)

	var buffer bytes.Buffer
	_, _ = fmt.Fprintf(&buffer, "From: %s\r\n", fromAddress)
	encodedTitle := mime.QEncoding.Encode("utf-8", deliverable.Title)
	_, _ = fmt.Fprintf(&buffer, "Subject: %s\r\n", encodedTitle)
	_, _ = buffer.WriteString("MIME-Version: 1.0\r\n")

	boundary := fmt.Sprintf("nginx-ignition-%d", time.Now().UnixNano())
	_, _ = fmt.Fprintf(&buffer, "Content-Type: multipart/alternative; boundary=%q\r\n", boundary)
	_, _ = buffer.WriteString("\r\n")

	writePart(&buffer, boundary, "text/plain; charset=utf-8", plainBody)
	writePart(&buffer, boundary, "text/html; charset=utf-8", htmlBody)

	_, _ = fmt.Fprintf(&buffer, "--%s--\r\n", boundary)

	return buffer.Bytes()
}

func writePart(buffer *bytes.Buffer, boundary, contentType, body string) {
	_, _ = fmt.Fprintf(buffer, "--%s\r\n", boundary)
	_, _ = fmt.Fprintf(buffer, "Content-Type: %s\r\n", contentType)
	_, _ = buffer.WriteString("Content-Transfer-Encoding: quoted-printable\r\n\r\n")

	writer := quotedprintable.NewWriter(buffer)
	_, _ = writer.Write([]byte(body))
	_ = writer.Close()
	_, _ = buffer.WriteString("\r\n")
}

func formatPlainBody(deliverable notification.Deliverable) string {
	var builder strings.Builder
	_, _ = builder.WriteString(deliverable.Summary)
	_, _ = builder.WriteString("\n\n")

	for _, section := range deliverable.Sections {
		if section.Title != nil {
			_, _ = builder.WriteString(*section.Title)
			_, _ = builder.WriteString("\n")
		}
		_, _ = builder.WriteString(section.Body)
		_, _ = builder.WriteString("\n\n")
	}

	for _, action := range deliverable.Actions {
		_, _ = builder.WriteString(action.Label)
		_, _ = builder.WriteString(": ")
		_, _ = builder.WriteString(action.URL)
		_, _ = builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String())
}
