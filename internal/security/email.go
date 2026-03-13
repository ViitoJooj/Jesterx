package security

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/resend/resend-go/v3"
)

func SendTicketResponseEmail(to, reporterName string, ticketNumber int, adminResponse, status string) error {
	statusLabel := map[string]string{
		"OPEN":        "Aberto",
		"IN_PROGRESS": "Em Análise",
		"RESOLVED":    "Resolvido",
		"DISMISSED":   "Encerrado",
	}
	statusText := statusLabel[status]
	if statusText == "" {
		statusText = status
	}

	html := `<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
	<meta charset="UTF-8" />
	<title>Atualização da sua denúncia</title>
	<style>
	body { margin:0; background:#1a1a1a; font-family:'Segoe UI',sans-serif; color:#f1f1f1; }
	.wrap { max-width:520px; margin:40px auto; background:#242424; border-radius:16px; overflow:hidden; }
	.header { background:linear-gradient(90deg,#ff3e00,#ff8a00); padding:28px 32px; }
	.header h1 { margin:0; font-size:22px; color:#fff; }
	.body { padding:28px 32px; }
	.ticket { background:#2e2e2e; border-radius:10px; padding:16px 20px; margin-bottom:20px; }
	.ticket span { font-size:13px; opacity:.7; }
	.ticket strong { display:block; font-size:18px; margin-top:4px; color:#ff8a00; }
	.status-badge { display:inline-block; padding:4px 14px; border-radius:20px; font-size:13px; font-weight:600; background:#ff3e00; color:#fff; margin-bottom:20px; }
	.response-box { background:#2e2e2e; border-left:3px solid #ff8a00; border-radius:8px; padding:16px 20px; font-size:14px; line-height:1.6; }
	.footer { padding:16px 32px; font-size:12px; opacity:.5; }
	</style>
	</head>
	<body>
	<div class="wrap">
	  <div class="header"><h1>JesterX · Atualização de Denúncia</h1></div>
	  <div class="body">
	    <p>Olá, <strong>` + reporterName + `</strong>!</p>
	    <p>Sua denúncia recebeu uma atualização da nossa equipe.</p>
	    <div class="ticket">
	      <span>Número do ticket</span>
	      <strong>#` + fmt.Sprintf("%05d", ticketNumber) + `</strong>
	    </div>
	    <div class="status-badge">` + statusText + `</div>
	    <p style="margin-bottom:10px;font-size:14px;opacity:.8;">Resposta da equipe JesterX:</p>
	    <div class="response-box">` + adminResponse + `</div>
	  </div>
	  <div class="footer">JesterX · Plataforma de lojas online</div>
	</div>
	</body>
	</html>`

	client := resend.NewClient(config.ResendKey)
	params := &resend.SendEmailRequest{
		From:    "JesterX <jesterx@resend.dev>",
		To:      []string{to},
		Subject: fmt.Sprintf("Atualização do ticket #%05d – JesterX", ticketNumber),
		Html:    html,
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Println("error sending ticket response email:", err)
		return errors.New("Internal error")
	}
	log.Println("Ticket response email sent! ID:", sent.Id)
	return nil
}

func SendSalesDigestEmail(to, subject, htmlBody string) error {
	client := resend.NewClient(config.ResendKey)
	params := &resend.SendEmailRequest{
		From:    "JesterX <jesterx@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Println("error sending sales digest email:", err)
		return errors.New("Internal error")
	}
	log.Println("Sales digest email sent! ID:", sent.Id)
	return nil
}

func SendOrderNotificationEmail(to, subject, htmlBody string) error {
	if config.ResendKey == "" {
		if config.IsDev {
			return nil
		}
		return errors.New("email service not configured")
	}
	client := resend.NewClient(config.ResendKey)
	params := &resend.SendEmailRequest{
		From:    "JesterX <jesterx@resend.dev>",
		To:      []string{to},
		Subject: subject,
		Html:    htmlBody,
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Println("error sending order notification email:", err)
		return errors.New("Internal error")
	}
	log.Println("Order notification email sent! ID:", sent.Id)
	return nil
}

func SendVerifyEmail(email string, token string, websiteID string) error {
	verifyURL := config.BackendURL + "/api/v1/auth/verify/" + token
	if websiteID != "" {
		verifyURL += "?website_id=" + url.QueryEscape(websiteID)
	}

	// In dev mode always log the URL so it can be used directly from the console
	// without needing real email delivery (Resend test-mode restrictions apply).
	if config.IsDev {
		log.Printf("[DEV] Verify email URL for %s → %s", email, verifyURL)
	}

	if config.ResendKey == "" {
		if config.IsDev {
			return nil // skip sending in dev when no key is set
		}
		return errors.New("email service not configured")
	}

	html := `<!DOCTYPE html>
		<html lang="pt-BR">
		<head>
		<meta charset="UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Verificar Email</title>
		<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap" rel="stylesheet">

		<style>
		body {
			margin: 0;
			height: 100vh;
			display: flex;
			justify-content: center;
			align-items: center;
			background: radial-gradient(circle at 30% 30%, #ff3e00, #ff8a00, #1a1a1a);
			font-family: 'Inter', sans-serif;
			color: white;
		}

		.container {
			text-align: center;
			max-width: 420px;
			padding: 40px;
			background: rgba(255,255,255,0.05);
			backdrop-filter: blur(20px);
			border-radius: 20px;
			box-shadow: 0 20px 40px rgba(0,0,0,0.4);
		}

		h1 {
			margin-bottom: 10px;
			font-weight: 700;
		}

		p {
			opacity: 0.8;
			margin-bottom: 30px;
			font-size: 14px;
		}

		a.button {
			display: inline-block;
			text-decoration: none;
			padding: 16px 40px;
			font-size: 16px;
			font-weight: 600;
			border: 0;
			border-radius: 14px;
			cursor: pointer;
			color: white;
			background: linear-gradient(90deg, #ff3e00, #ff8a00);
			transition: all 0.3s ease;
			box-shadow: 0 10px 30px rgba(255, 62, 0, 0.4);
		}

		a.button:hover {
			transform: translateY(-3px);
			box-shadow: 0 15px 40px rgba(255, 62, 0, 0.6);
		}

		a.button:active {
			transform: scale(0.98);
		}

		.fallback {
			margin-top: 18px;
			font-size: 12px;
			opacity: 0.65;
			word-break: break-all;
		}
		</style>
		</head>

		<body>
			<div class="container">
				<h1>Confirme seu Email</h1>
				<p>Clique no botão abaixo para ativar sua conta.</p>
				<a class="button" href="` + verifyURL + `">Verificar Email</a>
				<p class="fallback">Se o botão não funcionar, copie e cole este link no navegador:<br/>` + verifyURL + `</p>
			</div>
		</body>
		</html>`
	client := resend.NewClient(config.ResendKey)

	params := &resend.SendEmailRequest{
		From:    "JesterX <jesterx@resend.dev>",
		To:      []string{email},
		Subject: "Verify your email",
		Html:    html,
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Println("error on sending email:", err)
		return errors.New("Internal error")
	}

	log.Println("Email sended! ID:", sent.Id)
	return nil
}
