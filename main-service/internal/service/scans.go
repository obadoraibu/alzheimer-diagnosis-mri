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

	scan_id, err := s.repo.SaveScanMetadata(input.UserID, objectName, input.Header.Filename, contentType, input.Header.Size, input.PatientName, input.PatientGender, input.PatientAge, input.ScanDate)
	if err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	err = s.repo.EnqueueScanTask(scan_id, objectName)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	return nil
}

func (s *Service) GetScansByFilters(userID int64, filter *domain.ScanFilter) ([]*domain.MRIScan, error) {
	return s.repo.GetScansByFilters(userID, filter)
}

func (s *Service) GetScanByID(userId, scanId int64) (*domain.MRIScanDetail, error) {
	scan, err := s.repo.GetScanDetail(userId, scanId)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan detail: %w", err)
	}

	if scan.GradCAMURL != nil && *scan.GradCAMURL != "" {
		presignedURL, err := s.repo.PresignedGetObject(*scan.GradCAMURL)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned url: %w", err)
		}
		urlStr := presignedURL.String()
		//urlStr = strings.Replace(urlStr, "http://minio:9000", "http://localhost:9000", 1)
		scan.GradCAMURL = &urlStr
	}

	return scan, nil
}
