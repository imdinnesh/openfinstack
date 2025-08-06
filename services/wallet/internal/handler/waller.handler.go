package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/service"
)

type WalletHandler struct {
	Service service.WalletService
}

func New(service service.WalletService) *WalletHandler {
	return &WalletHandler{Service: service}
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.Service.CreateWallet(req.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "wallet created"})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userID"))

	wallet, err := h.Service.GetWallet(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if wallet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}
func (h *WalletHandler) AddFunds(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userID"))

	var req struct {
		Amount int64  `json:"amount" binding:"required"`
		Desc   string `json:"desc"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Service.AddFunds(c.Request.Context(), uint(userID), req.Amount, req.Desc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "funds added"})
}
func (h *WalletHandler) WithdrawFunds(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userID"))

	var req struct {
		Amount int64  `json:"amount" binding:"required"`
		Desc   string `json:"desc"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	if err := h.Service.WithdrawFunds(c.Request.Context(), uint(userID), req.Amount, req.Desc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "funds withdrawn"})
}
func (h *WalletHandler) Transfer(c *gin.Context) {
	var req struct {
		FromUserID uint  `json:"from_user_id" binding:"required"`
		ToUserID   uint  `json:"to_user_id" binding:"required"`
		Amount     int64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transfer request"})
		return
	}

	if err := h.Service.Transfer(c.Request.Context(), req.FromUserID, req.ToUserID, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}
func (h *WalletHandler) GetTransactions(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("userID"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	txs, err := h.Service.GetTransactions(uint(userID), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": txs})
}
