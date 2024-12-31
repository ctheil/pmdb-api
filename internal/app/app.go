package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ctheil/pmdb-api/internal/config"
	"github.com/ctheil/pmdb-api/internal/middleware"
	"github.com/ctheil/pmdb-api/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type App struct {
	TX     *sqlx.Tx
	DB     *sqlx.DB
	Routes *gin.Engine
}

func (a *App) CreateConnection() {
	// NOTE: without sqlx ORM:
	// connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.UNAMEDB, config.PASSDB, config.HOSTDB, config.DBNAME)
	// db, err := sql.Open("postgres", connStr)

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s  sslmode=disable", config.UNAMEDB, config.DBNAME, config.PASSDB, config.HOSTDB))
	if err != nil {
		log.Fatal(err)
	}
	a.TX = db.MustBegin()
	a.DB = db
}

func (a *App) CreateRoutes() {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello world"})
	})
	v1 := r.Group("/v1")
	{
		routes.TitleRoutes(v1)
		routes.AuthRoutes(v1, a.TX)
		routes.OAuthRoutes(v1, a.TX)
		routes.UserRoutes(v1, a.TX)
	}

	a.Routes = r
}

func (a *App) Run() {
	a.Routes.Run(":8080")
}
