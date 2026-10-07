package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/EliasBlind/EduFlow/internal/sso_service/config"
	"github.com/wneessen/go-mail"
)

//go:embed templates/*.html
var templateFS embed.FS

var (
	templatesDir     = "templates/"
	verificationTmpl = templatesDir + "verification.html"
)

type Mailer struct {
	log      *slog.Logger
	cfg      *config.MailtrapConfig
	client   *mail.Client
	template *template.Template
}

func New(log *slog.Logger, cfg *config.MailtrapConfig) (*Mailer, error) {
	const op = "mailer.New"
	logNew := log.With(
		"op", op,
		"host", cfg.Host,
		"port", cfg.Port,
	)

	tlsPolicy := mail.NoTLS
	if cfg.UseTls {
		tlsPolicy = mail.TLSMandatory
	}

	opts := []mail.Option{
		mail.WithPort(cfg.Port),
		mail.WithTLSPolicy(tlsPolicy),
	}

	if cfg.Username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthPlain),
			mail.WithUsername(cfg.Username),
			mail.WithPassword(cfg.Password),
		)
	}

	c, err := mail.NewClient(
		cfg.Host,
		opts...,
	)

	if err != nil {
		logNew.Error("failed to create mail client", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	tmpl, err := template.ParseFS(templateFS, verificationTmpl)
	if err != nil {
		logNew.Error("failed to pre-parse email template", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logNew.Info("mailer initialized successfully")

	return &Mailer{
		log:      log,
		cfg:      cfg,
		client:   c,
		template: tmpl,
	}, nil
}

func (m *Mailer) SendVerificationCode(
	ctx context.Context,
	email,
	code string,
) error {
	const op = "mailer.SendVerificationCode"
	log := m.log.With(
		"op", op,
		"email", email,
	)

	var body bytes.Buffer
	if err := m.template.Execute(&body, map[string]string{"Code": code}); err != nil {
		log.Error("failed to execute template", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	msg := mail.NewMsg()
	if err := msg.From(m.cfg.SenderEmail); err != nil {
		log.Error("failed to set sender", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	if err := msg.To(email); err != nil {
		log.Error("failed to set recipient", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}
	msg.Subject("EduFlow | Email Confirmation")
	msg.SetBodyString(mail.TypeTextHTML, body.String())

	log.Info("sending email")

	if err := m.client.DialAndSendWithContext(ctx, msg); err != nil {
		log.Error("failed to send email", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("email sent successfully")
	return nil
}
