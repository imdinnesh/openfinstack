package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/service"
)

type WalletHandler struct {
	Service *service.Service
}

func New(service *service.Service) *WalletHandler {
	return &WalletHandler{Service: service}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID := c.Param("userId")
	wallet, err := h.Service.Repo.GetWallet(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID, "balance": wallet.Balance})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	var req struct {
		FromUserID string `json:"from_user_id"`
		ToUserID   string `json:"to_user_id"`
		Amount     int64  `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if err := h.Service.Transfer(req.FromUserID, req.ToUserID, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "transfer successful"})
}

func (h *WalletHandler) Transactions(c *gin.Context) {
	userID := c.Param("userId")
	txns, err := h.Service.Repo.GetTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, txns)
}