package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/arkarsg/splitapp/token"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerTransaction(r gin.IRoutes) {
	r.GET("/transaction/:id", s.getTransactionById)
	r.GET("/transaction", s.listTransactions)
	r.POST("/transaction", s.createTransactionTx)
}

type getTransactionRequest struct {
	ID int64 `uri:"id" binding:"required,number,min=1"`
}

func (s *Server) getTransactionById(c *gin.Context) {
	var req getTransactionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	transaction, err := s.store.GetTransactionById(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if transaction.PayerUsername != authPayload.Username {
		err := errors.New("Transaction does not belong to the authorized user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, transaction)
}

type listTransactionsRequest struct {
	PayerID  int64 `form:"payer_id" binding:"required,min=1"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (s *Server) listTransactions(c *gin.Context) {
	var req listTransactionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	user, err := s.store.GetUserById(c, req.PayerID)
	if user.Username != authPayload.Username {
		err := errors.New("Unauthorized to view user's transactions")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	args := db.GetTransactionsByPayerParams{
		PayerID: req.PayerID,
		Limit:   req.PageSize,
		Offset:  (req.PageID - 1) * req.PageSize,
	}

	transactions, err := s.store.GetTransactionsByPayer(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, transactions)
}

type CreateTransactionDebtTxRequest struct {
	Amount   string      `json:"amount" binding:"required,moneyamount"`
	Currency db.Currency `json:"currency" binding:"required,supportedcurrency"`
	Title    string      `json:"title" binding:"required"`
	PayerID  int64       `json:"payer_id" binding:"required"`
}

func (s *Server) createTransactionTx(c *gin.Context) {
	var req CreateTransactionDebtTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.CreateTransactionDebtTxParams{
		Amount:   req.Amount,
		Currency: req.Currency,
		Title:    req.Title,
		PayerID:  req.PayerID,
	}

	res, err := s.store.CreateTransactionDebtTx(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, res)
}
