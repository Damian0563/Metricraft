package mail

import (
	"errors"
	"html"
	"net/smtp"
	"os"
	"strings"
)

const (
	fromAddress = "noreply.metricraft@gmail.com"
	fromHeader  = "Metricraft <" + fromAddress + ">"
	smtpAddress = "smtp.gmail.com:587"
	smtpHost    = "smtp.gmail.com"
)

func SendVerification(to, code string) error {
	body := strings.ReplaceAll(verificationEmailTemplate, "{{CODE}}", html.EscapeString(code))
	return send(to, "Your Metricraft verification code", body)
}

func SendPermissionRequest(owner, requester, appName string) error {
	body := strings.ReplaceAll(permissionRequestEmailTemplate, "{{REQUESTER}}", html.EscapeString(requester))
	body = strings.ReplaceAll(body, "{{APP}}", html.EscapeString(appName))
	return send(owner, "Someone is awaiting access to "+sanitizeHeaderValue(appName), body)
}

func NotifyDecision(to, decision, appName string) error {
	decisionText := "denied"
	decisionColor := "#EF4444"
	decisionShadow := "rgba(239,68,68,0.4)"
	decisionMessage := "You do not currently have access to this project. If you believe this was a mistake, contact the project owner."
	if decision == "true" {
		decisionText = "approved"
		decisionColor = "#00F376"
		decisionShadow = "rgba(0,243,118,0.4)"
		decisionMessage = "You can now sign in to this Metricraft project using your existing verification flow."
	}

	body := strings.ReplaceAll(decisionNotificationEmailTemplate, "{{APP}}", html.EscapeString(appName))
	body = strings.ReplaceAll(body, "{{DECISION}}", decisionText)
	body = strings.ReplaceAll(body, "{{DECISION_COLOR}}", decisionColor)
	body = strings.ReplaceAll(body, "{{DECISION_SHADOW}}", decisionShadow)
	body = strings.ReplaceAll(body, "{{DECISION_MESSAGE}}", html.EscapeString(decisionMessage))

	subject := "Your access request for " + sanitizeHeaderValue(appName) + " was " + decisionText
	return send(to, subject, body)
}

func SendRecovery(to, subject, linkURL string) error {
	body := strings.ReplaceAll(linkEmailTemplate, "{{LINK_URL}}", html.EscapeString(linkURL))
	return send(to, subject, body)
}

func SendInvite(to string) error {
	inviteURL := strings.TrimRight(os.Getenv("frontend"), "/")
	if inviteURL == "" {
		inviteURL = strings.TrimRight(os.Getenv("host"), "/")
	}
	body := strings.ReplaceAll(inviteEmailTemplate, "{{INVITE_URL}}", html.EscapeString(inviteURL))
	return send(to, "You have been invited to Metricraft", body)
}

func send(to, subject, body string) error {
	apiKey := os.Getenv("GOOGLE_APP_PASSWORD")
	if apiKey == "" {
		return errors.New("GOOGLE_APP_PASSWORD")
	}
	to = sanitizeHeaderValue(to)
	auth := smtp.PlainAuth("", fromAddress, apiKey, smtpHost)
	return smtp.SendMail(smtpAddress, auth, fromAddress, []string{to}, buildMessage(to, subject, body))
}

func buildMessage(to, subject, body string) []byte {
	headers := "From: " + fromHeader + "\r\n" +
		"To: " + sanitizeHeaderValue(to) + "\r\n" +
		"Subject: " + sanitizeHeaderValue(subject) + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n" +
		"\r\n"

	return []byte(headers + body)
}

func SendManualInvites(invitees []string) error {
	for _, invitee := range invitees {
		err := SendInvite(invitee)
		if err != nil {
			return err
		}
	}
	return nil
}

func sanitizeHeaderValue(value string) string {
	value = strings.ReplaceAll(value, "\r", "")
	return strings.ReplaceAll(value, "\n", "")
}
