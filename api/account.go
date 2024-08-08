package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/arkarsg/splitapp/db/sqlc"
	u "github.com/arkarsg/splitapp/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (s *Server) registerAccount(r *gin.Engine) {
	r.GET("/account/:username", s.getAccount)
	r.POST("/account", s.createAccount)
}

type getAccountRequest struct {
	Username string `uri:"username" binding:"required,alphanum"`
}

func (s *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)
}

type createAccountRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createAccountResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := u.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	createAccountArgs := db.CreateAccountParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.Fullname,
		Email:          req.Email,
	}

	account, err := s.store.CreateAccount(c, createAccountArgs)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse((err)))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// wrap Account in a response struct to hide hashed password
	resp := createAccountResponse{
		Username:          account.Username,
		Email:             account.Email,
		FullName:          account.FullName,
		PasswordChangedAt: account.PasswordChangedAt,
		CreatedAt:         account.CreatedAt,
	}
	c.JSON(http.StatusOK, resp)
}
