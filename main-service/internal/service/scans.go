package service

import (
	"context"
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

func (s *Service) GetScansByFilters(
	ctx context.Context,
	userID int64,
	filter *domain.ScanFilter,
) ([]*domain.MRIScan, error) {

	return s.repo.GetScansByFilters(ctx, userID, filter)
}

func (s *Service) GetScanByID(
	ctx context.Context,
	userID, scanID int64,
) (*domain.MRIScanDetail, error) {

	scan, err := s.repo.GetScanByID(ctx, userID, scanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan detail: %w", err)
	}

	// Преобразуем Grad-CAM URL в presigned ссылку (если имеется)
	if scan.GradCAMURL != nil && *scan.GradCAMURL != "" {
		presigned, err := s.repo.PresignedGetObject(*scan.GradCAMURL)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned url: %w", err)
		}
		url := presigned.String()
		scan.GradCAMURL = &url
	}
	return scan, nil
}
