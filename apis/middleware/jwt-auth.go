package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			response := utils.BuildFailedResponse("Token not found")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		token, err := jwtService.ValidateToken(authHeader)

		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				response := utils.BuildFailedResponse(err.Error())
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
				return
			}

			response := utils.BuildFailedResponse("Invalid token provided !")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			utils.TOKEN_ID = fmt.Sprintf("%v", claims["user_id"])
			utils.USER_TYPE = fmt.Sprintf("%v", claims["user_type"])

		}

	}
}
