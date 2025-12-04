# Sample User Story

**Ticket ID:** PROJ-1234
**Ticket URL:** https://jira.example.com/browse/PROJ-1234
**Status:** Ready for Specification

## User Story

As a registered user, I want to reset my password, so that I can regain access to my account if I forget my password.

## Acceptance Criteria

- [ ] User can request password reset by providing email
- [ ] System validates email exists in database
- [ ] System sends reset email with unique token
- [ ] Reset token expires after 24 hours
- [ ] User can set new password using valid token
- [ ] Old password is invalidated after successful reset
- [ ] User receives confirmation email after password change
- [ ] System logs all password reset attempts for security audit

## Notes

- **Security:** Reset tokens must be single-use only
- **Compliance:** Must log all attempts per security policy
- **Rate Limiting:** Maximum 5 reset requests per email per 15 minutes
- **Email Service:** Uses SendGrid (configured in environment)

## Dependencies

- User authentication service (internal)
- Email service (SendGrid - external)
- Audit logging service (internal)

## Business Rules (from business.md)

- Password must meet complexity requirements (min 8 chars, 1 uppercase, 1 number, 1 special)
- All security events must be logged
- Email must be verified before account activation
