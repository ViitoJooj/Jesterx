Create a production-grade security middleware stack for a Go REST API.

Requirements:

1. Global rate limiting
- Limit requests per IP
- Example: 10 requests per second with burst
- Return HTTP 429 when limit exceeded

2. Request body protection
- Limit request body size to prevent large payload abuse
- Example: maximum 10MB
- Reject with proper HTTP status

3. Timeout protection
- Configure HTTP server timeouts:
  - ReadTimeout
  - WriteTimeout
  - IdleTimeout

4. Abuse protection
- Track IP address
- Basic in-memory IP blocking when too many violations occur
- Temporary ban system

5. Logging middleware
- Log:
  - IP
  - request method
  - endpoint
  - response status
  - response time

6. Security headers
Add headers:
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection
- Strict-Transport-Security
- Content-Security-Policy

7. Upload protection
- Prevent large file uploads
- Validate Content-Length

8. Panic recovery middleware
- Recover from panic
- Return HTTP 500 safely
- Log the stack trace

9. Request size and multipart form limits

10. Efficient and low-memory implementation
- Avoid heavy dependencies
- Prefer standard library where possible

11. Modular design
- Each protection implemented as middleware
- Easy to plug into net/http router

12. Optional:
- simple in-memory rate limiter using golang.org/x/time/rate
- middleware chaining example

Output:
- clean Go code
- production ready
- comments explaining each protection

13. Protect against denial-of-wallet attacks
- prevent excessive downloads
- limit expensive endpoints
- protect CPU intensive routes

14. Add per-route rate limiting support

Create a complete security layer for a Go REST API.

Include protections for:

- global rate limiting
- per IP rate limiting
- per route rate limiting
- brute force protection
- user enumeration protection
- request body size limit
- JSON payload protection
- upload size limit
- timeout protection
- slowloris protection
- panic recovery
- request logging
- security headers
- abuse detection
- temporary IP ban system
- protection against denial-of-wallet attacks
- pagination enforcement
- protection against mass assignment
- path traversal protection
- middleware based architecture

Use standard Go net/http and minimal dependencies.
Code must be modular and production-ready.
remove all coments

If a user deletes their account, all websites they created will also be deleted, and any products registered in that store will also be deleted.