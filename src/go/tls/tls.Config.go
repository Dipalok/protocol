// package tls

// import (
// 	"bufio"
// 	"context"
// 	"crypto/tls"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strings"

// 	_ "github.com/jackc/pgx/v5/stdlib"
// 	"github.com/kovarike/protocol/src/go/db"
// )

// // Variáveis globais para conexões

// func ListenTLS(addr, certFile, keyFile string) error {
// 	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
// 	if err != nil {
// 		return err
// 	}
// 	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
// 	ln, err := tls.Listen("tcp", addr, cfg)
// 	if err != nil {
// 		return err
// 	}
// 	defer ln.Close()
// 	log.Println("TLS server listening on", addr)
// 	// err = db.InitDB()
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// worker.Init(db.Pg, db.Rdb)

// 	// ctx := context.Background()
// 	// modemPort := "/dev/ttyUSB0"
// 	// go worker.WorkerLoop(ctx, modemPort)
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			log.Println("accept:", err)
// 			continue
// 		}
// 		go handleConnection(conn)
// 	}
// }

// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	reader := bufio.NewReader(conn)
// 	writer := bufio.NewWriter(conn)

// 	writeResponse := func(code int, message string) {
// 		response := fmt.Sprintf("%d %s\r\n", code, message)
// 		writer.WriteString(response)
// 		writer.Flush()
// 	}

// 	writeResponse(220, "Unified Notify Protocol Server Ready")

// 	authenticated := false
// 	var mailFrom, rcptTo string
// 	var currentChannel string // "email" ou "sms"
// 	var dataBuffer []string   // armazena linhas DATA

// 	for {
// 		line, err := reader.ReadString('\n')
// 		if err != nil {
// 			log.Println("Connection closed:", err)
// 			return
// 		}
// 		line = strings.TrimSpace(line)
// 		log.Println("Received:", line)

// 		parts := strings.SplitN(line, " ", 2)
// 		command := strings.ToUpper(parts[0])
// 		var param string
// 		if len(parts) > 1 {
// 			param = parts[1]
// 		}

// 		switch command {
// 		case "EHLO":
// 			writeResponse(250, "Hello "+param)
// 		case "AUTH":
// 			// Aqui você pode implementar autenticação real
// 			authenticated = true
// 			writeResponse(250, "Authentication successful")
// 		case "MAIL":
// 			if !authenticated {
// 				writeResponse(530, "Authentication required")
// 				continue
// 			}
// 			// Exemplo: MAIL FROM:<user@example.com>
// 			mailFrom = extractEmail(param)
// 			currentChannel = "email" // assumindo email
// 			writeResponse(250, "OK")
// 		case "RCPT":
// 			if !authenticated {
// 				writeResponse(530, "Authentication required")
// 				continue
// 			}
// 			// Exemplo: RCPT TO:<recipient@example.com>
// 			rcptTo = extractEmail(param)
// 			writeResponse(250, "OK")
// 		case "DATA":
// 			if !authenticated {
// 				writeResponse(530, "Authentication required")
// 				continue
// 			}
// 			writeResponse(354, "Start input; end with <CRLF>.<CRLF>")
// 			dataBuffer = []string{}
// 			for {
// 				dataLine, err := reader.ReadString('\n')
// 				if err != nil {
// 					log.Println("Error reading data:", err)
// 					return
// 				}
// 				dataLine = strings.TrimRight(dataLine, "\r\n")
// 				if dataLine == "." {
// 					break
// 				}
// 				dataBuffer = append(dataBuffer, dataLine)
// 			}
// 			payload := strings.Join(dataBuffer, "\n")
// 			log.Printf("Received data for email: %s\n", payload)

// 			// Aqui chama sua função para salvar e enfileirar
// 			id, err := saveAndEnqueue(currentChannel, mailFrom, rcptTo, payload)
// 			if err != nil {
// 				log.Printf("Erro ao salvar mensagem: %v", err)
// 				writeResponse(451, "Requested action aborted: error saving message")
// 				continue
// 			}
// 			log.Printf("Mensagem salva com ID: %s\n", id)
// 			writeResponse(250, "Message accepted")
// 		case "SMS":
// 			if !authenticated {
// 				writeResponse(530, "Authentication required")
// 				continue
// 			}
// 			// Exemplo simples: SMS FROM:<+5511999999999>
// 			mailFrom = extractEmail(param) // aqui email é o telefone, adaptar se quiser
// 			currentChannel = "sms"
// 			writeResponse(250, "OK")
// 		case "SMS DATA":
// 			if !authenticated {
// 				writeResponse(530, "Authentication required")
// 				continue
// 			}
// 			writeResponse(354, "Start input; end with <CRLF>.<CRLF>")
// 			dataBuffer = []string{}
// 			for {
// 				dataLine, err := reader.ReadString('\n')
// 				if err != nil {
// 					log.Println("Error reading SMS data:", err)
// 					return
// 				}
// 				dataLine = strings.TrimRight(dataLine, "\r\n")
// 				if dataLine == "." {
// 					break
// 				}
// 				dataBuffer = append(dataBuffer, dataLine)
// 			}
// 			payload := strings.Join(dataBuffer, "\n")
// 			log.Printf("Received data for SMS: %s\n", payload)

// 			id, err := saveAndEnqueue(currentChannel, mailFrom, rcptTo, payload)
// 			if err != nil {
// 				log.Printf("Erro ao salvar SMS: %v", err)
// 				writeResponse(451, "Requested action aborted: error saving SMS")
// 				continue
// 			}
// 			log.Printf("SMS salvo com ID: %s\n", id)
// 			writeResponse(250, "SMS accepted")
// 		case "QUIT":
// 			writeResponse(221, "Bye")
// 			return
// 		default:
// 			writeResponse(500, "Command unrecognized")
// 		}
// 	}
// }

// // Função para extrair email ou telefone do comando MAIL FROM, RCPT TO, SMS FROM
// func extractEmail(param string) string {
// 	param = strings.TrimSpace(param)
// 	param = strings.TrimPrefix(param, "FROM:")
// 	param = strings.TrimPrefix(param, "TO:")
// 	param = strings.Trim(param, "<>")
// 	return param
// }

// // salvar mensagem e enfileirar
// func saveAndEnqueue(channel, from, to string, payload string) (string, error) {
// 	var id string
// 	err := db.Pg.QueryRow(
// 		`INSERT INTO messages(channel, sender, recipient, payload) VALUES($1,$2,$3,$4) RETURNING id`,
// 		channel, from, to, payload,
// 	).Scan(&id)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = db.Rdb.LPush(context.Background(), "queue:notifications", id).Err()
// 	if err != nil {
// 		return id, err
// 	}

// 	return id, nil
// }

package tls

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kovarike/protocol/src/go/db"
)

// func ListenTLS(addr, certFile, keyFile string) error {
// 	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
// 	if err != nil {
// 		return fmt.Errorf("failed to load key pair: %w", err)
// 	}

// 	cfg := &tls.Config{
// 		Certificates: []tls.Certificate{cert},
// 		MinVersion:   tls.VersionTLS12,
// 		ClientAuth:   tls.NoClientCert, // For development only
// 	}

// 	ln, err := tls.Listen("tcp", addr, cfg)
// 	if err != nil {
// 		return fmt.Errorf("failed to listen: %w", err)
// 	}
// 	defer ln.Close()

// 	log.Println("TLS server listening on", addr)
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			log.Println("accept:", err)
// 			continue
// 		}
// 		go handleConnection(conn)
// 	}
// }

func ListenTLS(addr, certFile, keyFile string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load key pair: %w", err)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	defer ln.Close()

	log.Println("TLS server listening on", addr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writeResponse := func(code int, message string) {
		response := fmt.Sprintf("%d %s\r\n", code, message)
		writer.WriteString(response)
		writer.Flush()
	}

	writeResponse(220, "Unified Notify Protocol Server Ready")

	authenticated := false
	var mailFrom, rcptTo string
	var currentChannel string
	var dataBuffer []string

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		line = strings.TrimSpace(line)
		log.Println("Received:", line)

		parts := strings.SplitN(line, " ", 2)
		command := strings.ToUpper(parts[0])
		var param string
		if len(parts) > 1 {
			param = parts[1]
		}

		switch command {
		case "EHLO":
			writeResponse(250, "Hello "+param)
		case "AUTH":
			// Implement real authentication here
			authenticated = true
			writeResponse(250, "Authentication successful")
		case "MAIL":
			if !authenticated {
				writeResponse(530, "Authentication required")
				continue
			}
			mailFrom = extractEmail(param)
			currentChannel = "email"
			writeResponse(250, "OK")
		case "RCPT":
			if !authenticated {
				writeResponse(530, "Authentication required")
				continue
			}
			rcptTo = extractEmail(param)
			writeResponse(250, "OK")
		case "DATA":
			if !authenticated {
				writeResponse(530, "Authentication required")
				continue
			}
			writeResponse(354, "Start input; end with <CRLF>.<CRLF>")
			dataBuffer = []string{}
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Error reading data:", err)
					return
				}
				dataLine = strings.TrimRight(dataLine, "\r\n")
				if dataLine == "." {
					break
				}
				dataBuffer = append(dataBuffer, dataLine)
			}
			payload := strings.Join(dataBuffer, "\n")
			log.Printf("Received data for email: %s\n", payload)

			id, err := saveAndEnqueue(currentChannel, mailFrom, rcptTo, payload)
			if err != nil {
				log.Printf("Erro ao salvar mensagem: %v", err)
				writeResponse(451, "Requested action aborted: error saving message")
				continue
			}
			log.Printf("Mensagem salva com ID: %s\n", id)
			writeResponse(250, "Message accepted")
		case "SMS":
			if !authenticated {
				writeResponse(530, "Authentication required")
				continue
			}
			mailFrom = extractEmail(param)
			currentChannel = "sms"
			writeResponse(250, "OK")
		case "SMSDATA":
			if !authenticated {
				writeResponse(530, "Authentication required")
				continue
			}
			writeResponse(354, "Start input; end with <CRLF>.<CRLF>")
			dataBuffer = []string{}
			for {
				dataLine, err := reader.ReadString('\n')
				if err != nil {
					log.Println("Error reading SMS data:", err)
					return
				}
				dataLine = strings.TrimRight(dataLine, "\r\n")
				if dataLine == "." {
					break
				}
				dataBuffer = append(dataBuffer, dataLine)
			}
			payload := strings.Join(dataBuffer, "\n")
			log.Printf("Received data for SMS: %s\n", payload)

			id, err := saveAndEnqueue(currentChannel, mailFrom, rcptTo, payload)
			if err != nil {
				log.Printf("Erro ao salvar SMS: %v", err)
				writeResponse(451, "Requested action aborted: error saving SMS")
				continue
			}
			log.Printf("SMS salvo com ID: %s\n", id)
			writeResponse(250, "SMS accepted")
		case "QUIT":
			writeResponse(221, "Bye")
			return
		default:
			writeResponse(500, "Command unrecognized")
		}
	}
}

func extractEmail(param string) string {
	param = strings.TrimSpace(param)
	if idx := strings.Index(param, ":"); idx != -1 {
		param = param[idx+1:]
	}
	return strings.Trim(param, "<> ")
}

func saveAndEnqueue(channel, from, to string, payload string) (string, error) {
	var id string
	err := db.Pg.QueryRow(
		`INSERT INTO messages(channel, sender, recipient, payload) VALUES($1,$2,$3,$4) RETURNING id`,
		channel, from, to, payload,
	).Scan(&id)
	if err != nil {
		return "", err
	}

	err = db.Rdb.LPush(context.Background(), "queue:notifications", id).Err()
	return id, err
}
