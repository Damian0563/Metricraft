package mail

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

const decisionNotificationEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Metricraft access decision</title>
</head>
<body style="margin:0;padding:0;background-color:#0A0E13;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#ffffff;">
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">
    Your access request for {{APP}} has been {{DECISION}}.
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
                    Access decision
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:40px 32px 8px 32px;">
              <h1 style="margin:0 0 12px 0;font-size:26px;line-height:1.3;color:#ffffff;font-weight:600;letter-spacing:-0.3px;">
                Your request was {{DECISION}}
              </h1>
              <p style="margin:0;font-size:15px;line-height:1.6;color:#9aa5b1;">
                The owner of this Metricraft project has reviewed your access request.
              </p>
            </td>
          </tr>
          <tr>
            <td style="padding:32px 32px 24px 32px;">
              <table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background:#0A0E13;border:1px solid #1A2633;border-radius:12px;">
                <tr>
                  <td style="padding:18px 22px;border-bottom:1px solid #1A2633;">
                    <div style="font-size:11px;color:#6b7282;letter-spacing:1px;text-transform:uppercase;font-weight:600;margin-bottom:6px;">
                      Project
                    </div>
                    <div style="font-size:16px;color:#00F376;font-weight:600;word-break:break-all;text-shadow:0 0 8px rgba(0,243,118,0.4);">
                      {{APP}}
                    </div>
                  </td>
                </tr>
                <tr>
                  <td style="padding:18px 22px;">
                    <div style="font-size:11px;color:#6b7282;letter-spacing:1px;text-transform:uppercase;font-weight:600;margin-bottom:6px;">
                      Decision
                    </div>
                    <div style="font-size:18px;color:{{DECISION_COLOR}};font-weight:700;text-transform:capitalize;text-shadow:0 0 8px {{DECISION_SHADOW}};">
                      {{DECISION}}
                    </div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:0 32px 32px 32px;">
              <div style="background:#0A0E13;border:1px solid #1A2633;border-left:3px solid {{DECISION_COLOR}};border-radius:8px;padding:14px 18px;">
                <p style="margin:0;font-size:13px;line-height:1.6;color:#9aa5b1;">
                  {{DECISION_MESSAGE}}
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

const linkEmailTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8" />
<meta name="viewport" content="width=device-width, initial-scale=1.0" />
<title>Metricraft link</title>
</head>
<body style="margin:0;padding:0;background-color:#0A0E13;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;color:#ffffff;">
  <div style="display:none;max-height:0;overflow:hidden;opacity:0;color:transparent;">
    Password Recovery
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
                    Link
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:40px 32px 8px 32px;">
              <h1 style="margin:0 0 12px 0;font-size:26px;line-height:1.3;color:#ffffff;font-weight:600;letter-spacing:-0.3px;">
                Password Recovery
              </h1>
            </td>
          </tr>
          <tr>
            <td align="center" style="padding:32px 32px 24px 32px;">
              <a href="{{LINK_URL}}" style="display:inline-block;padding:16px 28px;background:#00F376;border:1px solid #00F376;border-radius:12px;color:#0A0E13;font-size:15px;font-weight:800;text-decoration:none;letter-spacing:0.2px;box-shadow:0 0 0 1px rgba(0,243,118,0.15),0 0 24px rgba(0,243,118,0.25);">
                Reset Password
              </a>
            </td>
          </tr>
          <tr>
            <td style="padding:0 32px 32px 32px;">
              <div style="background:#0A0E13;border:1px solid #1A2633;border-left:3px solid #00F376;border-radius:8px;padding:14px 18px;">
                <p style="margin:0 0 8px 0;font-size:13px;line-height:1.6;color:#9aa5b1;">
                  If the button does not work, copy and paste this link into your browser:
                </p>
                <p style="margin:0;font-size:13px;line-height:1.6;color:#00F376;word-break:break-all;">
                  <a href="{{LINK_URL}}" style="color:#00F376;text-decoration:none;">{{LINK_URL}}</a>
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
