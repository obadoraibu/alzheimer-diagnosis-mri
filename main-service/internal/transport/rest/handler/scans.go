package handler

import (
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
		sendErrorResponse(c, http.StatusBadRequest, "missing form fields")
		return
	}

	patientAge, err := strconv.Atoi(patientAgeStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid patient_age")
		return
	}

	scanDate, err := time.Parse("2006-01-02", scanDateStr)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid scan_date, expected YYYY-MM-DD")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "failed to read file")
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
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "analysis created and scan uploaded"})
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
			sendErrorResponse(c, http.StatusBadRequest, "invalid scan id")
			return
		}
	}

	if fromStr := c.Query("uploaded_from"); fromStr != "" {
		if from, err := time.Parse("2006-01-02", fromStr); err == nil {
			f.UploadedFrom = &from
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "invalid uploaded_from date")
			return
		}
	}

	if toStr := c.Query("uploaded_to"); toStr != "" {
		if to, err := time.Parse("2006-01-02", toStr); err == nil {
			f.UploadedTo = &to
		} else {
			sendErrorResponse(c, http.StatusBadRequest, "invalid uploaded_to date")
			return
		}
	}

	scans, err := h.service.GetScansByFilters(userID, f)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "failed to retrieve scans")
		return
	}

	c.JSON(http.StatusOK, scans)
}

func (h *Handler) GetScanDetail(c *gin.Context) {
	token := c.MustGet("AccessToken").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	scanIDParam := c.Param("id")
	scanID, err := strconv.ParseInt(scanIDParam, 10, 64)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid scan id")
		return
	}

	scan, err := h.service.GetScanByID(userID, scanID)
	if err != nil {
		sendErrorResponse(c, http.StatusNotFound, "scan not found or access denied")
		return
	}

	c.JSON(http.StatusOK, scan)
}
