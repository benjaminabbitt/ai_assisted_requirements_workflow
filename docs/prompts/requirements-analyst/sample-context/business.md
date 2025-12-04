# Business Context (Sample)

## Business Rules

| ID | Rule | Exceptions |
|----|------|------------|
| BR-001 | All security events must be audited | None |
| BR-002 | Password reset tokens expire after 24 hours | None |
| BR-003 | Password must meet complexity: min 8 chars, 1 uppercase, 1 number, 1 special | Admin passwords require 12+ chars |
| BR-004 | Rate limiting: 5 password reset requests per email per 15 minutes | None |

## Compliance Requirements

| Requirement | Standard | Affected Features |
|-------------|----------|-------------------|
| Security audit logging | SOC 2 | All authentication, password changes |
| Data encryption at rest | GDPR | User credentials, tokens |
| Email verification | Internal Policy | All password reset flows |

## User Personas

| Persona | Description | Access Level |
|---------|-------------|--------------|
| Registered User | Standard user account | Basic access |
| Admin | System administrator | Full access |
| Support Staff | Customer support | Limited admin access |
