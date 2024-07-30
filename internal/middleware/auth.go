package middleware

import (
	"fmt"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func Protect(tx *sqlx.Tx) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := auth.AuthenticateRequest(c, tx); err != nil {
			if err.Type == auth.ErrUnauthorized {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.Next()
	}
}

func Authenticate(tx *sqlx.Tx) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("\n\n AUTHENTICATE ROUTE ")
		if err := auth.AuthenticateRequest(c, tx); err != nil {
			if err.Type == auth.ErrServerError {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fmt.Printf("UNAUTHENTICATED: %s", err.Type)
		}
		// Continue event if Unauthorized.
		c.Next()
	}
}
