# DeckForge Architecture Guide (For Learners)

## What is DeckForge?

DeckForge is an AI SaaS app. Users upload business documents, and the app creates a 7-slide investor pitch deck using OpenAI.

---

## High-Level Diagram

```
Browser (React)
    │
    │  HTTP + JSON
    │  Authorization: Bearer <JWT>
    ▼
Go API (Gin)
    ├── AuthService      → users table, bcrypt, JWT
    ├── ExtractService   → PDF / DOCX / TXT → plain text
    ├── OpenAIService    → text → structured slides JSON
    ├── PresentationService → save to DB
    └── ExportService    → PPTX (ZIP/XML) + PDF (gofpdf)
    │
    ├── PostgreSQL (users, files, presentations, slides)
    └── OpenAI API (GPT)
```

---

## Request Lifecycle Example: Generate Deck

1. **User** drops `business-plan.pdf` on Upload page.
2. **Frontend** `uploadFile()` sends `multipart/form-data` to `POST /api/v1/upload`.
3. **Auth middleware** reads JWT, sets `userID` on Gin context.
4. **PresentationHandler.Upload** reads file bytes.
5. **PresentationService.SaveUpload** writes file to `./uploads`, calls **ExtractService**.
6. **ExtractService** uses `go-fitz` for PDF, custom XML parser for DOCX.
7. Row inserted in `uploaded_files` with `extracted_text`.
8. User clicks **Generate** → `POST /api/v1/generate` with `{ "file_id": "..." }`.
9. **OpenAIService.GenerateSlides** sends text to GPT with JSON response format.
10. GPT returns `{ title, slides: [...] }` with 7 slide types.
11. **PresentationService** inserts `presentations` + 7 `slides` rows.
12. **Frontend** navigates to `/presentation/:id` and renders `SlideCard` components.

---

## Why Clean Architecture?

| Layer | Knows about | Does NOT know about |
|-------|-------------|---------------------|
| Handler | HTTP, JSON status codes | OpenAI prompts |
| Service | Business rules, DB, AI | HTTP headers |
| Model | Table columns | API routes |

This makes code easier to test and change. Example: swap OpenAI for Claude by only editing `openai_service.go`.

---

## JWT Flow

1. Login returns `{ token, user }`.
2. Frontend stores token in `localStorage`.
3. Axios interceptor adds `Authorization: Bearer ...` to every request.
4. `AuthMiddleware` validates signature with `JWT_SECRET`.
5. If invalid → 401 → frontend redirects to `/login`.

---

## Database Relationships

```
users (1) ──< uploaded_files
users (1) ──< presentations
presentations (1) ──< slides
uploaded_files (1) ──< presentations (optional FK)
```

---

## File Reference

### Backend

| File | Role |
|------|------|
| `cmd/server/main.go` | Wires everything and starts HTTP server |
| `internal/router/router.go` | URL → handler mapping |
| `internal/middleware/auth.go` | JWT gate for protected routes |
| `internal/handler/auth_handler.go` | Register/login HTTP |
| `internal/handler/presentation_handler.go` | Upload, generate, export HTTP |
| `internal/service/auth_service.go` | Password hashing, token creation |
| `internal/service/openai_service.go` | GPT integration |
| `internal/service/extract_service.go` | Document text extraction |
| `internal/service/presentation_service.go` | Core deck workflow |
| `internal/service/export_service.go` | PDF export |
| `internal/service/pptx_writer.go` | PPTX as Office Open XML ZIP |

### Frontend

| File | Role |
|------|------|
| `src/App.tsx` | Route definitions |
| `src/context/AuthContext.tsx` | Global login state |
| `src/api/client.ts` | Axios + JWT interceptor |
| `src/pages/Upload.tsx` | Drag-drop + generate UI |
| `src/pages/Presentation.tsx` | Slide preview + export buttons |
| `src/components/slides/SlideCard.tsx` | Visual slide card |

---

## Concepts Glossary

- **REST API** — Standard way for frontend to talk to backend using URLs and HTTP verbs.
- **JWT** — Signed token proving the user is logged in (stateless auth).
- **ORM (GORM)** — Maps Go structs to SQL tables.
- **Middleware** — Code that runs before your handler (like a security guard).
- **CORS** — Browser security rule; backend must allow frontend origin.
- **JSONB** — PostgreSQL column type for JSON arrays (slide bullets).
