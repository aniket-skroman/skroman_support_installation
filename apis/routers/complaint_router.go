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
	complaint_repo  repositories.ComplaintRepository
	allocation_repo repositories.ComplaintAllocationRepository
	allocation_ser  services.ComplaintAllocationService
	jwt_service     services.JWTService
	complaint_serv  services.ComplaintService
	complaint_cont  controller.ComplaintController
)

func ComplaintRouter(router *gin.Engine, store *apis.Store) {
	complaint_repo = repositories.NewComplaintRepository(store)
	allocation_repo = repositories.NewComplaintAllocationRepository(store)
	allocation_ser = services.NewComplaintAllocationService(allocation_repo)
	jwt_service = services.NewJWTService()
	complaint_serv = services.NewComplaintService(complaint_repo, jwt_service, allocation_ser)
	complaint_cont = controller.NewComplaintController(complaint_serv)
	complaint := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		complaint.POST("/create-complaint", complaint_cont.CreateComplaint)
		complaint.GET("/fetch-complaints/:page_id/:page_size/:tag_key", complaint_cont.FetchAllComplaints)
		complaint.GET("/fetch-complaint", complaint_cont.FetchComplaintDetailByComplaint)
		complaint.POST("/upload-device-image", complaint_cont.UploadDeviceImage)
		complaint.POST("/upload-device-video", complaint_cont.UploadDeviceVideo)

		complaint.PUT("/update-complaint", complaint_cont.UpdateComplaintInfo)
		complaint.DELETE("/delete-device-file/:file_id", complaint_cont.DeleteDeviceFiles)
		complaint.DELETE("/delete-complaint/:complaint_id", complaint_cont.DeleteComplaint)

		complaint.PUT("/complaint-resolved/:complaint_id", complaint_cont.ComplaintResolve)
		complaint.GET("/fetch-complaint-counts", complaint_cont.FetchAllComplaintCounts)
	}

	client := router.Group("/api", middleware.AuthorizeJWT(jwt_service))
	{
		client.POST("/client-registraion", complaint_cont.ClientRegistration)
		client.GET("/client-complaints/:client_id/:page_id/:page_size", complaint_cont.FetchComplaintsByClient)
	}

	device_img := router.Group("/api")
	{
		device_img.GET("/device-file/:directory/:image_path", complaint_cont.FetchDeviceImageURL)
	}
}
