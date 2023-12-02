package main

import (
	"database/sql"
	"log"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	"github.com/aniket-skroman/skroman_support_installation/apis/database"
	"github.com/aniket-skroman/skroman_support_installation/apis/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://13.234.110.115:3000"}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
	return corsConfig
}

const (
	ContentTypeBinary = "application/octet-stream"
	ContentTypeForm   = "application/x-www-form-urlencoded"
	ContentTypeJSON   = "application/json"
	ContentTypeHTML   = "text/html; charset=utf-8"
	ContentTypeText   = "text/plain; charset=utf-8"
)

type APIServer struct{}

func (api *APIServer) make_db_connection() (*sql.DB, error) {
	db, err := database.DB_INSTANCE()
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}

func (api *APIServer) init_app_route() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(CORSConfig()))
	//r.Use(middleware.ErrorHandleMiddleware())

	return r
}

func (api *APIServer) make_app_route(route *gin.Engine, db *sql.DB) {
	route.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	store := apis.NewStore(db)

	routers.ComplaintRouter(route, store)
	routers.ComplaintAllocationRouter(route, store)
	routers.InstallationUserRouter(route, store)
}

func (api *APIServer) run_app(route *gin.Engine) error {
	return route.Run(":9000")
}

func main() {

	app_server := APIServer{}

	db, _ := app_server.make_db_connection()

	defer database.CloseDBConnection(db)

	route := app_server.init_app_route()

	app_server.make_app_route(route, db)

	if err := app_server.run_app(route); err != nil {
		log.Fatal(err)
	}

}
