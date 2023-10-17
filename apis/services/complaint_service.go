package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/aniket-skroman/skroman_support_installation/apis/dto"
	"github.com/aniket-skroman/skroman_support_installation/apis/helper"
	proxycalls "github.com/aniket-skroman/skroman_support_installation/apis/proxy_calls"
	"github.com/aniket-skroman/skroman_support_installation/apis/repositories"
	"github.com/aniket-skroman/skroman_support_installation/connections"
	db "github.com/aniket-skroman/skroman_support_installation/sqlc_lib"
	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/google/uuid"
)

type ComplaintService interface {
	CreateComplaint(dto.CreateComplaintRequestDTO) (interface{}, error)
	FetchAllComplaints(dto.PaginationRequestParams) ([]dto.ComplaintInfoDTO, error)
	FetchComplaintDetailByComplaint(uuid.UUID) (dto.ComplaintInfoByComplaintDTO, error)
	UploadDeviceImage(file_path string, complaint_info_id string) error
	UploadDeviceVideo(file multipart.File, handler *multipart.FileHeader, complaint_info_id string) error
}

type complaint_service struct {
	complaint_repo repositories.ComplaintRepository
	jwt_service    JWTService
}

func NewComplaintService(repo repositories.ComplaintRepository, serv JWTService) ComplaintService {
	return &complaint_service{
		complaint_repo: repo,
		jwt_service:    serv,
	}
}

func (ser *complaint_service) CreateComplaint(req dto.CreateComplaintRequestDTO) (interface{}, error) {
	//init a complaint first
	comp_args := db.CreateComplaintParams{
		ClientID:  req.ClientID,
		CreatedBy: uuid.New(),
	}
	complaint, err := ser.complaint_repo.CreateComplaint(comp_args)

	if err != nil {
		return nil, err
	}

	// create a complaint info
	c_time, err := time.Parse("2006-01-02 15:04:05", req.ClientAvailable)
	avalibale_date, err := time.Parse("2006-01-02", req.ClientAvailableDate)
	time_slots := fmt.Sprintf("%s %s", req.ClientTimeSlots.From, req.ClientTimeSlots.To)
	if err != nil {
		return nil, err
	}
	info_args := db.CreateComplaintInfoParams{
		ComplaintID: complaint.ID,
		DeviceID:    req.DeviceID,
		DeviceType: sql.NullString{
			String: req.DeviceModel, Valid: true,
		},
		DeviceModel:      sql.NullString{String: req.DeviceModel, Valid: true},
		ProblemStatement: req.ProblemStatement,
		ProblemCategory:  sql.NullString{String: req.ProblemCategory, Valid: true},
		ClientAvailable:  c_time,
		Status:           "INIT",
		ClientAvailableDate: sql.NullTime{
			Time:  avalibale_date,
			Valid: true,
		},
		ClientAvailableTimeSlot: sql.NullString{String: time_slots, Valid: true},
	}

	complaint_info, err := ser.complaint_repo.CreateComplaintInfo(info_args)

	complaint_info_dto := dto.ComplaintInfoDTO{
		ID:               complaint_info.ID,
		ComplaintID:      complaint_info.ComplaintID,
		DeviceID:         complaint_info.DeviceID,
		ProblemStatement: complaint_info.ProblemStatement,
		ProblemCategory:  complaint_info.ProblemCategory.String,
		ClientAvailable:  complaint_info.ClientAvailable,
		Status:           complaint_info.Status,
		CreatedAt:        complaint_info.CreatedAt,
		UpdatedAt:        complaint_info.UpdatedAt,
		DeviceType:       complaint_info.DeviceType.String,
		DeviceModel:      complaint_info.DeviceModel.String,
	}

	return complaint_info_dto, err
}

func (ser *complaint_service) FetchAllComplaints(req dto.PaginationRequestParams) ([]dto.ComplaintInfoDTO, error) {
	args := db.FetchAllComplaintsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		count, err := ser.count_complaints()

		if err != nil {
			utils.SetPaginationData(int(req.PageID), 0)
			return
		}

		utils.SetPaginationData(int(req.PageID), count)

	}()

	complaints, err := ser.complaint_repo.FetchAllComplaints(args)

	if err != nil {
		return nil, err
	}

	if len(complaints) == 0 {
		return nil, errors.New("complaint not found")
	}

	complaint_info := new(dto.ComplaintInfoDTO).SetComplaintInfoData(complaints...)

	if _, ok := complaint_info.([]dto.ComplaintInfoDTO); ok {
		return complaint_info.([]dto.ComplaintInfoDTO), nil
	}

	s_complaint := complaint_info.(dto.ComplaintInfoDTO)
	wg.Wait()
	return []dto.ComplaintInfoDTO{s_complaint}, nil
}

func (ser *complaint_service) FetchComplaintDetailByComplaint(complaint_id uuid.UUID) (dto.ComplaintInfoByComplaintDTO, error) {
	complaint_info, err := ser.complaint_repo.FetchComplaintDetailByComplaint(complaint_id)

	if err != nil {
		return dto.ComplaintInfoByComplaintDTO{}, err
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	result := dto.ComplaintInfoByComplaintDTO{}

	// fetch client info by proxy call
	go func(client_id string) {
		defer wg.Done()

		client_info, err := ser.fetch_client_info(complaint_info.Client)
		if err != nil {
			result.ClientInfo = proxycalls.ClientInfoDTO{}
		} else {
			result.ClientInfo = client_info.Result
		}
	}(complaint_info.Client)

	// fetch all device images and videos for complaint
	go func(complaint_info_id string) {
		defer wg.Done()
		device_images, err := ser.FetchDeviceImagesByComplaintId(complaint_info.ComplaintInfoID.String())
		if err != nil {
			result.ComplaintInfo.DeviceInfo = []dto.ComplaintDeviceImagesDTO{}
		} else {
			result.ComplaintInfo.DeviceInfo = device_images
		}
	}(complaint_info.ComplaintInfoID.String())

	go func() {
		defer wg.Done()
		user_name, _ := ser.fetch_user_info(complaint_info.CreatedBy.String())
		result.ComplaintInfo.CreatedBy = user_name
	}()

	result.ComplaintInfo = dto.ComplaintFullDetailsDTO{
		Client:            complaint_info.Client,
		DeviceID:          complaint_info.DeviceID,
		DeviceModel:       complaint_info.DeviceModel.String,
		DeviceType:        complaint_info.DeviceType.String,
		ProblemStatement:  complaint_info.ProblemStatement,
		ProblemCategory:   complaint_info.ProblemCategory.String,
		ComplaintStatus:   complaint_info.ComplaintStatus,
		ComplaintRaisedAt: complaint_info.ComplaintRaisedAt,
		LastModifiedAt:    complaint_info.LastModifiedAt,
		ClientAvailable:   complaint_info.ClientAvailable,
	}

	wg.Wait()

	return result, nil
}

func (ser *complaint_service) FetchDeviceImagesByComplaintId(complaint_info_id string) ([]dto.ComplaintDeviceImagesDTO, error) {
	obj_id, _ := uuid.Parse(complaint_info_id)

	result, err := ser.complaint_repo.FetchDeviceImagesByComplaintId(obj_id)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, errors.New("device image/video not found")
	}

	device_images := make([]dto.ComplaintDeviceImagesDTO, len(result))
	for i, device := range result {
		device_images[i] = dto.ComplaintDeviceImagesDTO{
			File:      "http://" + utils.REQUEST_HOST + "/api/device-image/" + device.DeviceImage,
			CreatedAt: device.CreatedAt,
			FileType:  device.FileType.String,
		}
	}

	return device_images, nil
}

func (ser *complaint_service) UploadDeviceImage(file_path string, complaint_info_id string) error {

	complaint_obj_id, err := uuid.Parse(complaint_info_id)

	if err != nil {
		return err
	}

	// upload a image in s3 bucket first
	s3_connection := connections.NewS3Connection()
	_, path, err := s3_connection.UploadDeviceImage(file_path)
	if err != nil {
		return err
	}

	if path == "" {
		return errors.New("failed to image process")
	}

	// remove the file from local storege
	ser.remove_local_files(file_path)

	// make args to store a image ref in db
	args := db.UploadDeviceImagesParams{
		DeviceImage:     path,
		ComplaintInfoID: complaint_obj_id,
		FileType:        sql.NullString{String: "image/png", Valid: true},
	}

	_, err = ser.complaint_repo.UploadDeviceImage(args)

	err = helper.Handle_db_err(err)

	return err
}

func (ser *complaint_service) UploadDeviceVideo(file multipart.File, handler *multipart.FileHeader, complaint_info_id string) error {

	complaint_obj_id, err := uuid.Parse(complaint_info_id)

	if err != nil {
		return err
	}

	// upload a image in s3 bucket first
	s3_connection := connections.NewS3Connection()
	path, err := s3_connection.UploadDeviceVideo(file, handler)
	if err != nil {
		return err
	}

	// make args to store a image ref in db
	args := db.UploadDeviceImagesParams{
		DeviceImage:     path,
		ComplaintInfoID: complaint_obj_id,
		FileType:        sql.NullString{String: "video/mp4", Valid: true},
	}

	_, err = ser.complaint_repo.UploadDeviceImage(args)

	err = helper.Handle_db_err(err)
	fmt.Println("File Path from service : ", path, complaint_obj_id)
	return err
}

func (ser *complaint_service) remove_local_files(file_path string) {
	_ = os.Remove(file_path)
}

func (ser *complaint_service) count_complaints() (int64, error) {
	result, err := ser.complaint_repo.CountComplaints()

	if err != nil {
		return 0, err
	}

	affected_rows, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return affected_rows, nil
}

// from another server
func (ser *complaint_service) fetch_client_info(client_id string) (proxycalls.ClientByIdResponse, error) {
	proxy_call := proxycalls.ProxyCalls{}
	proxy_call.ReqEndpoint = "profileapi/profileuser/userId"
	proxy_call.RequestMethod = http.MethodPost

	temp_body := struct {
		UserId string `json:"userId"`
	}{
		UserId: client_id,
	}

	request_body, err := json.Marshal(temp_body)

	if err != nil {
		return proxycalls.ClientByIdResponse{}, err
	}

	proxy_call.RequestBody = request_body

	response, err := proxy_call.MakeRequestWithBody()

	defer func() {
		if err := response.Body.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return proxycalls.ClientByIdResponse{}, err
	}

	response_body, err := io.ReadAll(response.Body)

	if err != nil {
		return proxycalls.ClientByIdResponse{}, err
	}

	client_info := proxycalls.ClientByIdResponse{}

	if err := json.Unmarshal(response_body, &client_info); err != nil {
		return proxycalls.ClientByIdResponse{}, err
	}

	if reflect.DeepEqual(client_info, proxycalls.ClientByIdResponse{}) || client_info.Msg != "success get the user Profile" {
		return proxycalls.ClientByIdResponse{}, errors.New("client info not found")
	}

	return client_info, nil
}

// fetch user data by user id from user-service
func (ser *complaint_service) fetch_user_info(user_id string) (string, error) {
	// generate auth_token for user id
	token := ser.jwt_service.GenerateToken(user_id, "EMP")

	requrl := "http://15.207.19.172:8080/api/fetch-user"

	request, err := http.NewRequest(http.MethodGet, requrl, nil)

	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", token)

	response, err := http.DefaultClient.Do(request)

	defer func() {
		if err := response.Body.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return "NOT AVAILABEL", nil
	}

	if response.StatusCode == http.StatusOK {
		response_body, err := io.ReadAll(response.Body)

		if err != nil {
			return "NOT AVAILABEL", nil
		}

		user_data := struct {
			Error    string `json:"error"`
			Message  string `json:"message"`
			Status   bool   `json:"status"`
			UserData struct {
				ID        string    `json:"id"`
				FullName  string    `json:"full_name"`
				Email     string    `json:"email"`
				Contact   string    `json:"contact"`
				UserType  string    `json:"user_type"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"user_data"`
		}{}

		json.Unmarshal(response_body, &user_data)
		if user_data.UserData.FullName != "" {
			return user_data.UserData.FullName, nil
		}

		return "NOT AVAILABEL", nil
	}

	return "NOT AVAILABEL", nil
}
