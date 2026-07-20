package email

import (
	"fmt"
	"log"
	"net/smtp"
	"multicliente-backend/internal/platform/config"
)

type EmailService interface {
	SendPasswordResetEmail(toEmail string, token string) error
}

type emailService struct {
	cfg *config.Config
}

func NewEmailService(cfg *config.Config) EmailService {
	return &emailService{cfg: cfg}
}

func (s *emailService) SendPasswordResetEmail(toEmail string, token string) error {
	// Always print to console/log in local environment for easy debugging
	log.Printf("\n========================================================")
	log.Printf("  [DEV LOCAL EMAIL RECOVERY LOG]")
	log.Printf("  Para: %s", toEmail)
	log.Printf("  Código de recuperación: %s", token)
	log.Printf("========================================================\n")

	// If SMTP host is configured (and not empty), attempt to send real email
	if s.cfg.SMTPHost != "" {
		subject := "Subject: Restablecimiento de Contraseña - Base Multicliente\r\n"
		mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
		from := fmt.Sprintf("From: %s\r\n", s.cfg.SMTPFrom)
		to := fmt.Sprintf("To: %s\r\n", toEmail)

		body := fmt.Sprintf(`
			<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<style>
					body { font-family: Arial, sans-serif; background-color: #f4f6f9; margin: 0; padding: 20px; }
					.container { max-width: 500px; margin: 0 auto; background: #ffffff; padding: 30px; border-radius: 12px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
					.header { text-align: center; margin-bottom: 20px; }
					.code-box { background: #eef2ff; border: 1px dashed #6366f1; color: #4f46e5; font-size: 28px; font-weight: bold; text-align: center; padding: 15px; border-radius: 8px; letter-spacing: 5px; margin: 20px 0; }
					.footer { margin-top: 25px; font-size: 12px; color: #888888; text-align: center; }
				</style>
			</head>
			<body>
				<div class="container">
					<div class="header">
						<h2>Recuperación de Contraseña</h2>
					</div>
					<p>Hola,</p>
					<p>Hemos recibido una solicitud para restablecer la contraseña de tu cuenta en <strong>Base Multicliente</strong>.</p>
					<p>Tu código de recuperación es:</p>
					<div class="code-box">%s</div>
					<p>Ingresa este código en la aplicación para configurar tu nueva contraseña. Este código expirará en 15 minutos.</p>
					<p>Si no solicitaste este cambio, puedes ignorar este correo de manera segura.</p>
					<div class="footer">
						&copy; Base Multicliente - Todos los derechos reservados.
					</div>
				</div>
			</body>
			</html>
		`, token)

		msg := []byte(from + to + subject + mime + "\r\n" + body)

		addr := fmt.Sprintf("%s:%s", s.cfg.SMTPHost, s.cfg.SMTPPort)

		var auth smtp.Auth
		if s.cfg.SMTPUsername != "" && s.cfg.SMTPPassword != "" {
			auth = smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
		}

		err := smtp.SendMail(addr, auth, s.cfg.SMTPFrom, []string{toEmail}, msg)
		if err != nil {
			log.Printf("[EmailService Error] No se pudo enviar el correo SMTP a %s: %v", toEmail, err)
			// Returning nil so the API response still works in dev even if SMTP server connection fails
		} else {
			log.Printf("[EmailService Success] Correo enviado correctamente a %s", toEmail)
		}
	}

	return nil
}
