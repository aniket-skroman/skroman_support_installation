package routers

import (
	"github.com/aniket-skroman/skroman_support_installation/apis"
	"github.com/aniket-skroman/skroman_support_installation/apis/controller"
	"github.com/aniket-skroman/skroman_support_installation/apis/middleware"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/gin-gonic/gin"
)

func ComplaintRouter(router *gin.Engine, store *apis.Store) {
	var (
		complaint_repo = repositories.NewComplaintRepository(store)
		jwt_service    = services.NewJWTService()
		complaint_serv = services.NewComplaintService(complaint_repo)
		complaint_cont = controller.NewComplaintController(complaint_serv)
	)

	complaint := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		complaint.POST("/create-complaint", complaint_cont.CreateComplaint)
		complaint.GET("/fetch-complaints/:page_id/:page_size", complaint_cont.FetchAllComplaints)
	}

	dummy_img := router.Group("/api")
	{
		dummy_img.GET("/dummy/:directory/:image_path", complaint_cont.FetchDummyURL)
	}
}
