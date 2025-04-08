package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/obadoraibu/go-auth/internal/domain"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SignUp(c *gin.Context) {
	r := &domain.UserSignUpInput{}
	if err := c.ShouldBindJSON(&r); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.service.SignUp(c, r)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) SignIn(c *gin.Context) {
	r := &domain.UserSignInInput{}
	if err := c.ShouldBindJSON(&r); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.service.SignIn(c, r)
	if err != nil {
		if err == domain.ErrWrongEmailOrPassword {
			sendErrorResponse(c, http.StatusConflict, err.Error())
			return
		}
		if err == domain.ErrEmailIsNotConfirmed {
			sendErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}
		sendErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "refresh",
		Value:    resp.RefreshToken,
		Path:     "/",
		MaxAge:   86400 * 60,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Refresh(c *gin.Context) {
	type request struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	req := &request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	cookie, err := c.Cookie("refresh")
	if err != nil {
		sendErrorResponse(c, http.StatusUnauthorized, "no authorization cookie")
	}

	response, err := h.service.Refresh(cookie, req.Fingerprint)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	logrus.Println(cookie)

	newCookie := &http.Cookie{
		Name:     "refresh",
		Value:    response.RefreshToken,
		Path:     "/",
		MaxAge:   86400 * 60,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, newCookie)

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Revoke(c *gin.Context) {
	type request struct {
		Fingerprint string `json:"fingerprint" binding:"required"`
	}

	req := &request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	cookie, err := c.Cookie("refresh")
	if err != nil {
		sendErrorResponse(c, http.StatusUnauthorized, "no authorization cookie")
	}

	err = h.service.Revoke(cookie, req.Fingerprint)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ConfirmEmail(c *gin.Context) {
	code := c.Param("code")

	err := h.service.ConfirmEmail(code)
	if err != nil {
		if err == domain.ErrWrongEmailConfirmationCode {
			sendErrorResponse(c, http.StatusBadRequest, "wrong confirmation code")
			return
		}
		sendErrorResponse(c, http.StatusInternalServerError, "cannot confirm email")
		return
	}
	logrus.Println("email confirmed")
	c.Redirect(http.StatusFound, "http://localhost:3000/sign-in")
}
