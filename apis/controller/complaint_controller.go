package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
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
	ctx.JSON(http.StatusCreated, response)
}

func (cont *complaint_controller) FetchAllComplaints(ctx *gin.Context) {
	var req dto.PaginationRequestParams

	if err := ctx.ShouldBindUri(&req); err != nil {
		response := utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	complaint_info, err := cont.comp_serv.FetchAllComplaints(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.BuildResponseWithPagination(utils.FETCHED_SUCCESS, "", utils.COMPLAINT_DATA, complaint_info)
	ctx.JSON(http.StatusOK, response)
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
		if strings.Contains(err.Error(), "no rows in result set") {
			cont.response = utils.BuildFailedResponse("conplaint not found")
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
	var req dto.ImageRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		cont.response = utils.BuildFailedResponse(err.Error())
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}
	respo := generateSignedS3URL(req.ImagePath)
	defer respo.Body.Close()
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
		ctx.JSON(http.StatusBadRequest, cont.response)
		return
	}

	cont.response = utils.BuildSuccessResponse("File upload successfully", utils.COMPLAINT_DATA, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, cont.response)
}

func generateSignedS3URL(img string) *s3.GetObjectOutput {
	s3_connection := connections.NewS3Connection()
	sess, err := s3_connection.MakeNewSession()
	if err != nil {
		return &s3.GetObjectOutput{}
	}

	bucket := s3_connection.GetBucketName()
	item := "/media/" + img
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
