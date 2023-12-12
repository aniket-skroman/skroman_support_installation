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
	installation_serv = services.NewInstallationUserService(installation_repo, complaint_serv, complaint_repo)
	jwt_service := services.NewJWTService()
	installation_cont = controller.NewInstallationUserController(installation_serv)

	allocated_data := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		allocated_data.GET("/fetch-allocated-complaints", installation_cont.FetchAllocatedComplaintByEmp)
	}

	complaint_progress := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		complaint_progress.POST("/complaint_progress", installation_cont.CreateComplaintProgress)
		complaint_progress.GET("/complaint_progress/:obj_id", installation_cont.FetchComplaintProgress)
		complaint_progress.DELETE("/complaint_progress/:progress_id", installation_cont.DeleteComplaintProgress)
		complaint_progress.PUT("/complaint_progress/:complaint_id", installation_cont.MakeVerificationPendingStatus)
		complaint_progress.GET("/complaint_progress/complete/:allocate_id", installation_cont.FetchAllocatedCompletComplaint)
	}
}
