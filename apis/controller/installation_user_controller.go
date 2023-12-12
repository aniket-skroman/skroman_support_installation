package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InstallationUserController interface {
	FetchAllocatedComplaintByEmp(ctx *gin.Context)
	CreateComplaintProgress(ctx *gin.Context)
	FetchComplaintProgress(ctx *gin.Context)
	DeleteComplaintProgress(ctx *gin.Context)
	MakeVerificationPendingStatus(ctx *gin.Context)
	FetchAllocatedCompletComplaint(ctx *gin.Context)
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
	var req dto.FetchAllocatedComplaintRequestDTO

	if err := ctx.ShouldBindQuery(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.installation_serv.FetchAllocatedComplaintByEmp(req)

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

func (cont *installation_cont) CreateComplaintProgress(ctx *gin.Context) {
	var req dto.CreateComplaintProgressRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.installation_serv.CreateComplaintProgress(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DATA_INSERTED, utils.COMPLAINT_DATA, result)
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *installation_cont) FetchComplaintProgress(ctx *gin.Context) {
	req := ctx.Param("obj_id")

	if req == "" {
		cont.response = utils.BuildFailedResponse(utils.REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.installation_serv.FetchComplaintProgress(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse(helper.Err_Data_Not_Found.Error())
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, result)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *installation_cont) DeleteComplaintProgress(ctx *gin.Context) {
	progress_id := ctx.Param("progress_id")

	if progress_id == "" {
		cont.response = utils.BuildFailedResponse(utils.REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.installation_serv.DeleteComplaintProgress(progress_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == helper.Err_Delete_Failed {
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DELETE_SUCCESS, utils.COMPLAINT_DATA, nil)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *installation_cont) MakeVerificationPendingStatus(ctx *gin.Context) {
	complaint_id := ctx.Param("complaint_id")

	if complaint_id == "" {
		cont.response = utils.BuildFailedResponse(helper.ERR_REQUIRED_PARAMS.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.installation_serv.MakeVerificationPendingStatus(complaint_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if err == helper.ERR_INVALID_ID || err == helper.Err_Update_Failed {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.UPDATE_SUCCESS, utils.COMPLAINT_DATA, nil)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *installation_cont) FetchAllocatedCompletComplaint(ctx *gin.Context) {
	allocated_id := ctx.Param("allocate_id")

	if allocated_id == "" {
		cont.response = utils.BuildFailedResponse(utils.REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.installation_serv.FetchAllocatedCompletComplaint(allocated_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if err == helper.ERR_INVALID_ID {
			ctx.JSON(http.StatusBadRequest, cont.response)
			return
		} else if err == sql.ErrNoRows {
			cont.response = utils.BuildFailedResponse(helper.Err_Data_Not_Found.Error())
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, result)
	ctx.JSON(http.StatusOK, cont.response)
}
