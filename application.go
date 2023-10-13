package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/aniket-skroman/skroman_support_installation/apis"
	"github.com/aniket-skroman/skroman_support_installation/apis/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var DB_DRIVER = "postgres"
var DB_SOURCE = "postgresql://postgres:support12@skroman-user.ckwveljlsuux.ap-south-1.rds.amazonaws.com:5432/skroman_client_complaints"

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
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

func main() {
	db, err := sql.Open(DB_DRIVER, DB_SOURCE)

	if err != nil {
		log.Fatal(err)
	}

	store := apis.NewStore(db)

	router := gin.New()

	router.GET("/", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, ContentTypeHTML, []byte("<html>Program file run...</html>"))
	})

	router.Use(cors.New(CORSConfig()))
	router.Static("static", "static")

	routers.ComplaintRouter(router, store)

	router.Run(":8181")
}
