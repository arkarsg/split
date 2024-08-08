package api

import (
	"database/sql"
	"net/http"

	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerPayment(r *gin.Engine) {
	r.POST("/payments", s.createPaymentTx)
}

type createPaymentTxRequest struct {
	DebtId   int64       `json:"debt_id" binding:"required,min=1"`
	DebtorId int64       `json:"debtor_id" binding:"required,min=1"`
	Amount   string      `json:"amount" binding:"required,moneyamount"`
	Currency db.Currency `json:"currency" binding:"required,supportedcurrency"`
}

func (s *Server) createPaymentTx(c *gin.Context) {
	var req createPaymentTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// validation here
	if !s.isValidPaymentRequest(c, req) {
		return
	}

	arg := db.SettleDebtPaymentTxParams{
		DebtId:   req.DebtId,
		DebtorId: req.DebtorId,
		Amount:   req.Amount,
		Currency: req.Currency,
	}

	res, err := s.store.SettleDebtPaymentsTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, res)
}

func (s *Server) isValidPaymentRequest(c *gin.Context, req createPaymentTxRequest) bool {
	return s.debtExists(c, req.DebtId) && s.debtorExists(c, req.DebtorId)
}

func (s *Server) debtExists(c *gin.Context, debtID int64) bool {
	_, err := s.store.GetDebtById(c, debtID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	return true
}

func (s *Server) debtorExists(c *gin.Context, debtorId int64) bool {
	_, err := s.store.GetUserById(c, debtorId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	return true
}
