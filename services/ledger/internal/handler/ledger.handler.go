package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/imdinnesh/openfinstack/services/ledger/internal/service"
)

type LedgerHandler struct{ svc service.LedgerService }

func NewLedgerHandler(svc service.LedgerService) *LedgerHandler { return &LedgerHandler{svc: svc} }

func (h *LedgerHandler) CreateTransaction(c *gin.Context) {
	var req service.CreateTxnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Allow Idempotency-Key header to map to ExternalRef if not provided
	if req.ExternalRef == "" {
		if hdr := c.GetHeader("Idempotency-Key"); hdr != "" {
			req.ExternalRef = hdr
		}
	}

	res, err := h.svc.CreateAndPost(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *LedgerHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	res, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *LedgerHandler) ReverseTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	reason := c.Query("reason")
	res, err := h.svc.Reverse(id, reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *LedgerHandler) ListEntriesByAccount(c *gin.Context) {
	accountID := c.Param("accountId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	entries, err := h.svc.ListEntries(accountID, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}
