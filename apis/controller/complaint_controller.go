package controller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/services"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

type ComplaintController interface {
	CreateComplaint(*gin.Context)
	FetchAllComplaints(*gin.Context)
	FetchDummyURL(ctx *gin.Context)
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

	_, err := cont.comp_serv.FetchAllComplaints(req)

	if err != nil {
		response := utils.BuildFailedResponse(err.Error())
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	// respo := generateSignedS3URL()
	// fmt.Println("Content Type : ", *respo.ContentType)
	// ctx.DataFromReader(http.StatusOK, *respo.ContentLength, *respo.ContentType, respo.Body, nil)
	path := fmt.Sprintf("%s/%s", "media", "test.png")
	r_url := "http://localhost:8181/api/dummy/" + path

	response := utils.BuildSuccessResponse(utils.FETCHED_SUCCESS, utils.COMPLAINT_DATA, r_url)
	ctx.JSON(http.StatusOK, response)
}

func (cont *complaint_controller) FetchDummyURL(ctx *gin.Context) {
	var req dto.ImageRequestDTO

	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image Path", req.ImagePath)

	respo := generateSignedS3URL(req.ImagePath)
	fmt.Println("Content Type : ", *respo.ContentType)
	defer respo.Body.Close()
	ctx.DataFromReader(http.StatusOK, *respo.ContentLength, *respo.ContentType, respo.Body, nil)

}

var aws_access_key_id = "AKIA3VMV3LWIQ6EL63WU"
var aws_secret_access_key = "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"
var region = "ap-south-1"

var bucket_name = "skromansupportbucket"

func generateSignedS3URL(img string) *s3.GetObjectOutput {
	bucket := bucket_name
	item := "/media/" + img

	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, "")
	_, err := creds.Get()
	if err != nil {
		log.Fatal(err)
	}

	cfg := aws.NewConfig().WithRegion("ap-south-1").WithCredentials(creds)
	sess, _ := session.NewSession(cfg)
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
