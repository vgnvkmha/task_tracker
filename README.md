# Task Tracker

Task Tracker is a Go backend service for managing users and teams. It exposes a JSON API and a small server-rendered HTML UI for user registration, login, and profile cabinet management.

The service uses:

- Go + Gin for HTTP routing
- PostgreSQL via `database/sql` and `pgx`
- JWT access tokens for authentication
- Optional temporary legacy header authentication for migration
- HTMX-backed HTML pages for the user UI

## Requirements

- Go 1.25+
- PostgreSQL
- TLS certificate files in the project root:
  - `cert.pem`
  - `key.pem`

The server runs with TLS on:

```text
https://localhost:8080
```

## Configuration

The application loads environment variables from `.env` using `godotenv`.

Example `.env`:

```env
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable

# Auth
JWT_SECRET=dev-secret-with-at-least-thirty-two-bytes
JWT_ACCESS_TTL_MINUTES=15
AUTH_LEGACY_HEADERS_ENABLED=false
```

Available variables:

| Variable | Default | Description |
| --- | --- | --- |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | empty | PostgreSQL user |
| `DB_PASSWORD` | empty | PostgreSQL password |
| `DB_NAME` | `postgres` | PostgreSQL database name |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | dev-only fallback outside production | Secret used to sign access tokens. Must be at least 32 bytes. Required in production. |
| `JWT_ACCESS_TTL_MINUTES` | `15` | Access token lifetime in minutes |
| `AUTH_LEGACY_HEADERS_ENABLED` | enabled outside production | Enables temporary fallback to legacy `X-User-ID` / `X-User-Role` auth |
| `APP_ENV` / `ENV` | `development` | If set to `production`, legacy auth defaults to disabled and `JWT_SECRET` is required |
| `GIN_MODE` | empty | If set to `release`, production-like auth defaults are used |

## Running The Service

Start the API:

```bash
go run ./cmd/api
```

If you use a self-signed local certificate, browsers and Postman may warn about TLS. For local testing, accept the browser warning or disable SSL certificate verification in Postman.

## Authentication

### JWT Access Tokens

The primary authentication method is:

```text
Authorization: Bearer <access_token>
```

Access tokens are issued by:

```text
POST /auth/login
```

Access tokens contain:

- `sub`: user ID
- `role`: user role
- `team_id`: user team ID, when present
- `iss`: token issuer
- `aud`: token audience
- `iat`: issued-at timestamp
- `exp`: expiration timestamp
- `jti`: token ID

### Browser UI Auth

The HTML UI uses the same access token, but stores it in an `HttpOnly`, `Secure`, `SameSite=Lax` cookie named:

```text
access_token
```

After a successful UI login, the browser receives this cookie and can open the user cabinet until the JWT expires.

If the UI token is missing or invalid, the user is redirected to:

```text
/ui/users/create?auth=required
```

If the UI token is expired, the user is redirected to:

```text
/ui/users/create?auth=expired
```

### Temporary Legacy Header Auth

During migration, the API can fall back to legacy headers when `AUTH_LEGACY_HEADERS_ENABLED=true`:

```text
X-User-ID: <uuid>
X-User-Role: admin|captain|user|guest
```

JWT has priority. If an `Authorization` header is present but invalid, the request is rejected and does not fall back to legacy headers.

Each legacy fallback is logged as a warning.

## Roles

Supported roles:

- `admin`
- `captain`
- `user`
- `guest`

Manager roles are:

- `admin`
- `captain`

Some user operations are restricted to manager roles or to the user updating/deleting their own account.

## UI Endpoints

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| `GET` | `/ui/users/create` | Public | Login and registration page |
| `POST` | `/ui/users` | Public | Register a user from the HTML form |
| `GET` | `/ui/users/login` | Public | Login from the HTML form. On success, sets `access_token` cookie and redirects to cabinet |
| `GET` | `/ui/users/success` | Public | Simple success page |
| `GET` | `/ui/users/cabinet` | JWT cookie | User cabinet for the current authenticated user |
| `POST` | `/ui/users/cabinet/update` | JWT cookie | Update current user's profile from the cabinet |
| `POST` | `/ui/users/cabinet/delete` | JWT cookie | Delete current user's account from the cabinet |

UI login flow:

1. Open `https://localhost:8080/ui/users/create`.
2. Log in with email and password.
3. The server sets the `access_token` cookie.
4. The browser redirects to `/ui/users/cabinet?...`.
5. While the JWT is valid, `/ui/users/cabinet` can be opened without logging in again.
6. After the JWT expires, the browser is redirected to `/ui/users/create?auth=expired`.

## API Endpoints

All `/user/*` and `/team/*` API endpoints require JWT auth or legacy header auth when the legacy flag is enabled.

### Auth

#### `POST /auth/login`

Request:

```json
{
  "email": "user@example.com",
  "password": "Strong1!"
}
```

Response:

```json
{
  "access_token": "<jwt>",
  "token_type": "Bearer",
  "expires_in": 900
}
```

Use the token in API requests:

```text
Authorization: Bearer <jwt>
```

### Users

#### `POST /user/create_register`

Creates a user through the API.

Request:

```json
{
  "email": "user@example.com",
  "password": "Strong1!",
  "role": "user",
  "team_name": "Backend",
  "first_name": "Ivan",
  "last_name": "Petrov",
  "age": 30,
  "birth_date": "1996-05-20T00:00:00Z"
}
```

#### `POST /user/create_by_actor`

Creates a user as the authenticated actor. Manager role is required.

Request body is the same as `/user/create_register`.

#### `PATCH /user/update`

Updates a user.

Non-manager users can update only themselves and cannot change `role` or `team_id`.

Request:

```json
{
  "user_id": "<uuid>",
  "email": "new@example.com",
  "password": "Strong1!",
  "role": "user",
  "team_id": "<uuid>",
  "first_name": "Ivan",
  "last_name": "Petrov",
  "age": 31,
  "birth_date": "1995-05-20T00:00:00Z"
}
```

All fields except `user_id` are optional.

#### `GET /user/get_by_id/:user_id`

Returns one user by ID.

#### `GET /user/list_active`

Returns active users.

#### `GET /user/list`

Returns users.

#### `DELETE /user/delete/:id`

Deletes a user. Manager users can delete any user. Non-manager users can delete only themselves.

### Teams

#### `POST /team/create`

Creates a team.

Request:

```json
{
  "name": "Backend",
  "timezone": "Europe/Moscow",
  "leader_id": "<uuid>"
}
```

`timezone` and `leader_id` are optional.

#### `GET /team/get_by_id/:id`

Returns a team by ID.

#### `GET /team/get_by_name/:team_name`

Returns a team by name.

#### `GET /team/list_active`

Returns active teams.

#### `GET /team/list`

Returns teams.

#### `PATCH /team/update/:id`

Updates a team.

Request:

```json
{
  "name": "Platform",
  "timezone": "Europe/Moscow",
  "leader_id": "<uuid>"
}
```

All fields are optional.

#### `DELETE /team/delete_by_id/:id`

Deletes a team by ID.

## Testing With curl

Log in:

```bash
curl -k -X POST https://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"Strong1!"}'
```

Call a protected endpoint:

```bash
curl -k https://localhost:8080/user/list \
  -H "Authorization: Bearer <access_token>"
```

Check missing-token behavior:

```bash
curl -k https://localhost:8080/user/list
```

Expected response when legacy auth is disabled:

```json
{
  "error": "missing bearer token"
}
```

Check expired-token behavior by setting a short TTL:

```env
JWT_ACCESS_TTL_MINUTES=1
```

Then log in, wait more than one minute, and call a protected endpoint again.

## Testing With Postman

1. Disable SSL verification for local self-signed certificates:
   `Settings -> General -> SSL certificate verification -> OFF`.
2. Send `POST https://localhost:8080/auth/login`.
3. Copy `access_token` from the response.
4. For protected requests, set:
   - Authorization type: `Bearer Token`
   - Token: the copied access token
5. Call endpoints such as:
   - `GET https://localhost:8080/user/list`
   - `GET https://localhost:8080/team/list`

## Running Tests

Run all tests:

```bash
go test ./...
```

Current known issue: the full suite may fail in `internal/application/team` on `TestServiceUpdateReturnsInvalidInputForInvalidTimezone`. This is unrelated to the JWT authentication layer.

Focused auth/UI tests:

```bash
go test ./internal/infrastracture/auth ./internal/transport/http/middleware ./internal/transport/http/user ./internal/transport/http/auth ./internal/app
```

## Security Notes

- Do not commit real `.env` secrets.
- Use a strong `JWT_SECRET` in staging and production.
- Keep `AUTH_LEGACY_HEADERS_ENABLED=false` in production.
- Access tokens are intentionally short-lived.
- Refresh tokens are not implemented yet.
- Browser UI uses `HttpOnly` cookies instead of localStorage.
- The current app still registers `MockActorMiddleware` globally; verify this before production hardening.
