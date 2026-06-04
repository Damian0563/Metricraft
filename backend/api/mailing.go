package api

import (
	"backend/db"
	"backend/types"
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"net/smtp"
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
<body style="margin:0;padding:0;background-color:#0A0E13;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#ffffff;">
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">
    Your Metricraft verification code is {{CODE}}. It expires in 10 minutes.
  </div>
  <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0A0E13;padding:32px 16px;">
    <tr>
      <td align="center">
        <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;background:#0D131A;border:1px solid #1A2633;border-radius:14px;overflow:hidden;">
          <tr>
            <td style="padding:28px 32px;border-bottom:1px solid #1A2633;background:#0A0E13;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td align="left" style="font-size:20px;font-weight:500;letter-spacing:-0.3px;color:#ffffff;">
                    <span style="display:inline-block;width:10px;height:10px;border-radius:2px;background:#00F376;box-shadow:0 0 10px #00F376,0 0 20px rgba(0,243,118,0.6);margin-right:12px;vertical-align:middle;"></span>
                    <span style="color:#ffffff;">Metric</span><span style="color:#00F376;font-weight:900;">raft</span>
                  </td>
                  <td align="right" style="font-size:11px;color:#00F376;letter-spacing:1px;text-transform:uppercase;font-weight:600;">
                    Verification
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:40px 32px 8px 32px;">
              <h1 style="margin:0 0 12px 0;font-size:26px;line-height:1.3;color:#ffffff;font-weight:600;letter-spacing:-0.3px;">
                Confirm your email
              </h1>
              <p style="margin:0;font-size:15px;line-height:1.6;color:#9aa5b1;">
                Use the verification code below to finish creating your Metricraft account.
                The code expires in <strong style="color:#00F376;">10 minutes</strong>.
              </p>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:32px 32px 24px 32px;">
              <div style="
                display:inline-block;
                padding:20px 32px;
                background:#0A0E13;
                border:1px solid #00F376;
                border-radius:12px;
                font-family:'SFMono-Regular',ui-monospace,Menlo,Consolas,'Liberation Mono',monospace;
                font-size:36px;
                font-weight:700;
                letter-spacing:14px;
                color:#00F376;
                text-shadow:0 0 4px #00F376,0 0 12px rgba(0,243,118,0.7);
                box-shadow:0 0 0 1px rgba(0,243,118,0.15),0 0 24px rgba(0,243,118,0.25);
              ">
                {{CODE}}
              </div>
              <div style="margin-top:16px;font-size:11px;color:#6b7282;letter-spacing:1px;text-transform:uppercase;font-weight:600;">
                Your verification code
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:8px 32px 32px 32px;">
              <div style="background:#0A0E13;border:1px solid #1A2633;border-left:3px solid #00F376;border-radius:8px;padding:14px 18px;">
                <p style="margin:0;font-size:13px;line-height:1.6;color:#9aa5b1;">
                  Didn't request this code? You can safely ignore this email &mdash; someone may have entered your address by mistake.
                </p>
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 32px 28px 32px;border-top:1px solid #1A2633;background:#0A0E13;">
              <p style="margin:0;font-size:12px;line-height:1.6;color:#6b7282;">
                Sent by <span style="color:#ffffff;">Metric</span><span style="color:#00F376;font-weight:700;">raft</span> &middot; This is an automated message, please do not reply.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`

const permissionRequestEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Metricraft access request</title>
</head>
<body style="margin:0;padding:0;background-color:#0A0E13;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#ffffff;">
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">
    {{REQUESTER}} is requesting access to {{APP}} on Metricraft.
  </div>
  <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:#0A0E13;padding:32px 16px;">
    <tr>
      <td align="center">
        <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;background:#0D131A;border:1px solid #1A2633;border-radius:14px;overflow:hidden;">
          <tr>
            <td style="padding:28px 32px;border-bottom:1px solid #1A2633;background:#0A0E13;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0">
                <tr>
                  <td align="left" style="font-size:20px;font-weight:500;letter-spacing:-0.3px;color:#ffffff;">
                    <span style="display:inline-block;width:10px;height:10px;border-radius:2px;background:#00F376;box-shadow:0 0 10px #00F376,0 0 20px rgba(0,243,118,0.6);margin-right:12px;vertical-align:middle;"></span>
                    <span style="color:#ffffff;">Metric</span><span style="color:#00F376;font-weight:900;">raft</span>
                  </td>
                  <td align="right" style="font-size:11px;color:#00F376;letter-spacing:1px;text-transform:uppercase;font-weight:600;">
                    Access request
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:40px 32px 8px 32px;">
              <h1 style="margin:0 0 12px 0;font-size:26px;line-height:1.3;color:#ffffff;font-weight:600;letter-spacing:-0.3px;">
                Someone is awaiting permission
              </h1>
              <p style="margin:0;font-size:15px;line-height:1.6;color:#9aa5b1;">
                A user is requesting access to your Metricraft project and is waiting for you to grant them permission.
              </p>
            </td>
          </tr>
          <tr>
            <td style="padding:32px 32px 24px 32px;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background:#0A0E13;border:1px solid #1A2633;border-radius:12px;">
                <tr>
                  <td style="padding:18px 22px;border-bottom:1px solid #1A2633;">
                    <div style="font-size:11px;color:#6b7282;letter-spacing:1px;text-transform:uppercase;font-weight:600;margin-bottom:6px;">
                      Requested by
                    </div>
                    <div style="font-size:16px;color:#ffffff;font-weight:500;word-break:break-all;">
                      {{REQUESTER}}
                    </div>
                  </td>
                </tr>
                <tr>
                  <td style="padding:18px 22px;">
                    <div style="font-size:11px;color:#6b7282;letter-spacing:1px;text-transform:uppercase;font-weight:600;margin-bottom:6px;">
                      Project
                    </div>
                    <div style="font-size:16px;color:#00F376;font-weight:600;word-break:break-all;text-shadow:0 0 8px rgba(0,243,118,0.4);">
                      {{APP}}
                    </div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:0 32px 32px 32px;">
              <div style="background:#0A0E13;border:1px solid #1A2633;border-left:3px solid #00F376;border-radius:8px;padding:14px 18px;">
                <p style="margin:0;font-size:13px;line-height:1.6;color:#9aa5b1;">
                  Open your Metricraft dashboard to review this request and approve or deny access. If you don't recognise this user, you can safely ignore this email.
                </p>
              </div>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 32px 28px 32px;border-top:1px solid #1A2633;background:#0A0E13;">
              <p style="margin:0;font-size:12px;line-height:1.6;color:#6b7282;">
                Sent by <span style="color:#ffffff;">Metric</span><span style="color:#00F376;font-weight:700;">raft</span> &middot; This is an automated message, please do not reply.
              </p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`

func renderVerificationEmail(to, code string) []byte {
	// Sanitize to prevent header injection
	to = strings.ReplaceAll(strings.ReplaceAll(to, "\r", ""), "\n", "")
	body := strings.ReplaceAll(verificationEmailTemplate, "{{CODE}}", html.EscapeString(code))
	headers := "From: Metricraft <noreply.metricraft@gmail.com>\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Your Metricraft verification code\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n" +
		"\r\n"
	return []byte(headers + body)
}

func permissionRequestCreator(to, requester, appName string) []byte {
	body := strings.ReplaceAll(permissionRequestEmailTemplate, "{{REQUESTER}}", html.EscapeString(requester))
	body = strings.ReplaceAll(body, "{{APP}}", html.EscapeString(appName))
	headers := "From: Metricraft <noreply.metricraft@gmail.com>\r\n" +
		"To: " + to + "\r\n" +
		"Subject: Someone is awaiting access to " + appName + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n" +
		"\r\n"
	return []byte(headers + body)
}

func sendMail(mail string, code string) error {
	apiKey := os.Getenv("GOOGLE_APP_PASSWORD")
	if apiKey == "" {
		return errors.New("GOOGLE_APP_PASSWORD")
	}
	auth := smtp.PlainAuth("", "noreply.metricraft@gmail.com", apiKey, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "noreply.metricraft@gmail.com", []string{mail}, renderVerificationEmail(mail, code))
	return err
}

func SendVerification(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	type sendVerificationPayload struct {
		Mail string `json:"mail"`
	}
	var payload sendVerificationPayload
	json.NewDecoder(r.Body).Decode(&payload)
	var routine = make(chan types.ExistsErrResponse, 1)
	mail := payload.Mail
	go db.CheckUserExists(routine, mail)
	response := <-routine
	if response.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if response.Exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := db.GenerateCode()
	err := db.SetCodeValidity(mail, code)
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

func sendPermissionRequest(owner string, mail string, appName string) error {
	apiKey := os.Getenv("GOOGLE_APP_PASSWORD")
	if apiKey == "" {
		return errors.New("GOOGLE_APP_PASSWORD")
	}
	auth := smtp.PlainAuth("", "noreply.metricraft@gmail.com", apiKey, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "noreply.metricraft@gmail.com", []string{mail}, permissionRequestCreator(owner, mail, appName))
	return err
}

func CheckVerification(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	type checkVerificationPayload struct {
		AppName string `json:"appName"`
		Mail    string `json:"mail"`
		Code    string `json:"code"`
	}
	var payload checkVerificationPayload
	json.NewDecoder(r.Body).Decode(&payload)
	var routine = make(chan types.ExistsErrResponse, 2)
	go db.CheckAllowed(routine, payload.Mail, payload.AppName)
	go db.CheckCodeValidity(routine, payload.Mail, payload.Code)
	var (
		codeValid   bool
		permitted   bool
		internalErr bool
		malicious   bool
	)
	for received := 2; received > 0; received-- {
		select {
		case response := <-routine:
			switch response.Origin {
			case "checkCodeValidity":
				if response.Err != nil {
					internalErr = true
				} else {
					codeValid = response.Exists
				}
			case "checkAllowed":
				msg := ""
				if response.Err != nil {
					msg = response.Err.Error()
				}
				switch {
				case msg == "Owner is allowed to sign in", response.Err == nil && response.Exists:
					permitted = true
				case msg == "Permission needed from the owner", response.Err == nil && !response.Exists:
					if err := sendPermissionRequest(response.Owner, payload.Mail, payload.AppName); err != nil {
						internalErr = true
					} else {
						malicious = true
					}
				case msg == "App name verification needed.", response.Err == nil && response.Exists:
					permitted = verifyAppName(payload.AppName)
				default:
					internalErr = true
				}
			}
		}
	}
	switch {
	case internalErr:
		http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
	case !codeValid:
		http.Error(w, "Invalid or expired verification code", http.StatusBadRequest)
	case !permitted:
		http.Error(w, "Permission needed from the owner", http.StatusUnauthorized)
	case malicious:
		http.Error(w, "Invalid app name", http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verification successful"))
	}
}
