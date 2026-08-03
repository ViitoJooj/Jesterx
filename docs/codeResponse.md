# General
- `RAX-001` -> internal error.
- `RAX-002` -> service unavailable.
- `RAX-003` -> invalid request.
- `RAX-004` -> malformed body.
- `RAX-005` -> resource not found.
- `RAX-006` -> operation not allowed.
- `RAX-007` -> request timeout.
- `RAX-008` -> too many requests.
- `RAX-009` -> feature unavailable.
- `RAX-010` -> unknown error.

# Authentication
- `RBX-001` -> invalid credentials.
- `RBX-002` -> invalid email.
- `RBX-003` -> invalid password.
- `RBX-004` -> password too weak.
- `RBX-005` -> email already exists.
- `RBX-006` -> email not verified.
- `RBX-007` -> invalid verification code.
- `RBX-008` -> verification code expired.
- `RBX-009` -> invalid token.
- `RBX-010` -> token expired.
- `RBX-011` -> refresh token expired.
- `RBX-012` -> unauthorized.
- `RBX-013` -> forbidden.
- `RBX-014` -> invalid email.
- `RBX-015` -> account disabled.
- `RBX-016` -> account locked.
- `RBX-017` -> too many login attempts.
- `RBX-018` -> invalid session.
- `RBX-019` -> session expired.
- `RBX-020` -> logout required.

# User
- `RCX-001` -> user not found.
- `RCX-002` -> user already exists.
- `RCX-003` -> invalid user identifier.
- `RCX-004` -> invalid name.
- `RCX-005` -> invalid cpf.
- `RCX-006` -> invalid avatar.
- `RCX-007` -> invalid role.
- `RCX-008` -> role not allowed.
- `RCX-009` -> profile incomplete.
- `RCX-010` -> user inactive.

# Website
- `RDX-001` -> website not found.
- `RDX-002` -> website already exists.
- `RDX-003` -> invalid website identifier.
- `RDX-004` -> invalid domain.
- `RDX-005` -> website disabled.
- `RDX-006` -> website limit reached.

# Database
- `RSI-001` -> database error.
- `RSI-002` -> duplicate value.
- `RSI-003` -> foreign key violation.
- `RSI-004` -> record not found.
- `RSI-005` -> transaction failed.
- `RSI-006` -> connection failed.
- `RSI-007` -> query failed.
- `RSI-008` -> constraint violation.

# Validation
- `RDI-001` -> invalid uuid.
- `RDI-002` -> invalid input.
- `RDI-003` -> missing required field.
- `RDI-004` -> invalid length.
- `RDI-005` -> value out of range.
- `RDI-006` -> invalid format.
- `RDI-007` -> invalid json.
- `RDI-008` -> unsupported value.

# Files
- `R8-001` -> file not found.
- `R8-002` -> file too large.
- `R8-003` -> invalid file type.
- `R8-004` -> upload failed.
- `R8-005` -> download failed.
- `R8-006` -> image processing failed.

# Email
- `R9-001` -> email delivery failed.
- `R9-002` -> invalid recipient.
- `R9-003` -> template not found.
- `R9-004` -> smtp unavailable.
- `R9-005` -> email already verified.

# Permission
- `R10-001` -> permission denied.
- `R10-002` -> insufficient permissions.
- `R10-003` -> administrator required.
- `R10-004` -> owner required.

# Products
- `R11-001` -> product not found.
- `R11-002` -> product already exists.
- `R11-003` -> product unavailable.
- `R11-004` -> out of stock.
- `R11-005` -> invalid price.
- `R11-006` -> invalid quantity.
- `R11-007` -> stock limit exceeded.

# Orders
- `R12-001` -> order not found.
- `R12-002` -> order already completed.
- `R12-003` -> order already canceled.
- `R12-004` -> payment required.
- `R12-005` -> payment failed.
- `R12-006` -> refund failed.
- `R12-007` -> insufficient balance.

# Rate Limits
- `R13-001` -> rate limit exceeded.
- `R13-002` -> temporarily blocked.
- `R13-003` -> ip blocked.

# API
- `R14-001` -> invalid api key.
- `R14-002` -> api key expired.
- `R14-003` -> api key revoked.
- `R14-004` -> endpoint disabled.
- `R14-005` -> unsupported api version.