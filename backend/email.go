// package main

// import (
// 	"fmt"
// 	"net/smtp"
// 	"os"
// )

// func sendEmail(subject, body string) error {
// 	host := os.Getenv("SMTP_HOST")
// 	port := os.Getenv("SMTP_PORT")
// 	from := os.Getenv("SMTP_EMAIL")
// 	pass := os.Getenv("SMTP_PASSWORD")
// 	to := os.Getenv("MANAGER_EMAIL")

// 	if host == "" || port == "" || from == "" || pass == "" || to == "" {
// 		return fmt.Errorf("smtp or manager env vars not set")
// 	}

// 	auth := smtp.PlainAuth("", from, pass, host)

// 	msg := []byte(
// 		"To: " + to + "\r\n" +
// 			"Subject: " + subject + "\r\n" +
// 			"MIME-Version: 1.0\r\n" +
// 			"Content-Type: text/plain; charset=\"utf-8\"\r\n\r\n" +
// 			body + "\r\n",
// 	)

// 	return smtp.SendMail(host+":"+port, auth, from, []string{to}, msg)
// }

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type brevoPayload struct {
	Sender struct {
		Email string `json:"email"`
	} `json:"sender"`
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	Subject     string `json:"subject"`
	TextContent string `json:"textContent"`
}

func sendEmail(subject, body string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY not set")
	}

	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		// fallback if you keep old naming
		from = os.Getenv("SMTP_EMAIL")
	}
	if from == "" {
		return fmt.Errorf("EMAIL_FROM not set")
	}

	to := os.Getenv("MANAGER_EMAIL")
	if to == "" {
		return fmt.Errorf("MANAGER_EMAIL not set")
	}

	var p brevoPayload
	p.Sender.Email = from
	p.To = append(p.To, struct {
		Email string `json:"email"`
	}{Email: to})
	p.Subject = subject
	p.TextContent = body

	jsonBody, err := json.Marshal(p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("brevo send failed: %s", resp.Status)
	}

	return nil
}
