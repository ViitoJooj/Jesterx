package security

import (
	"errors"
	"log"

	"github.com/ViitoJooj/Jesterx/internal/config"
	"github.com/resend/resend-go/v3"
)

func SendVerifyEmail(email string, token string) error {
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

button {
	padding: 16px 40px;
	font-size: 16px;
	font-weight: 600;
	border: none;
	border-radius: 14px;
	cursor: pointer;
	color: white;
	background: linear-gradient(90deg, #ff3e00, #ff8a00);
	transition: all 0.3s ease;
	box-shadow: 0 10px 30px rgba(255, 62, 0, 0.4);
}

button:hover {
	transform: translateY(-3px);
	box-shadow: 0 15px 40px rgba(255, 62, 0, 0.6);
}

button:active {
	transform: scale(0.98);
}

button.loading {
	opacity: 0.7;
	pointer-events: none;
}

.message {
	margin-top: 20px;
	font-weight: 600;
	min-height: 24px;
}
.success {
	color: #00ffb3;
}
.error {
	color: #ff4d4d;
}
</style>
</head>

<body>
	<div class="container">
		<h1>Confirme seu Email</h1>
		<p>Clique no botão abaixo para ativar sua conta.</p>
		<button id="verifyBtn">Verificar Email</button>
		<div id="message" class="message"></div>
	</div>

<script>
const button = document.getElementById("verifyBtn");
const message = document.getElementById("message");

button.addEventListener("click", async () => {
	button.classList.add("loading");
	button.textContent = "Verificando...";
	message.textContent = "";

	try {
		const response = await fetch("http://localhost:8080/api/v1/auth/verify/` + token + `", {
			method: "POST"
		});

		if (response.ok) {
			button.textContent = "Email Verificado ✓";
			message.textContent = "Sua conta foi ativada com sucesso!";
			message.className = "message success";
		} else {
			button.textContent = "Verificar Email";
			message.textContent = "Token inválido ou expirado.";
			message.className = "message error";
		}
	} catch (error) {
		button.textContent = "Verificar Email";
		message.textContent = "Erro ao conectar com o servidor.";
		message.className = "message error";
	}

	button.classList.remove("loading");
});
</script>

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
