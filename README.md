# Kannada Balaga Booking System

End-to-end documentation for the booking platform.

This repository is split into two apps:

- `booking-backend` - Go API, database, auth, booking, email verification, and admin check-in.
- `booking-frontend/ganeshbooking` - Next.js app for users and admins.

## What This App Does

- User signup and signin.
- Email verification flow.
- Ticket booking with QR generation.
- Admin dashboard for verification and check-in.
- QR scanner for gate entry.
- Optional Cloudflare Turnstile CAPTCHA for auth endpoints.
- SMTP mail sending for verification and booking emails.

## Prerequisites

- Go 1.22 or newer.
- Node.js 20 or newer.
- PostgreSQL.
- SMTP credentials for Brevo or Gmail.
- Cloudflare Turnstile site key and secret key if CAPTCHA is enabled.

## Repository Structure

```text
booking-backend/
  cmd/main.go
  models/
  router/
  scripts/
booking-frontend/ganeshbooking/
  src/app/
  src/components/
```

## Local Setup Overview

1. Configure the backend environment variables.
2. Start the Go API on port `8080`.
3. Configure the frontend environment variables.
4. Start the Next.js app on port `3000`.
5. Open the frontend and use the booking or admin flows.

## Backend Setup

### 1) Create `booking-backend/.env`

Use the following template:

```env
PORT=8080
DB_URL=postgres://postgres:postgres@localhost:5432/kannada_booking?sslmode=disable

# or use DATABASE_URL instead of DB_URL
# DATABASE_URL=postgres://postgres:postgres@localhost:5432/kannada_booking?sslmode=disable

FRONTEND_URL=http://localhost:3000
JWT_SECRET=replace_with_a_long_random_secret
SESSION_SECRET=replace_with_a_long_random_secret

Adminname=admin
Adminpass=adminpassword

# Mail settings
Email=verified-sender@example.com
Emailuser=your_smtp_username
Emailpass=your_smtp_password

# CAPTCHA settings
CAPTCHA_ENABLED=false
CAPTCHA_SECRET_KEY=your_turnstile_secret_key
```

### 2) Install and run the backend

```bash
cd booking-backend
go mod download
go run ./cmd
```

The API runs on `http://localhost:8080` by default.

### 3) Backend health check

Open:

```text
GET http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```

## Frontend Setup

### 1) Create `booking-frontend/ganeshbooking/.env.local`

Use the following template:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_USER_API_URL=http://localhost:8080/user
NEXT_PUBLIC_ADMIN_API_URL=http://localhost:8080/admin

NEXT_PUBLIC_CAPTCHA_SITE_KEY=temporary_dev_captcha_site_key
NEXT_PUBLIC_CAPTCHA_ENABLED=false
```

### 2) Install and run the frontend

```bash
cd booking-frontend/ganeshbooking
npm install
npm run dev
```

The frontend runs on `http://localhost:3000`.

## Mail Verification

The backend already includes:

- `GET /user/verify?token=...` to verify a user account.
- `scripts.EmailVerifymail(to, token)` to send the verification email.
- `FRONTEND_URL` to build the redirect target after verification.

### How the flow works

1. Generate a verification token in your onboarding flow.
2. Send the token through `scripts.EmailVerifymail(to, token)`.
3. The email contains a link like:

```text
http://localhost:3000/auth/verified?query=<token>
```

4. That page can forward the token to the backend verification endpoint.
5. The backend marks the user as verified and redirects back to the frontend.

### Important note

The helper and endpoint are already present, but if you want automatic email verification during signup, make sure your signup flow actually sends the email after creating the token.

## CAPTCHA Setup

The auth routes support optional Cloudflare Turnstile verification.

### Disable CAPTCHA for local development

Set:

```env
CAPTCHA_ENABLED=false
NEXT_PUBLIC_CAPTCHA_ENABLED=false
```

### Enable CAPTCHA in production

Set:

```env
CAPTCHA_ENABLED=true
CAPTCHA_SECRET_KEY=your_turnstile_secret_key
NEXT_PUBLIC_CAPTCHA_ENABLED=true
NEXT_PUBLIC_CAPTCHA_SITE_KEY=your_turnstile_site_key
```

### Behavior

- When CAPTCHA is enabled, `/user/signin` and `/user/signup` require a valid Turnstile token.
- When CAPTCHA is disabled, those endpoints skip verification.

## Booking Flow

1. User signs in or signs up on the frontend.
2. User opens the booking page.
3. The booking form sends data to `POST /user/book`.
4. The backend stores the booking and generates a QR ticket.
5. The booking page displays the QR for the UTR value.

## Admin Flow

1. Admin signs in with `Adminname` and `Adminpass`.
2. The admin dashboard is protected by the admin token.
3. The dashboard shows bookings, verification status, and attendance.
4. The scanner page at `/scanner` is used for gate check-in.
5. If a booking is already marked present, the backend returns `Already entered` and does not update it again.

## API Endpoints

### User

- `GET /user/`
- `POST /user/signin`
- `POST /user/signup`
- `GET /user/verify?token=...`
- `GET /user/verifytoken`
- `GET /user/mybooking`
- `POST /user/book`

### Admin

- `GET /admin/`
- `POST /admin/adminauth`
- `POST /admin/verifyadmintoken`
- `POST /admin/verifybook`
- `POST /admin/enter`
- `PUT /admin/statuscount/:utr`

### Health

- `GET /health`

## Environment Variables Reference

### Backend variables

| Variable | Purpose | Required |
| --- | --- | --- |
| `PORT` | Backend listen port | No |
| `DB_URL` | PostgreSQL connection string | Yes, or `DATABASE_URL` |
| `DATABASE_URL` | Alternate PostgreSQL connection string | Yes, or `DB_URL` |
| `FRONTEND_URL` | Frontend base URL for redirects and email links | Recommended |
| `JWT_SECRET` | Admin token signing secret | Yes |
| `SESSION_SECRET` | Session or app secret used by your broader app setup | Recommended |
| `Adminname` | Admin username | Yes |
| `Adminpass` | Admin password | Yes |
| `Email` | Verified sender address | Yes for email sending |
| `Emailuser` | SMTP login | Yes for email sending |
| `Emailpass` | SMTP password | Yes for email sending |
| `CAPTCHA_ENABLED` | Turns CAPTCHA verification on or off | No |
| `CAPTCHA_SECRET_KEY` | Cloudflare Turnstile secret key | Yes if CAPTCHA is enabled |

### Frontend variables

| Variable | Purpose |
| --- | --- |
| `NEXT_PUBLIC_API_URL` | Base backend URL |
| `NEXT_PUBLIC_USER_API_URL` | User API base URL |
| `NEXT_PUBLIC_ADMIN_API_URL` | Admin API base URL |
| `NEXT_PUBLIC_CAPTCHA_SITE_KEY` | Turnstile site key |
| `NEXT_PUBLIC_CAPTCHA_ENABLED` | Enables CAPTCHA UI on the frontend |

## Production Notes

- Add your production frontend URL to the backend CORS allowlist in `cmd/main.go`.
- Set `FRONTEND_URL` to the deployed frontend domain.
- Use a strong `JWT_SECRET` and `SESSION_SECRET`.
- Enable CAPTCHA in production if you want extra auth protection.
- Use a real SMTP provider and verified sender domain for email delivery.

## Troubleshooting

### Backend will not start

- Check that `DB_URL` or `DATABASE_URL` is set.
- Check that PostgreSQL is running.
- Check the `.env` file is in `booking-backend/`.

### CAPTCHA errors

- Make sure `CAPTCHA_ENABLED=true` only when both frontend and backend keys are set.
- Keep CAPTCHA disabled for local development unless you have valid keys.

### Email verification links do not work

- Verify `FRONTEND_URL` points to the correct frontend domain.
- Check the SMTP credentials in `Emailuser` and `Emailpass`.
- Confirm the sender email in `Email` is allowed by your SMTP provider.

### Scanner does not open the camera

- Allow browser camera permission.
- Use HTTPS in production.
- Try the explicit `Allow camera access` button on the scanner page.

## Notes On Current Behavior

- User-side QR scanning is intentionally disabled.
- QR scanning is admin-only.
- The admin entry endpoint is idempotent, so scanning the same QR again returns `Already entered`.
