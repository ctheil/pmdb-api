package app

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ctheil/pmdb-api/internal/config"
	"github.com/ctheil/pmdb-api/internal/routes"
	"github.com/gin-gonic/gin"
)

type App struct {
	DB     *sql.DB
	Routes *gin.Engine
}

func (a *App) CreateConnection() {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.UNAMEDB, config.PASSDB, config.HOSTDB, config.DBNAME)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
}

func (a *App) CreateRoutes() {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		routes.MovieRoutes(v1)
	}

	a.Routes = r
}

func (a *App) Run() {
	a.Routes.Run(":8080")
}
