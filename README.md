# Attendly

Attendly is an attendance and billing platform for recurring sessions that share booked resources, such as sports courts, meeting rooms, coworking desks, studios or lab equipment. It records which members participated in each time slot, calculates individual charges from actual participation, and maintains a running ledger of billed amounts, payments and outstanding balances.

Bookings typically span several hours, while members arrive and leave at different times. Rather than splitting the total cost evenly across everyone who attended, Attendly bills each member only for the time slots they participated in, at their proportional share of the units booked for those slots. Unpaid balances carry forward automatically to subsequent sessions.

Built with **Go (Gin, GORM, PostgreSQL)** and **React (TypeScript, Vite)**.

## Features

- **Session and time slot management.** Sessions can be divided into hourly slots automatically, with the number of units configured per slot.
- **Slot-level attendance.** An attendance grid (members × time slots) records participation for each slot.
- **Proportional billing.** Each slot's booking cost is divided evenly among the members in that slot.
- **Payments and ledger.** Payments are recorded per member, with balances carried forward across sessions.
- **Member history.** Each member has a per-session breakdown of charges and payments, with a WhatsApp link for payment reminders.
- **Authentication.** JWT-based authentication with token expiry and configurable self-registration.
- **API documentation.** OpenAPI specification served through Swagger UI at `/docs`.

## Billing model

A *unit* is whatever is booked: a court, a room, a desk or a machine. For every time slot:

```
slot cost       = units booked × unit price per hour
cost per member = slot cost ÷ members in that slot
```

A member's charge for a session is the sum of their charges across all slots they participated in. Slots with no members are not billed.

**Example** (unit price ₹300 per hour):

| Time slot   | Units booked  | Members in slot | Charge per member |
| ----------- | ------------- | --------------- | ----------------- |
| 7pm – 8pm   | 2             | 3               | ₹200              |
| 8pm – 9pm   | 1             | 2               | ₹150              |

A member present in both slots is billed ₹350; a member present only in the first slot is billed ₹200.

Payments are recorded against a member and session rather than individual time slots. Each member's balance is calculated as:

```
balance = total paid − total billed      (negative = still owes)
```

A member billed ₹350 who pays ₹300 carries a balance of −₹50 into subsequent sessions.

## Project structure

```
attendly/
├── cmd/api/main.go            # entry point and routes
├── internal/
│   ├── config/                # env loading and database connection
│   ├── domain/                # models and request types
│   ├── services/              # business logic: billing, ledger, sessions, auth (+ tests)
│   ├── handlers/              # HTTP layer only
│   ├── middleware/            # JWT authentication
│   └── docs/                  # OpenAPI spec and Swagger UI
└── frontend/                  # React app
    └── src/
        ├── api/               # API client and typed endpoint functions
        ├── pages/             # Dashboard, Sessions, Session detail, Members, ...
        ├── components/
        └── lib/               # formatting and helpers
```

Handlers deal with HTTP; anything that decides money or business rules lives in `services`, where it can be unit tested without a database.

## Running locally

**Requirements:** Go 1.26+, Node 22+, PostgreSQL.

### 1. Database

```bash
createdb attendly
```

Tables are created automatically when the API starts.

### 2. API

```bash
cp .env.example .env      # then set DB_PASSWORD and JWT_SECRET
go run ./cmd/api
```

The API listens on `http://localhost:8080`. Check it with `curl localhost:8080/health`.

### 3. Frontend

```bash
cd frontend
npm install
npm run dev
```

Open `http://localhost:5173` and create an account. In development, Vite forwards `/api/*` to the Go server, so no CORS setup is needed.

### Tests

```bash
go test ./...
cd frontend && npm run lint && npm run build
```

## Configuration

All settings come from environment variables, or from a `.env` file in the directory you start the API from. See [`.env.example`](.env.example).

| Variable             | Default                  | Notes                                             |
| -------------------- | ------------------------ | ------------------------------------------------- |
| `DATABASE_URL`       | (empty)                  | Full connection string; overrides the `DB_*` values |
| `DB_HOST`            | `localhost`              |                                                   |
| `DB_PORT`            | `5432`                   |                                                   |
| `DB_USER`            | `postgres`               |                                                   |
| `DB_PASSWORD`        | (empty)                  |                                                   |
| `DB_NAME`            | `attendly`             |                                                   |
| `DB_SSLMODE`         | `disable`                | Use `require` for hosted databases                |
| `SERVER_PORT`        | `8080`                   | Falls back to `PORT` if set                       |
| `JWT_SECRET`         | (required)               | The API refuses to start without it               |
| `JWT_TTL_HOURS`      | `168`                    | How long a login lasts (7 days)                   |
| `CORS_ORIGINS`       | `http://localhost:5173`  | Comma-separated                                   |
| `ALLOW_REGISTRATION` | `true`                   | Set to `false` once your own account exists       |

## API

Full interactive docs are served at **`http://localhost:8080/docs`**.
Everything except `/health`, `/register`, `/login` and `/docs` needs an `Authorization: Bearer <token>` header from `/login`.

| Method | Path                                  | Description                                    |
| ------ | ------------------------------------- | ---------------------------------------------- |
| POST   | `/register`                           | Create an account                              |
| POST   | `/login`                              | Get a JWT                                      |
| GET    | `/sessions`                           | List sessions with slots and members           |
| POST   | `/sessions`                           | Create a session (optionally with hourly slots)|
| GET    | `/sessions/:id`                       | One session                                    |
| DELETE | `/sessions/:id`                       | Delete a session (refused if it has payments)  |
| GET    | `/sessions/:id/billing`               | What each member owes for the session          |
| POST   | `/sessions/:id/timeslots`             | Add a time slot                                |
| PUT    | `/timeslots/:id`                      | Change units booked                           |
| DELETE | `/timeslots/:id`                      | Delete a slot                                  |
| POST   | `/timeslots/:id/members`              | Mark a member as present in a slot             |
| DELETE | `/timeslots/:id/members/:memberId`    | Remove a member from a slot                    |
| GET    | `/members`                            | List members                                   |
| POST   | `/members`                            | Add a member                                   |
| GET    | `/members/:id/ledger`                 | A member's balance, sessions and payments      |
| GET    | `/ledger`                             | Every member's balance                         |
| GET    | `/payments`                           | List payments (`?member_id=`, `?session_id=`)  |
| POST   | `/payments`                           | Record a payment                               |
| DELETE | `/payments/:id`                       | Delete a payment                               |

## Deployment

The frontend is published to GitHub Pages and the API runs on Render, backed by a Neon Postgres database.

**Database (Neon).** Create a project and copy its connection string (`postgres://…?sslmode=require`).

**API (Render).** In the Render dashboard choose *New → Blueprint* and select this repository. [`render.yaml`](render.yaml) defines the service; when prompted, paste the Neon connection string as `DATABASE_URL`. `JWT_SECRET` is generated automatically. Once deployed, note the service URL, e.g. `https://attendly-api.onrender.com`.

**Frontend (GitHub Pages).**
1. In the repository settings, set *Pages → Source* to **GitHub Actions**.
2. Add a repository variable `VITE_API_URL` (*Settings → Secrets and variables → Actions → Variables*) containing the Render URL.
3. Push to `master` or run the *Deploy frontend* workflow manually. The app is served at `https://<user>.github.io/attendly/`.

After creating your account, set `ALLOW_REGISTRATION=false` on Render so no one else can sign up.

Render's free tier sleeps after a period of inactivity, so the first request after a pause can take up to a minute.
