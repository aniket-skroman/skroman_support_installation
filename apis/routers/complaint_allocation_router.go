package routers

import (
	"github.com/aniket-skroman/skroman_support_installation/apis"
	"github.com/aniket-skroman/skroman_support_installation/apis/controller"
	"github.com/aniket-skroman/skroman_support_installation/apis/middleware"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/gin-gonic/gin"
)

func ComplaintAllocationRouter(router *gin.Engine, store *apis.Store) {
	var (
		allocation_repo = repositories.NewComplaintAllocationRepository(store)
		jwtService      = services.NewJWTService()
		allocation_ser  = services.NewComplaintAllocationService(allocation_repo)
		allocation_cont = controller.NewComplaintAllocationController(allocation_ser)
	)

	allocation := router.Group("/api", middleware.AuthorizeJWT(jwtService))
	{
		allocation.POST("/allocate-complaint", allocation_cont.AllocateComplaint)
	}
}
