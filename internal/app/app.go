package app

import (
	"fmt"
	"log"

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
	// connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", config.UNAMEDB, config.PASSDB, config.HOSTDB, config.DBNAME)
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s  sslmode=disable", config.UNAMEDB, config.DBNAME, config.PASSDB, config.HOSTDB))
	// db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	a.TX = db.MustBegin()
	a.DB = db
}

func (a *App) CreateRoutes() {
	r := gin.Default()

	v1 := r.Group("/v1")
	r.Use(middleware.CORS())
	{
		routes.TitleRoutes(v1)
		routes.UserRoutes(v1, a.TX)
	}

	a.Routes = r
}

func (a *App) Run() {
	a.Routes.Run(":8080")
}
