package controller

import (
	"net/http"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/gin-gonic/gin"
)

type ComplaintAllocationController interface {
	AllocateComplaint(ctx *gin.Context)
}

type allocation_controller struct {
	allocation_serv services.ComplaintAllocationService
	response        map[string]interface{}
}

func NewComplaintAllocationController(ser services.ComplaintAllocationService) ComplaintAllocationController {
	return &allocation_controller{
		allocation_serv: ser,
		response:        map[string]interface{}{},
	}
}

func (cont *allocation_controller) AllocateComplaint(ctx *gin.Context) {
	var req dto.CreateAllocationRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.allocation_serv.AllocateComplaint(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DATA_INSERTED, utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusCreated, cont.response)
}
