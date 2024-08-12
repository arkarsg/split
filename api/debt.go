package api

import (
	"database/sql"
	"net/http"

	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerDebt(r gin.IRoutes) {
	r.GET("/debt/:id", s.getDebt)
	r.GET("/debt", s.getDebtByTransactionID)
	r.GET("/debt/:id/debtors", s.getDebtorsByDebt)
	r.POST("/debt/:id/debtors", s.createDebtor)
	r.POST("/debt", s.createDebt)
}

type getDebtRequest struct {
	ID int64 `uri:"id" binding:"required,number,min=1"`
}

func (s *Server) getDebt(c *gin.Context) {
	var req getDebtRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	debt, err := s.store.GetDebtById(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, debt)
}

type getDebtByTransactionRequest struct {
	TransactionID int64 `form:"transaction_id" binding:"required,min=1"`
}

func (s *Server) getDebtByTransactionID(c *gin.Context) {
	var req getDebtByTransactionRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	debt, err := s.store.GetDebtByTransactionId(c, req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, debt)
}

func (s *Server) getDebtorsByDebt(c *gin.Context) {
	var req getDebtRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	debtors, err := s.store.GetDebtDebtorsByDebtId(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, debtors)
}

type createDebtorRequest struct {
	DebtID   int64       `uri:"id" binding:"required,number,min=1"`
	DebtorID int64       `json:"transaction_id" binding:"required"`
	Amount   string      `json:"amount" binding:"required,moneyamount"`
	Currency db.Currency `json:"currency" binding:"required,supportedcurrency"`
}

func (s *Server) createDebtor(c *gin.Context) {
	var req createDebtorRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateDebtDebtorsParams{
		DebtID:   req.DebtID,
		DebtorID: req.DebtorID,
		Amount:   req.Amount,
		Currency: req.Currency,
	}

	dd, err := s.store.CreateDebtDebtors(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, dd)
}

type createDebtRequest struct {
	TransactionID int64 `json:"transaction_id" binding:"required"`
}

func (s *Server) createDebt(c *gin.Context) {
	var req createDebtRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	debt, err := s.store.CreateDebt(c, req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, debt)
}
