package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
)

func (h *Handler) CreateMRIAnalysis(c *gin.Context) {
	token := c.MustGet("AccessToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	patientName := c.PostForm("patient_name")
	patientGender := c.PostForm("patient_gender")
	patientAgeStr := c.PostForm("patient_age")
	scanDateStr := c.PostForm("scan_date")

	if patientName == "" || patientGender == "" || patientAgeStr == "" || scanDateStr == "" {
		sendErrorResponse(c, http.StatusBadRequest, "MISSING_FIELDS", "Required form fields are missing")
		return
	}

	patientAge, err := strconv.Atoi(patientAgeStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_AGE", "Patient age must be a valid number")
		return
	}

	scanDate, err := time.Parse("2006-01-02", scanDateStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_SCAN_DATE", "Scan date must be in YYYY-MM-DD format")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "FILE_REQUIRED", "Scan file is required")
		return
	}
	defer file.Close()

	input := &domain.CreateAnalysisInput{
		UserID:        userID,
		File:          file,
		Header:        header,
		PatientName:   patientName,
		PatientGender: patientGender,
		PatientAge:    patientAge,
		ScanDate:      scanDate,
	}

	err = h.service.CreateMRIAnalysis(c, input)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			sendErrorResponse(c, httpStatusFromAppError(appErr), appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create analysis")
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
	})
}

func (h *Handler) ListScans(c *gin.Context) {
	token := c.MustGet("AccessToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	f := &domain.ScanFilter{}

	if idStr := c.Query("id"); idStr != "" {
		if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			f.ScanID = &id
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "INVALID_SCAN_ID", "Scan ID must be a valid integer")
			return
		}
	}

	if fromStr := c.Query("uploaded_from"); fromStr != "" {
		if from, err := time.Parse("2006-01-02", fromStr); err == nil {
			f.UploadedFrom = &from
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "INVALID_DATE_FROM", "Invalid uploaded_from date format")
			return
		}
	}

	if toStr := c.Query("uploaded_to"); toStr != "" {
		if to, err := time.Parse("2006-01-02", toStr); err == nil {
			f.UploadedTo = &to
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "INVALID_DATE_TO", "Invalid uploaded_to date format")
			return
		}
	}

	scans, err := h.service.GetScansByFilters(c, userID, f)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			sendErrorResponse(c, httpStatusFromAppError(appErr), appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve scans")
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    scans,
	})
}

func (h *Handler) GetScanDetail(c *gin.Context) {
	token := c.MustGet("AccessToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	scanIDParam := c.Param("id")
	scanID, err := strconv.ParseInt(scanIDParam, 10, 64)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "INVALID_SCAN_ID", "Scan ID must be a valid integer")
		return
	}

	scan, err := h.service.GetScanByID(c, userID, scanID)
	if err != nil {
		var appErr *domain.AppError
		if errors.As(err, &appErr) {
			sendErrorResponse(c, httpStatusFromAppError(appErr), appErr.Code, appErr.Message)
		} else {
			sendErrorResponse(c, http.StatusNotFound, "SCAN_NOT_FOUND", "Scan not found or access denied")
		}
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    scan,
	})
}
