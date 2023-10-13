package controller

import (
	"net/http"
	"strings"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/gin-gonic/gin"
)

type ComplaintController interface {
	CreateComplaint(*gin.Context)
	FetchAllComplaints(*gin.Context)
}

type complaint_controller struct {
	comp_serv services.ComplaintService
}

func NewComplaintController(serv services.ComplaintService) ComplaintController {
	return &complaint_controller{
		comp_serv: serv,
	}
}

func (cont *complaint_controller) CreateComplaint(ctx *gin.Context) {
	var req dto.CreateComplaintRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	complaint_info, err := cont.comp_serv.CreateComplaint(req)
	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.COMPLAINT_CREATED, utils.COMPLAINT_DATA, complaint_info)
	ctx.JSON(http.StatusOK, response)
}

func (cont *complaint_controller) FetchAllComplaints(ctx *gin.Context) {
	var req dto.PaginationRequestParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	complaints, err := cont.comp_serv.FetchAllComplaints(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, complaints)
	ctx.JSON(http.StatusOK, response)
}
