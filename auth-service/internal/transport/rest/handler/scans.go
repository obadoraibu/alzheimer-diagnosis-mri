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
