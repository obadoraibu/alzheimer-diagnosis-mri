package domain

import (
	"mime/multipart"
	"time"
)

type CreateAnalysisInput struct {
	UserID        int64
	File          multipart.File
	Header        *multipart.FileHeader
	PatientName   string
	PatientGender string
	PatientAge    int
	ScanDate      time.Time
}

// [x]
type ScanFilter struct {
	ScanID       *int64
	UploadedFrom *time.Time
	UploadedTo   *time.Time
}

type MRIScan struct {
	ID            int64
	UserID        int64
	PatientName   string
	PatientGender string
	PatientAge    int
	ScanDate      time.Time
	CreatedAt     time.Time
	Status        string
}

type MRIScanDetail struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	PatientName   string    `json:"patient_name"`
	PatientGender string    `json:"patient_gender"`
	PatientAge    int       `json:"patient_age"`
	ScanDate      time.Time `json:"scan_date"`
	ObjectName    string    `json:"object_name"`
	OriginalName  string    `json:"original_filename"`
	ContentType   string    `json:"content_type"`
	Size          int64     `json:"size"`
	CreatedAt     time.Time `json:"created_at"`
	Status        string    `json:"status"`

	// Анализ
	Diagnosis   *int       `json:"diagnosis,omitempty"`
	Confidence  *float32   `json:"confidence,omitempty"`
	GradCAMURL  *string    `json:"gradcam_url,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
