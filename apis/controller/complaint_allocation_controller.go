package controller

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/gin-gonic/gin"
)

type ComplaintAllocationController interface {
	AllocateComplaint(ctx *gin.Context)
	UpdateAllocateComplaint(ctx *gin.Context)
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
		cont.response = utils.RequestParamsMissingResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	// if jwt service failed to extract user id
	if utils.TOKEN_ID == "" {
		cont.response = utils.BuildFailedResponse("failed to allocate service")
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	// set a allocate by who handle the operation
	req.AllocateBy = utils.TOKEN_ID

	err := cont.allocation_serv.AllocateComplaint(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if strings.Contains(err.Error(), "already allocated") {
			ctx.JSON(http.StatusUnprocessableEntity, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DATA_INSERTED, utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *allocation_controller) UpdateAllocateComplaint(ctx *gin.Context) {
	var req dto.UpdateAllocateComplaintRequestDTO

	if err := ctx.ShouldBindQuery(&req); err != nil {
		cont.response = utils.RequestParamsMissingResponse(helper.Error_handler(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	// if jwt service failed to extract user id
	if utils.TOKEN_ID == "" {
		cont.response = utils.BuildFailedResponse("failed to update allocate complaint")
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	// set a allocate by who handle the operation
	req.AllocateBy = utils.TOKEN_ID

	err := cont.allocation_serv.UpdateComplaintAllocation(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse("invalid allocation detect")
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if strings.Contains(err.Error(), "completed complaint") {
			ctx.JSON(http.StatusUnprocessableEntity, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.UPDATE_SUCCESS, utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, cont.response)
}
