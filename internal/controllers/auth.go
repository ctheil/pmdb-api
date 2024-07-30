package controllers

import (
	"fmt"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/auth"
	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/internal/repository"
	"github.com/ctheil/pmdb-api/internal/services"
	"github.com/ctheil/pmdb-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type AuthController struct {
	TX *sqlx.Tx
}

func NewAuthController(tx *sqlx.Tx) *AuthController {
	return &AuthController{TX: tx}
}

func (a *AuthController) Login(c *gin.Context) {
	type LoginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	lr := LoginReq{}
	if err := c.ShouldBindJSON(&lr); err != nil {
		// if err := json.NewDecoder(c.Request.Body).Decode(&lr); err != nil {
		msg := fmt.Sprintf("Invalid username or password: %v", lr)
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}

	user_repo := repository.NewUserRepository(a.TX)
	user, err := user_repo.GetByUsername(lr.Username)
	if err != nil || user.Id == 0 {
		fmt.Printf("error: %e", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect username or password.", "error": err})
		return
	}
	if !utils.VerifyPassowrd([]byte(user.Password), []byte(lr.Password)) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Incorrect password.", "error": err})
		return
	}

	tokens, err := auth.NewTokens(user, user_repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error generating auth tokens", "error": err})
		return
	}
	auth.SetTokenCookies(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Welcome, %s", user.Username), "user": user, "access_token": tokens.AccessToken, "refresh_token": tokens.RefreshToken})
}

func (a *AuthController) Signup(c *gin.Context) {
	type SignupRequest struct {
		Username string `json:"username" validate:"required,min=3,max=32"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	sr := SignupRequest{}
	if err := c.ShouldBindJSON(&sr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid username, password, and/or email"})
		return
	}

	errs := services.ValidateStruct(sr)
	if errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errs})
		return
	}

	// Email and username check
	user_repo := repository.NewUserRepository(a.TX)

	user, err := user_repo.GetByUsername(sr.Username)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if user.Id != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is already taken"})
		return
	}
	user, err = user_repo.GetByEmail(sr.Email)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if user.Id != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is already taken"})
		return
	}

	// ALL GOOD!
	hashed_pw, err := utils.HashPassword(sr.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	new_user := model.PostUser{
		Username: sr.Username,
		Email:    sr.Email,
		Password: hashed_pw,
	}
	_, ok := user_repo.Insert(new_user)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save new user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success! Please login."})
}
