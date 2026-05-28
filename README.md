# DeckForge — AI Pitch Deck Generator

DeckForge is a full-stack SaaS app that turns PDFs, Word documents, or pasted text into professional 7-slide pitch decks using OpenAI GPT.

## Architecture Overview

```
┌─────────────┐     REST + JWT      ┌─────────────┐     SQL      ┌──────────────┐
│   React     │ ◄─────────────────► │  Go (Gin)   │ ◄──────────► │  PostgreSQL  │
│  Frontend   │                     │   Backend   │              │              │
└─────────────┘                     └──────┬──────┘              └──────────────┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │  OpenAI API │
                                    └─────────────┘
```

### Clean Architecture Layers (Backend)

| Layer          | Folder                | Responsibility                     |
| -------------- | --------------------- | ---------------------------------- |
| **Entry**      | `cmd/server`          | Starts the app, wires dependencies |
| **Router**     | `internal/router`     | Maps URLs → handlers               |
| **Handlers**   | `internal/handler`    | HTTP: parse request, return JSON   |
| **Services**   | `internal/service`    | Business logic: auth, AI, files    |
| **Models**     | `internal/models`     | Database table structs             |
| **Middleware** | `internal/middleware` | JWT validation on protected routes |

Handlers should stay thin; services contain the real logic.

### Frontend Structure

| Folder           | Purpose                                        |
| ---------------- | ---------------------------------------------- |
| `src/pages`      | Full page components (Dashboard, Upload, etc.) |
| `src/components` | Reusable UI (Button, DropZone, SlideCard)      |
| `src/api`        | Axios calls to backend                         |
| `src/context`    | Global auth state (React Context)              |
| `src/hooks`      | Reusable logic (protected route redirect)      |

---

## Data Flow (Step by Step)

1. **Register/Login** — User submits email/password → Go hashes password with bcrypt → returns JWT → React stores token in `localStorage`.
2. **Upload** — User drops a PDF/DOCX/TXT → `POST /api/v1/upload` → file saved to disk → text extracted → row in `uploaded_files`.
3. **Generate** — `POST /api/v1/generate` with `file_id` or `text` → OpenAI returns JSON slides → saved to `presentations` + `slides` tables.
4. **Preview** — `GET /api/v1/presentation/:id` → frontend renders `SlideCard` components.
5. **Export** — `GET /api/v1/presentation/:id/export/pptx|pdf` → file built on disk → downloaded in browser.

---

## API Endpoints

| Method | Path                                   | Auth | Description                         |
| ------ | -------------------------------------- | ---- | ----------------------------------- |
| POST   | `/api/v1/register`                     | No   | Create account                      |
| POST   | `/api/v1/login`                        | No   | Get JWT token                       |
| POST   | `/api/v1/upload`                       | Yes  | Upload PDF/DOCX/TXT                 |
| POST   | `/api/v1/generate`                     | Yes  | Generate deck (`file_id` or `text`) |
| GET    | `/api/v1/presentations`                | Yes  | List user's decks                   |
| GET    | `/api/v1/presentation/:id`             | Yes  | Get deck + slides                   |
| GET    | `/api/v1/presentation/:id/export/pptx` | Yes  | Download PowerPoint                 |
| GET    | `/api/v1/presentation/:id/export/pdf`  | Yes  | Download PDF                        |
| GET    | `/health`                              | No   | Health check                        |

**Auth header:** `Authorization: Bearer <your-jwt-token>`

---

## Quick Start (Local Development)

### Prerequisites

- Go 1.22+
- Node.js 20+
- PostgreSQL 16 (or Docker)
- OpenAI API key

### 1. Clone and configure

```bash
cd DeckForger
cp .env.example .env
# Edit .env — set OPENAI_API_KEY and JWT_SECRET
```

### 2. Start PostgreSQL

```bash
docker compose up db -d
# Or use your own Postgres and update DATABASE_URL
```

### 3. Run the backend

```bash
cd backend
go mod download
go run ./cmd/server
# API: http://localhost:8080
```

### 4. Run the frontend

```bash
cd frontend
npm install
npm run dev
# App: http://localhost:5173
```

The Vite dev server proxies `/api` to `localhost:8080`.

---

## Docker (Full Stack)

```bash
cp .env.example .env
# Set OPENAI_API_KEY in .env

docker compose up --build
```

- Frontend: http://localhost:5173
- Backend: http://localhost:8080
- Postgres: localhost:5433

---

## Environment Variables

See [`.env.example`](.env.example) for all variables.

| Variable         | Description                         |
| ---------------- | ----------------------------------- |
| `DATABASE_URL`   | PostgreSQL connection string        |
| `JWT_SECRET`     | Secret for signing JWT tokens       |
| `OPENAI_API_KEY` | Your OpenAI API key                 |
| `OPENAI_MODEL`   | Model name (default: `gpt-4o-mini`) |
| `VITE_API_URL`   | Frontend → backend URL              |

---

## Database Schema

See [`database/schema.sql`](database/schema.sql).

- **users** — accounts
- **uploaded_files** — source documents + extracted text
- **presentations** — generated decks
- **slides** — individual slides (JSON bullets in `content`)

---

## Slide Types Generated

1. `title` — Company name & tagline
2. `problem` — Pain points
3. `solution` — Your product
4. `market` — Opportunity
5. `features` — Key capabilities
6. `roadmap` — Timeline
7. `conclusion` — Call to action

---

## Project Structure

```
DeckForger/
├── backend/
│   ├── cmd/server/main.go          # Entry point
│   ├── internal/
│   │   ├── config/                 # Env loading
│   │   ├── database/               # Postgres connection
│   │   ├── handler/                # HTTP handlers
│   │   ├── middleware/             # JWT auth
│   │   ├── models/                 # DB models
│   │   ├── router/                 # Route definitions
│   │   └── service/                # Business logic
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── api/                    # Axios API client
│   │   ├── components/             # UI components
│   │   ├── context/                # Auth provider
│   │   ├── hooks/
│   │   ├── pages/
│   │   └── types/
│   ├── Dockerfile
│   └── package.json
├── database/schema.sql
├── docker-compose.yml
└── .env.example
```

---

## Key Concepts for Learners

### JWT Authentication

After login, the server returns a signed token. The frontend sends it on every request. Middleware verifies the signature before allowing access to `/upload`, `/generate`, etc.

### ORM (GORM)

Go structs in `models/` map to SQL tables. `db.Create()`, `db.Where().First()` replace hand-written SQL for common operations.

### React Context

`AuthProvider` wraps the app so any component can call `useAuth()` without passing props through every level ("prop drilling").

### Multipart Upload

Files are sent as `multipart/form-data`, not JSON. The backend uses `c.FormFile("file")` to read them.

---

## License

MIT — use freely for learning and projects.
