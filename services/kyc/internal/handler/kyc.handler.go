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

// SubmitKYC handles submission of KYC by user
func (h *KYCHandler) SubmitKYC(c *gin.Context) {
	var req struct {
		DocumentType string `json:"document_type" binding:"required"`
		DocumentURL  string `json:"document_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.Request.Header.Get("X-User-ID")
	userID64, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID64)
	input := &models.KYC{	
		UserID:       userID,
		DocumentType: req.DocumentType,
		DocumentURL:  req.DocumentURL,
		Status:       "pending",
	}

	if err := h.service.SubmitKYC(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "KYC submitted"})
}

// GetUserKYC returns KYC records of the logged-in user
func (h *KYCHandler) GetUserKYC(c *gin.Context) {
	userIDStr := c.Request.Header.Get("X-User-ID")
	userID64, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID64)

	kycs, err := h.service.GetUserKYC(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kycs)
}

// ListPending lists all pending KYCs for admin review
func (h *KYCHandler) ListPending(c *gin.Context) {
	kycs, err := h.service.ListPending()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, kycs)
}

// VerifyKYC handles KYC verification by admin
func (h *KYCHandler) VerifyKYC(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid KYC ID"})
		return
	}

	var req struct {
		Status string  `json:"status" binding:"required"`
		Reason *string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// If rejecting, Reason must be provided
	if req.Status == "rejected" && (req.Reason == nil || *req.Reason == "") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reason must be provided when rejecting KYC"})
		return
	}

	adminIDStr := c.Request.Header.Get("X-User-ID")
	adminID64, _ := strconv.ParseUint(adminIDStr, 10, 64)
	adminID := uint(adminID64)

	if err := h.service.VerifyKYC(uint(id), req.Status, req.Reason, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "KYC status updated"})
}

func (h *KYCHandler) GetKYCStatusByUserID(c *gin.Context) {
	userIDStr := c.Request.Header.Get("X-User-ID")
	userID64, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID64)

	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	status, err := h.service.GetKYCStatusByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

func (h *KYCHandler) UpdateKYCStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed","message":err.Error()})
		return
	}

	var req struct {
		Status string  `json:"status" binding:"required"`
		Reason *string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status":"failed","message":err.Error()})
		return
	}

	adminIDStr := c.Request.Header.Get("X-Admin-ID")
	adminID64, _ := strconv.ParseUint(adminIDStr, 10, 64)
	adminID := uint(adminID64)

	if err := h.service.UpdateKYCStatus(uint(id), req.Status, req.Reason, adminID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status":"failed","message":err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status":"success","message": "KYC status updated"})
}
