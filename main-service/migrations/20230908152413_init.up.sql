CREATE TABLE "users" (
                        id SERIAL PRIMARY KEY,
                        username VARCHAR(50) NOT NULL,
                        email VARCHAR(100) UNIQUE NOT NULL,
                        password_hash VARCHAR(255),
                        role VARCHAR(50) NOT NULL DEFAULT 'doctor',
                        status VARCHAR(50) NOT NULL DEFAULT 'invited',
                        invite_token VARCHAR(255),
                        invite_token_expires_at TIMESTAMPTZ
);

INSERT INTO users (username, email, password_hash, role, status, invite_token)
VALUES (
    'Admin',
    'admin@example.com',
    '$2a$14$qE3FhvHU5w.lUQkOOOMe8urygdTAwuQIXqj6JjBzkz5AOnOJxpaMe',
    'admin',
    'active',
    NULL
);

CREATE TABLE mri_scans (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    patient_name VARCHAR(100) NOT NULL,
    patient_gender VARCHAR(10),
    patient_age INT,
    scan_date DATE,
    object_name VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255),
    content_type VARCHAR(100),
    size BIGINT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    status VARCHAR(50) DEFAULT 'queued'
);

CREATE TABLE mri_analysis_results (
    id SERIAL PRIMARY KEY,
    scan_id INTEGER NOT NULL REFERENCES mri_scans(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    status VARCHAR(50) DEFAULT 'pending',  -- pending, processing, done, failed
    diagnosis INTEGER CHECK (diagnosis IN (0, 1, 2)), -- 0: Healthy, 1: Mild Cognitive Impairment, 2: Alzheimer's
    confidence REAL CHECK (confidence >= 0 AND confidence <= 1),
    gradcam_url TEXT,
    error_message TEXT
);
