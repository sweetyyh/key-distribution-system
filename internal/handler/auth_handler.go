package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"key-distribution-system/internal/config"
	"key-distribution-system/internal/pkg/jwtutil"
	"key-distribution-system/internal/pkg/response"
	"key-distribution-system/internal/store"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40020, "invalid request: "+err.Error())
		return
	}

	user, err := store.Global.Register(store.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, http.StatusConflict, 40021, err.Error())
		return
	}

	token, err := jwtutil.Sign(user.ID, string(user.Role), h.cfg.JWT.BuyerSecret, h.cfg.JWT.ExpireHours)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "sign token failed")
		return
	}

	response.Created(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40022, "invalid request: "+err.Error())
		return
	}

	user, err := store.Global.LoginByUsername(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 40101, "invalid credentials")
		return
	}

	secret := h.cfg.JWT.BuyerSecret
	if user.Role == "admin" {
		secret = h.cfg.JWT.AdminSecret
	}

	token, err := jwtutil.Sign(user.ID, string(user.Role), secret, h.cfg.JWT.ExpireHours)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 50001, "sign token failed")
		return
	}

	response.OK(c, gin.H{
		"token":    token,
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
	})
}
