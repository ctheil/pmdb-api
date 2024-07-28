package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/internal/repository"
	"github.com/ctheil/pmdb-api/internal/services"
	"github.com/ctheil/pmdb-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
)

type UserController struct {
	TX *sqlx.Tx
}

func NewUserController(tx *sqlx.Tx) *UserController {
	return &UserController{TX: tx}
}

func (u *UserController) GetUser(c *gin.Context) {
	user_id := c.Param("id")
	if user_id == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "User ID is required."})
		return
	}
	tx := u.TX
	user_repo := repository.NewUserRepository(tx)
	user, err := user_repo.GetById(user_id)
	if err != nil {
		msg := fmt.Sprintf("No user found with id: %s", user_id)
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"result": user})
	}
}

func (u *UserController) InsertUser(c *gin.Context) {
	tx := u.TX
	var post model.PostUser

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}
	user_repo := repository.NewUserRepository(tx)
	hashed_pwd, err := utils.HashPassword(post.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to hash password..."})
	}
	post.Password = hashed_pwd
	post.RefreshTokenVersion = 0
	id, ok := user_repo.Insert(post)

	if ok {
		c.JSON(http.StatusOK, gin.H{"message": "user inserted successfully", "user_id": id})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user failed to insert"})
	}
}

func (u *UserController) Login(c *gin.Context) {
	type LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	lr := LoginReq{}
	if err := json.NewDecoder(c.Request.Body).Decode(&lr); err != nil {
		msg := fmt.Sprintf("Invalid username or password: %v", lr)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	user_repo := repository.NewUserRepository(u.TX)
	user, err := user_repo.GetByUsername(lr.Username)
	if err != nil {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect username or password.", "error": err})
		return
	}
	if !utils.VerifyPassowrd([]byte(user.Password), []byte(lr.Password)) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect password.", "error": err})
		return
	}

	userClaims := services.UserClaims{
		Id:       user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}

	signedAccessToken, err := services.NewAccessToken(userClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate access token"})
		return
	}

	refreshClaims := services.RefreshClaims{
		Version: user.RefreshTokenVersion + 1,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 48)),
		},
	}

	signedRefreshToken, err := services.NewRefreshToken(refreshClaims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate refresh token"})
		return
	}
	user.RefreshTokenVersion++
	if err := user.Save(*u.TX); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to update refresh token version."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome, %s", user.Username), "user": user, "accessToken": signedAccessToken, "refreshToken": signedRefreshToken})
}
