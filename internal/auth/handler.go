package auth

import (
	"be-golang-poliklinik/internal/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service Service
}

func NewHandler(service Service) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input PasienRegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, helper.FormatValidatorError(err))
		return
	}

	newUser, err := h.service.Register(input)
	if err != nil {
		helper.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	helper.Success(c, http.StatusCreated, "Registrasi akun berhasil!", newUser)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		helper.Error(c, http.StatusBadRequest, helper.FormatValidatorError(err))
		return
	}

	token, userData, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		helper.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Login berhasil", gin.H{
		"access_token": token,
		"user":         userData,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString := c.GetString("clean_token")
	if tokenString == "" {
		helper.Error(c, http.StatusInternalServerError, "Gagal memproses token sistem")
		return
	}

	err := h.service.Logout(tokenString)
	if err != nil {
		helper.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.Success(c, http.StatusOK, "Logout berhasil, token telah dinonaktifkan", nil)
}
