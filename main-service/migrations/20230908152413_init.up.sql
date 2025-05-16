
BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id                   SERIAL PRIMARY KEY,
    username             VARCHAR(50)  NOT NULL,
    email                VARCHAR(100) UNIQUE NOT NULL,
    password_hash        VARCHAR(255),
    role                 VARCHAR(50)  NOT NULL DEFAULT 'doctor',
    status               VARCHAR(50)  NOT NULL DEFAULT 'invited',
    invite_token         VARCHAR(255),
    invite_token_expires_at TIMESTAMPTZ
);


INSERT INTO users (username, email, password_hash, role, status)
VALUES ('Admin',
        'admin@example.com',
        '$2a$14$qE3FhvHU5w.lUQkOOOMe8urygdTAwuQIXqj6JjBzkz5AOnOJxpaMe',
        'admin',
        'active')
ON CONFLICT (email) DO NOTHING;

CREATE TABLE IF NOT EXISTS mri_scans (
    id                SERIAL PRIMARY KEY,
    user_id           INTEGER      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    patient_name      VARCHAR(100) NOT NULL,
    patient_gender    VARCHAR(10),
    patient_age       INT,
    scan_date         DATE,
    object_name       VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255),
    content_type      VARCHAR(100),
    size              BIGINT,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    status            VARCHAR(50)  NOT NULL DEFAULT 'queued'
);

CREATE TABLE IF NOT EXISTS mri_analysis_results (
    id          SERIAL PRIMARY KEY,
    scan_id     INTEGER      NOT NULL REFERENCES mri_scans(id) ON DELETE CASCADE,
    started_at  TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    status      VARCHAR(50)  NOT NULL DEFAULT 'pending',
    diagnosis   INTEGER      CHECK (diagnosis IN (0,1,2)),
    confidence  REAL         CHECK (confidence BETWEEN 0 AND 1),
    gradcam_url TEXT,
    error_message TEXT
);

-- mri_scans ------------------------------------------------
ALTER TABLE mri_scans ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS mri_scans_doctor_select ON mri_scans;
CREATE POLICY mri_scans_doctor_select
    ON mri_scans
    FOR SELECT
    USING (user_id = current_setting('app.user_id')::int);

DROP POLICY IF EXISTS mri_scans_doctor_insert ON mri_scans;
CREATE POLICY mri_scans_doctor_insert
    ON mri_scans
    FOR INSERT
    WITH CHECK (TRUE); 

-- mri_analysis_results ------------------------------------
ALTER TABLE mri_analysis_results ENABLE ROW LEVEL SECURITY;

DROP POLICY IF EXISTS mri_analysis_results_doctor_select ON mri_analysis_results;
CREATE POLICY mri_analysis_results_doctor_select
    ON mri_analysis_results
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1
            FROM mri_scans s
            WHERE s.id = mri_analysis_results.scan_id
              AND s.user_id = current_setting('app.user_id')::int
        )
    );

COMMIT;
