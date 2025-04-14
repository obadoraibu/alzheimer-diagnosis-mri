package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

func (s *Service) CreateMRIAnalysis(ctx *gin.Context, input *domain.CreateAnalysisInput) error {
	objectName := fmt.Sprintf("%d/%d_%s", input.UserID, time.Now().Unix(), input.Header.Filename)
	contentType := input.Header.Header.Get("Content-Type")

	err := s.repo.UploadScanToMinIO(ctx, objectName, input.File, input.Header.Size, contentType)
	if err != nil {
		return fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	err = s.repo.SaveScanMetadata(input.UserID, objectName, input.Header.Filename, contentType, input.Header.Size, input.PatientName, input.PatientGender, input.PatientAge, input.ScanDate)
	if err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	err = s.repo.EnqueueScanTask(input.UserID, objectName)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}
