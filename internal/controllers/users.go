package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/model"
	"github.com/ctheil/pmdb-api/internal/repository"
	"github.com/ctheil/pmdb-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	DB *sql.DB
}

func NewUserController(db *sql.DB) UserControllerInterface {
	return &UserController{DB: db}
}

func (u *UserController) GetUser(c *gin.Context) {
	user_id := c.Param("id")
	if user_id == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "User ID is required."})
		return
	}
	db := u.DB
	user_repo := repository.NewUserRepository(db)
	user, err := user_repo.SelectUserByID(user_id)
	if err != nil {
		msg := fmt.Sprintf("No user found with id: %s", user_id)
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"result": user})
	}
}

func (u *UserController) InsertUser(c *gin.Context) {
	db := u.DB
	var post model.PostUser

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err})
		return
	}
	user_repo := repository.NewUserRepository(db)
	hashed_pwd, err := utils.HashPassword(post.HashedPW)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to hash password..."})
	}
	post.HashedPW = hashed_pwd
	post.RefreshTokenVersion = 1
	insert := user_repo.InsertUser(post)

	if insert {
		c.JSON(http.StatusOK, gin.H{"message": "user inserted successfully"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "user failed to insert"})
	}
}
