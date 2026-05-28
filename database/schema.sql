-- =============================================================================
-- DeckForge Database Schema
-- Tables: users, uploaded_files, presentations, slides
-- =============================================================================

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- -----------------------------------------------------------------------------
-- USERS: stores registered accounts
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- -----------------------------------------------------------------------------
-- UPLOADED_FILES: raw documents users upload (PDF, DOCX, text)
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS uploaded_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    original_name VARCHAR(500) NOT NULL,
    stored_path VARCHAR(1000) NOT NULL,
    file_type VARCHAR(50) NOT NULL, -- pdf, docx, txt
    file_size BIGINT NOT NULL DEFAULT 0,
    extracted_text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_uploaded_files_user_id ON uploaded_files(user_id);

-- -----------------------------------------------------------------------------
-- PRESENTATIONS: a generated pitch deck belongs to one user
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS presentations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    uploaded_file_id UUID REFERENCES uploaded_files(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    source_summary TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_presentations_user_id ON presentations(user_id);

-- -----------------------------------------------------------------------------
-- SLIDES: individual slides within a presentation
-- slide_type: title, problem, solution, market, features, roadmap, conclusion
-- -----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS slides (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    presentation_id UUID NOT NULL REFERENCES presentations(id) ON DELETE CASCADE,
    slide_order INT NOT NULL,
    slide_type VARCHAR(50) NOT NULL,
    title VARCHAR(500) NOT NULL,
    subtitle VARCHAR(500),
    content JSONB NOT NULL DEFAULT '[]', -- bullet points or extra fields
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_slides_presentation_id ON slides(presentation_id);
