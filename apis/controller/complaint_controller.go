package controller

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/connections"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ComplaintController interface {
	CreateComplaint(*gin.Context)
	FetchAllComplaints(*gin.Context)
	FetchDeviceImageURL(*gin.Context)
	FetchComplaintDetailByComplaint(*gin.Context)
	UploadDeviceImage(*gin.Context)
	UploadDeviceVideo(*gin.Context)
	UpdateComplaintInfo(ctx *gin.Context)
	DeleteDeviceFiles(ctx *gin.Context)
	DeleteComplaint(ctx *gin.Context)
	ComplaintResolve(ctx *gin.Context)
	FetchAllComplaintCounts(ctx *gin.Context)
	FetchComplaintsByClient(ctx *gin.Context)

	ClientRegistration(*gin.Context)
	DeleteClient(*gin.Context)

	FetchPDFFile(ctx *gin.Context)
}

type complaint_controller struct {
	comp_serv services.ComplaintService
	response  map[string]interface{}
}

func NewComplaintController(serv services.ComplaintService) ComplaintController {
	return &complaint_controller{
		comp_serv: serv,
		response:  map[string]interface{}{},
	}
}

func (cont *complaint_controller) CreateComplaint(ctx *gin.Context) {
	var req dto.CreateComplaintRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	complaint_info, err := cont.comp_serv.CreateComplaint(req)
	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.COMPLAINT_CREATED, utils.COMPLAINT_DATA, complaint_info)
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *complaint_controller) FetchAllComplaints(ctx *gin.Context) {
	var req dto.PaginationRequestParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	complaint_info, err := cont.comp_serv.FetchAllComplaints(req)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse("complaints not found")
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildResponseWithPagination(utils.FETCHED_SUCCESS, "", utils.COMPLAINT_DATA, complaint_info)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *complaint_controller) FetchComplaintDetailByComplaint(ctx *gin.Context) {
	complaint_id := ctx.Request.URL.Query().Get("complaint_id")

	if complaint_id == "" {
		cont.response = utils.BuildFailedResponse("please provide a required params")
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	obj_id, err := uuid.Parse(complaint_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	utils.REQUEST_HOST = ctx.Request.Host

	complaint, err := cont.comp_serv.FetchComplaintDetailByComplaint(obj_id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse("conplaint not founds")
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, complaint)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *complaint_controller) FetchDeviceImageURL(ctx *gin.Context) {

	var req dto.VideoRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	respo := generateSignedS3URL(req.FilePath, req.Directory)
	defer respo.Body.Close()
	ctx.DataFromReader(http.StatusOK, *respo.ContentLength, *respo.ContentType, respo.Body, nil)
}

func (cont *complaint_controller) FetchPDFFile(ctx *gin.Context) {

	var req dto.VideoRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}
	// ctx.JSON(http.StatusOK, "run.."+req.Directory+req.FilePath)
	respo := generateSignedS3URLForPDF(req.FilePath, req.Directory)
	//defer respo.Body.Close()
	ctx.DataFromReader(http.StatusOK, *respo.ContentLength, *respo.ContentType, respo.Body, nil)
}

func (cont *complaint_controller) UploadDeviceImage(ctx *gin.Context) {
	file, _, err := ctx.Request.FormFile("device_image")
	complaint_info_id := ctx.PostForm("complaint_info_id")

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	if complaint_info_id == "" {
		cont.response = utils.BuildFailedResponse("invalid comaplaint id")
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	if ctx.Request.ContentLength > 5*1024*1024 {
		cont.response = utils.BuildFailedResponse("image should be less that 5 MB")
		ctx.JSON(http.StatusRequestEntityTooLarge, cont.response)
		return
	}

	tempFile, err := ioutil.TempFile("media", "upload-*.png")

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	defer tempFile.Close()

	fileBytes, fileReader := ioutil.ReadAll(file)

	if fileReader != nil {
		cont.response = utils.BuildFailedResponse(fileReader.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	tempFile.Write(fileBytes)
	defer file.Close()
	defer tempFile.Close()
	_ = path.Base(tempFile.Name())

	err = cont.comp_serv.UploadDeviceImage(tempFile.Name(), complaint_info_id)
	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())

		if uuid.IsInvalidLengthError(err) {
			cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else if strings.Contains(err.Error(), "completed complaint") {
			ctx.JSON(http.StatusUnprocessableEntity, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse("File upload successfully", utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *complaint_controller) UploadDeviceVideo(ctx *gin.Context) {
	file, handler, err := ctx.Request.FormFile("device_video")
	complaint_info_id := ctx.PostForm("complaint_info_id")

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	if complaint_info_id == "" {
		cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err = cont.comp_serv.UploadDeviceVideo(file, handler, complaint_info_id)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if uuid.IsInvalidLengthError(err) {
			cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else if strings.Contains(err.Error(), "completed complaint") {
			ctx.JSON(http.StatusUnprocessableEntity, cont.response)
			return
		} else if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse("complaint not found to upload video")
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse("Video has been upload successfully", utils.COMPLAINT_DATA, handler.Filename)
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *complaint_controller) UpdateComplaintInfo(ctx *gin.Context) {
	var req dto.UpdateComplaintRequestDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	result, err := cont.comp_serv.UpdateComplaintInfo(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		if uuid.IsInvalidLengthError(err) {
			cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
			ctx.JSON(http.StatusConflict, cont.response)
			return
		} else if strings.Contains(err.Error(), "completed complaint") {
			ctx.JSON(http.StatusUnprocessableEntity, cont.response)
			return
		}
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.UPDATE_SUCCESS, utils.COMPLAINT_DATA, result)
	ctx.JSON(http.StatusOK, cont.response)
}

func generateSignedS3URL(img string, folder_name string) *s3.GetObjectOutput {
	s3_connection := connections.NewS3Connection()
	sess, err := s3_connection.MakeNewSession()
	if err != nil {
		return &s3.GetObjectOutput{}
	}

	bucket := s3_connection.GetBucketName()
	item := "/" + folder_name + "/" + img
	svc := s3.New(sess)

	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func generateSignedS3URLForPDF(img string, folder_name string) *s3.GetObjectOutput {
	s3_connection := connections.NewS3Connection()
	sess, err := s3_connection.MakeNewSession()
	if err != nil {
		return &s3.GetObjectOutput{}
	}

	bucket := s3_connection.GetBucketName()
	item := "/" + folder_name + "/" + img
	svc := s3.New(sess)
	content := "application/pdf"
	ContentDisposition := "inline"

	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket:                     aws.String(bucket),
		Key:                        aws.String(item),
		ResponseContentType:        &content,
		ResponseContentDisposition: &ContentDisposition,
	})
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func (cont *complaint_controller) DeleteDeviceFiles(ctx *gin.Context) {
	device_file_id := ctx.Param("file_id")

	err := cont.comp_serv.DeleteDeviceFiles(device_file_id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse("device file not found")
			ctx.JSON(http.StatusNotFound, cont.response)
			return

		} else if uuid.IsInvalidLengthError(err) {
			cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
			ctx.JSON(http.StatusConflict, cont.response)
			return
		}
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.DELETE_SUCCESS, utils.COMPLAINT_DATA, device_file_id)
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *complaint_controller) DeleteComplaint(ctx *gin.Context) {
	complaint_id := ctx.Param("complaint_id")

	if complaint_id == "" {
		cont.response = utils.BuildFailedResponse(utils.REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.comp_serv.DeleteComplaint(complaint_id)

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

	cont.response = utils.BuildSuccessResponse(utils.DELETE_SUCCESS, utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *complaint_controller) ComplaintResolve(ctx *gin.Context) {
	complaint_id := ctx.Param("complaint_id")

	if complaint_id == "" {
		cont.response = utils.RequestParamsMissingResponse(utils.REQUIRED_PARAMS)
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.comp_serv.ComplaintResolve(complaint_id)

	if err != nil {
		if uuid.IsInvalidLengthError(err) {
			cont.response = utils.BuildFailedResponse(utils.INVALID_PARAMS)
			ctx.JSON(http.StatusConflict, cont.response)
			return
		}

		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse("complaint has been resolved", utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, cont.response)
}

func (cont *complaint_controller) FetchAllComplaintCounts(ctx *gin.Context) {
	result := cont.comp_serv.FetchAllComplaintCounts()

	cont.response = utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, result)
	ctx.JSON(http.StatusOK, result)
}

func (cont *complaint_controller) ClientRegistration(ctx *gin.Context) {
	var req dto.ClientRegistration

	if err := ctx.ShouldBindJSON(&req); err != nil {
		cont.response = utils.BuildFailedResponse(helper.Handle_required_param_error(err))
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.comp_serv.ClientRegistration(req)

	if err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse(utils.USER_REGISTRATION_SUCCESS, utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusCreated, cont.response)
}

func (cont *complaint_controller) FetchComplaintsByClient(ctx *gin.Context) {
	var req dto.FetchComplaintsByClientRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	complaint_info, err := cont.comp_serv.FetchComplaintsByClient(req)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse("complaints not found")
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}

	cont.response = utils.BuildResponseWithPagination(utils.FETCHED_SUCCESS, "", utils.COMPLAINT_DATA, complaint_info)
	ctx.JSON(http.StatusOK, cont.response)
}

/* delete a client means simply deactivate a account */
func (cont *complaint_controller) DeleteClient(ctx *gin.Context) {
	client_id := ctx.Param("client_id")

	if client_id == "" {
		cont.response = utils.RequestParamsMissingResponse("provide a required params")
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	err := cont.comp_serv.DeleteClient(client_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			cont.response = utils.BuildFailedResponse("user not found to delete")
			ctx.JSON(http.StatusNotFound, cont.response)
			return
		}

		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusInternalServerError, cont.response)
		return
	}
	cont.response = utils.BuildSuccessResponse(utils.DELETE_SUCCESS, utils.COMPLAINT_DATA, nil)
	ctx.JSON(http.StatusOK, cont.response)
}
