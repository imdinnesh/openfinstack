package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/imdinnesh/openfinstack/services/ledger/config"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/handler"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/repository"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/service"
)

func RegisterLedgerRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	ledgerRepo := repository.NewLedgerRepository(db)
	ledgerService := service.NewLedgerService(db, ledgerRepo)
	ledgerHandler := handler.NewLedgerHandler(ledgerService)

	ledger := r.Group("/ledger")
	ledger.POST("/transaction", ledgerHandler.CreateTransaction)
	ledger.GET("/transaction/:id", ledgerHandler.GetTransaction)
	ledger.POST("/transaction/:id/reverse", ledgerHandler.ReverseTransaction)
	ledger.GET("/accounts/:accountId/entries", ledgerHandler.ListEntriesByAccount)
}