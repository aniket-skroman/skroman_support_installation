package routers

import (
	"github.com/aniket-skroman/skroman_support_installation/apis"
	"github.com/aniket-skroman/skroman_support_installation/apis/controller"
	"github.com/aniket-skroman/skroman_support_installation/apis/middleware"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/gin-gonic/gin"
)

var (
	installation_repo repositories.InstallationUserRepository
	installation_serv services.InstallationUserService
	installation_cont controller.InstallationUserController
)

func InstallationUserRouter(router *gin.Engine, db *apis.Store) {
	installation_repo = repositories.NewInstallationUserRepository(db)
	installation_serv = services.NewInstallationUserService(installation_repo, complaint_serv)
	jwt_service := services.NewJWTService()
	installation_cont = controller.NewInstallationUserController(installation_serv)

	allocated_data := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		allocated_data.GET("/fetch-allocated-complaints/:allocated_to", installation_cont.FetchAllocatedComplaintByEmp)
	}
}
