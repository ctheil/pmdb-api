package controllers

import (
	"net/http"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserController struct {
	TX *sqlx.Tx
}

func NewUserController(tx *sqlx.Tx) *UserController {
	return &UserController{TX: tx}
}

func (u *UserController) GetUser(c *gin.Context) {
	user, exists := c.Get("user")
	user, ok := user.(model.User)
	if !exists || !ok {
		c.JSON(http.StatusOK, gin.H{"message": "unauthenticated hello world!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "authenticated hello world!"})
}
