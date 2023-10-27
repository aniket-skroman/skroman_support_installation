package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstallationUserController interface {
	FetchAllocatedComplaintByEmp(ctx *gin.Context)
}

type installation_cont struct {
	installation_serv services.InstallationUserService
	response          map[string]interface{}
}

func NewInstallationUserController(ser services.InstallationUserService) InstallationUserController {
	return &installation_cont{
		installation_serv: ser,
		response:          make(map[string]interface{}),
	}
}

func (cont *installation_cont) FetchAllocatedComplaintByEmp(ctx *gin.Context) {
	allocated_id := ctx.Param("allocated_to")

	if allocated_id == " " {
		cont.response = utils.BuildFailedResponse("please provide a required params")
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}
	result, err := cont.installation_serv.FetchAllocatedComplaintByEmp(allocated_id)

	if err != nil {
		if uuid.IsInvalidLengthError(err) {
			cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse(utils.DATA_NOT_FOUND)
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, result)
	ctx.JSON(http.StatusOK, cont.response)
}
