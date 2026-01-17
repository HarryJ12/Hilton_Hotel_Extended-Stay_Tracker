package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// brevoPayload represents the JSON payload structure for Brevo email API
// Brevo is used to send transactional emails
// - For documentation, see: https://developers.brevo.com/docs/send-a-transactional-email
type brevoPayload struct {
	// Sender email address (must be authorized in Brevo)
	Sender struct {
		Email string `json:"email"`
	} `json:"sender"`
	// Recipient email addresses
	To []struct {
		Email string `json:"email"`
	} `json:"to"`
	Subject     string `json:"subject"`     // Email subject
	TextContent string `json:"textContent"` // Plain text email body
}

/*
sendEmail sends an email using the Brevo transactional email API
It requires the following environment variables to be set:
  - BREVO_API_KEY: API key for Brevo
  - EMAIL_FROM: Authorized sender email address
  - MANAGER_EMAIL: Recipient email address

Returns an error if:
- Required environment variables are missing
- JSON encoding fails
- HTTP request fails
- Brevo returns a non-2xx response
*/
func sendEmail(subject, body string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY not set")
	}
	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		from = os.Getenv("SMTP_EMAIL")
	}
	if from == "" {
		return fmt.Errorf("EMAIL_FROM not set")
	}
	to := os.Getenv("MANAGER_EMAIL")
	if to == "" {
		return fmt.Errorf("MANAGER_EMAIL not set")
	}

	// Build Brevo email
	var p brevoPayload
	p.Sender.Email = from
	p.To = append(p.To, struct {
		Email string `json:"email"`
	}{Email: to})
	p.Subject = subject
	p.TextContent = body

	// Encode as JSON, if fails return error
	jsonBody, err := json.Marshal(p)
	if err != nil {
		return err
	}

	// Create HTTP POST request to Brevo API, if fails return error
	req, err := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return err
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-key", apiKey)

	// Send the HTTP request, if fails return error
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Treat any non-2xx response as a failure
	if resp.StatusCode >= 300 {
		return fmt.Errorf("brevo send failed: %s", resp.Status)
	}

	return nil
}
