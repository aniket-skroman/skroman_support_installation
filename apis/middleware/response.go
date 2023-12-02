package middleware

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		fmt.Println("Error handler get called..")
		for _, err := range ctx.Errors {
			fmt.Println("Error from middleware : ", err)

			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
				ctx.Abort()
				return
			}
		}
		//
	}
}
