package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/service"
	"github.com/imdinnesh/openfinstack/services/kyc/models"
)

type KYCHandler struct {
  service service.KYCService
}

func NewKYCHandler(s service.KYCService) *KYCHandler {
  return &KYCHandler{service: s}
}

func (h *KYCHandler) SubmitKYC(c *gin.Context) {
  var input models.KYC
  if err := c.ShouldBindJSON(&input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  // Extract user ID from JWT context
  userID := c.GetUint("userID")
  input.UserID = userID

  if err := h.service.SubmitKYC(&input); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusCreated, gin.H{"message": "KYC submitted"})
}

func (h *KYCHandler) GetUserKYC(c *gin.Context) {
  userID := c.GetUint("userID")

  kycs, err := h.service.GetUserKYC(userID)
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, kycs)
}

func (h *KYCHandler) ListPending(c *gin.Context) {
  kycs, err := h.service.ListPending()
  if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }
  c.JSON(http.StatusOK, kycs)
}

func (h *KYCHandler) VerifyKYC(c *gin.Context) {
  idParam := c.Param("id")
  id, _ := strconv.Atoi(idParam)

  var req struct {
    Status string  `json:"status" binding:"required"`
    Reason *string `json:"reason"`
  }
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }

  adminID := c.GetUint("userID") // Admin ID

  if err := h.service.VerifyKYC(uint(id), req.Status, req.Reason, adminID); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
  }

  c.JSON(http.StatusOK, gin.H{"message": "KYC status updated"})
}
