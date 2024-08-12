package api

import (
	"database/sql"
	"net/http"

	db "github.com/arkarsg/splitapp/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (s *Server) registerUser(r gin.IRoutes) {
	r.GET("/user", s.listUsers)
	r.GET("/user/:id", s.getUser)
	r.POST("/user", s.createUser)
}

type listUsersRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=50"`
}

func (s *Server) listUsers(c *gin.Context) {
	var req listUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	args := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := s.store.ListUsers(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, users)
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,number,min=1"`
}

func (s *Server) getUser(c *gin.Context) {
	var req getUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUserById(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	createUserArgs := db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
	}

	user, err := s.store.CreateUser(c, createUserArgs)
	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {
			case "foreign_key_violation":
				c.JSON(http.StatusForbidden, errorResponse((err)))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, user)
}
