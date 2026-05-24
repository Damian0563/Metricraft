package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/resend/resend-go/v3"
	"html"
	"net/http"
	"os"
	"strings"
)

const verificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Metricraft verification code</title>
</head>
<body style="margin:0;padding:0;background-color:#0b0d12;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#e6e8ee;">
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">
    Your Metricraft verification code is {{CODE}}. It expires in 10 minutes.
  </div>
  <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0b0d12;padding:32px 16px;">
    <tr>
      <td align="center">
        <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;background:#11141b;border:1px solid #1f2530;border-radius:14px;overflow:hidden;">
          <tr>
            <td style="padding:28px 32px;border-bottom:1px solid #1f2530;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td align="left" style="font-size:18px;font-weight:600;letter-spacing:0.2px;color:#ffffff;">
                    <span style="display:inline-block;width:10px;height:10px;border-radius:3px;background:linear-gradient(135deg,#7c5cff,#3aa9ff);margin-right:10px;vertical-align:middle;"></span>
                    Metricraft
                  </td>
                  <td align="right" style="font-size:12px;color:#7a8190;letter-spacing:0.4px;text-transform:uppercase;">
                    Account verification
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:36px 32px 8px 32px;">
              <h1 style="margin:0 0 12px 0;font-size:24px;line-height:1.3;color:#ffffff;font-weight:600;">
                Confirm your email
              </h1>
              <p style="margin:0;font-size:15px;line-height:1.6;color:#aab1c0;">
                Use the verification code below to finish creating your Metricraft account.
                The code expires in <strong style="color:#e6e8ee;">10 minutes</strong>.
              </p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:28px 32px;">
              <div style="
                display:inline-block;
                padding:18px 28px;
                background:linear-gradient(180deg,#1a1f2c 0%,#141826 100%);
                border:1px solid #2a3142;
                border-radius:12px;
                font-family:'SFMono-Regular',ui-monospace,Menlo,Consolas,'Liberation Mono',monospace;
                font-size:34px;
                font-weight:700;
                letter-spacing:14px;
                color:#ffffff;
                text-shadow:0 0 18px rgba(124,92,255,0.45);
              ">
                {{CODE}}
              </div>
              <div style="margin-top:14px;font-size:12px;color:#7a8190;letter-spacing:0.4px;text-transform:uppercase;">
                Your verification code
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:8px 32px 32px 32px;">
              <div style="background:#0d1018;border:1px solid #1f2530;border-left:3px solid #7c5cff;border-radius:8px;padding:14px 16px;">
                <p style="margin:0;font-size:13px;line-height:1.6;color:#aab1c0;">
                  Didn't request this code? You can safely ignore this email &mdash; someone may have entered your address by mistake.
                </p>
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 32px 28px 32px;border-top:1px solid #1f2530;">
              <p style="margin:0;font-size:12px;line-height:1.6;color:#6b7282;">
                Sent by Metricraft &middot; This is an automated message, please do not reply.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`

func renderVerificationEmail(code string) string {
	return strings.ReplaceAll(verificationEmailTemplate, "{{CODE}}", html.EscapeString(code))
}

func sendMail(mail string, code string) error {
	apiKey := os.Getenv("API_RESEND")
	if apiKey == "" {
		return errors.New("Missing API_RESEND")
	}
	client := resend.NewClient(apiKey)
	htmlBody := renderVerificationEmail(code)
	textBody := fmt.Sprintf(
		"Your Metricraft verification code is: %s\n\nThis code expires in 10 minutes.\nIf you didn't request it, you can safely ignore this email.",
		code,
	)
	params := &resend.SendEmailRequest{
		From:    "metricraft-noreply@gmail.com",
		To:      []string{mail},
		Html:    htmlBody,
		Text:    textBody,
		Subject: fmt.Sprintf("Your Metricraft verification code: %s", code),
	}
	_, err := client.Emails.Send(params)
	fmt.Println(err)
	return err
}

func sendVerification(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	type sendVerificationPayload struct {
		Mail string `json:"mail"`
	}
	var payload sendVerificationPayload
	json.NewDecoder(r.Body).Decode(&payload)
	exists, err := checkUserExists(payload.Mail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println(payload)
	mail := payload.Mail
	code := generateCode()
	err = setCodeValidity(mail, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = sendMail(mail, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
