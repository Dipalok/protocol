package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
	"github.com/redis/go-redis/v9"
	"github.com/tarm/serial"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Inicializa conexões ao Postgres e Redis (exemplo)
var (
	pg  *sql.DB
	rdb *redis.Client
)

func Init(pgDB *sql.DB, rdbClient *redis.Client) {
	pg = pgDB
	rdb = rdbClient
}

// Payloads para diferentes canais
type EmailPayload struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type SmsPayload struct {
	Message string `json:"message"`
}

// Envia SMS via modem serial (AT commands)
func sendSmsModem(portName, sender, recipient, message string) error {
	config := &serial.Config{Name: portName, Baud: 115200, ReadTimeout: 5 * time.Second}
	s, err := serial.OpenPort(config)
	if err != nil {
		return fmt.Errorf("open serial port: %w", err)
	}
	defer s.Close()

	cmds := []string{
		"AT\r",        // Teste conexão
		"AT+CMGF=1\r", // Modo texto SMS
		fmt.Sprintf(`AT+CMGS="%s"`+"\r", recipient), // Número destinatário
	}

	for _, cmd := range cmds {
		if _, err := s.Write([]byte(cmd)); err != nil {
			return fmt.Errorf("write command %q: %w", cmd, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Enviar mensagem + Ctrl+Z (char 26)
	if _, err := s.Write([]byte(message + string(rune(26)))); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	// Espera envio
	time.Sleep(5 * time.Second)
	return nil
}

// Envia e-mail via SMTP (exemplo SendGrid)
func sendEmail(sender, recipient string, payload EmailPayload) error {
	e := email.NewEmail()
	e.From = sender
	e.To = []string{recipient}
	e.Subject = payload.Subject
	e.Text = []byte(payload.Body)

	auth := smtp.PlainAuth("", "apikey", "SENDGRID_API_KEY", "smtp.sendgrid.net")
	err := e.Send("smtp.sendgrid.net:587", auth)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// Atualiza status da mensagem no banco
func updateStatus(id string, err error) {
	if err == nil {
		_, _ = pg.Exec(`UPDATE messages SET status='sent', updated_at=now() WHERE id=$1`, id)
	} else {
		_, _ = pg.Exec(`UPDATE messages SET status='failed', attempts = attempts + 1, updated_at=now() WHERE id=$1`, id)
		// Pode implementar fila com delay (sorted set no Redis) aqui para retry backoff
	}
}

// Loop principal do worker processando notificações
func WorkerLoop(ctx context.Context, modemPort string) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker loop finalizado por contexto cancelado")
			return
		default:
		}

		// BRPOP bloqueante (timeout 0 = espera indefinidamente)
		res, err := rdb.BRPop(ctx, 0, "queue:notifications").Result()
		if err != nil {
			log.Println("Erro no BRPop do Redis:", err)
			continue
		}
		id := res[1]

		// Buscar dados no Postgres
		var channel, sender, recipient string
		var payloadJson []byte
		err = pg.QueryRowContext(ctx, `SELECT channel,sender,recipient,payload FROM messages WHERE id=$1`, id).
			Scan(&channel, &sender, &recipient, &payloadJson)
		if err != nil {
			log.Println("Erro ao buscar mensagem no Postgres:", err)
			continue
		}

		switch channel {
		case "email":
			var payload EmailPayload
			if err := json.Unmarshal(payloadJson, &payload); err != nil {
				log.Println("Erro no unmarshal payload email:", err)
				updateStatus(id, err)
				continue
			}
			err = sendEmail(sender, recipient, payload)
			updateStatus(id, err)

		case "sms":
			var payload SmsPayload
			if err := json.Unmarshal(payloadJson, &payload); err != nil {
				log.Println("Erro no unmarshal payload sms:", err)
				updateStatus(id, err)
				continue
			}
			err = sendSmsModem(modemPort, sender, recipient, payload.Message)
			updateStatus(id, err)

		default:
			log.Printf("Canal desconhecido: %s\n", channel)
			updateStatus(id, fmt.Errorf("canal desconhecido"))
		}
	}
}
