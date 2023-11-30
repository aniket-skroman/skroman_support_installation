package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
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
	UpdateComplaintInfo(req dto.UpdateComplaintRequestDTO) (dto.ComplaintInfoDTO, error)
	DeleteDeviceFiles(file_id string) error
	DeleteComplaint(complaint_id string) error
	ComplaintResolve(complaint_id string) error
	FetchAllComplaintCounts() dto.AllComplaintsCount
	ClientRegistration(req dto.ClientRegistration) error
	FetchComplaintsByClient(req dto.FetchComplaintsByClientRequestDTO) (dto.ComplaintByClientDTO, error)
	Fetch_user_info(string, string, string) (dto.AllocatedEmpDetailsDTO, error)
	Fetch_client_info(client_id string) (interface{}, error)
	DeleteClient(client_id string) error
}

type complaint_service struct {
	complaint_repo  repositories.ComplaintRepository
	jwt_service     JWTService
	allocation_serv ComplaintAllocationService
}

func NewComplaintService(repo repositories.ComplaintRepository, serv JWTService, allocation_serv ComplaintAllocationService) ComplaintService {
	return &complaint_service{
		complaint_repo:  repo,
		jwt_service:     serv,
		allocation_serv: allocation_serv,
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
	avalibale_date, err := time.Parse("2006-01-02", req.ClientAvailableDate)
	time_slots := fmt.Sprintf("%s-%s", req.ClientTimeSlots.From, req.ClientTimeSlots.To)
	if err != nil {
		return nil, err
	}
	info_args := db.CreateComplaintInfoParams{
		ComplaintID:             complaint.ID,
		DeviceID:                req.DeviceID,
		DeviceType:              sql.NullString{String: req.DeviceType, Valid: true},
		DeviceModel:             sql.NullString{String: req.DeviceModel, Valid: true},
		ProblemStatement:        req.ProblemStatement,
		ProblemCategory:         sql.NullString{String: req.ProblemCategory, Valid: true},
		ClientAvailable:         time.Now(),
		Status:                  "INIT",
		ClientAvailableDate:     sql.NullTime{Time: avalibale_date, Valid: true},
		ClientAvailableTimeSlot: sql.NullString{String: time_slots, Valid: true},
	}

	complaint_info, err := ser.complaint_repo.CreateComplaintInfo(info_args)

	complaint_info_dto := dto.ComplaintInfoDTO{
		ID:                  complaint_info.ID,
		ComplaintID:         complaint_info.ComplaintID,
		DeviceID:            complaint_info.DeviceID,
		ProblemStatement:    complaint_info.ProblemStatement,
		ProblemCategory:     complaint_info.ProblemCategory.String,
		ClientAvailableDate: req.ClientAvailableDate,
		ClientTimeSlots:     req.ClientTimeSlots.From + "-" + req.ClientTimeSlots.To,
		Status:              complaint_info.Status,
		CreatedAt:           complaint_info.CreatedAt,
		UpdatedAt:           complaint_info.UpdatedAt,
		DeviceType:          complaint_info.DeviceType.String,
		DeviceModel:         complaint_info.DeviceModel.String,
	}

	return complaint_info_dto, err
}

func (ser *complaint_service) FetchAllComplaints(req dto.PaginationRequestParams) ([]dto.ComplaintInfoDTO, error) {
	args := db.FetchAllComplaintsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Status: req.TagKey,
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		count, err := ser.count_complaints(req.TagKey)

		if err != nil {
			helper.SetPaginationData(int(req.PageID), 0)
		} else {
			helper.SetPaginationData(int(req.PageID), count)
		}
	}()
	var complaints []db.ComplaintInfo
	var err error

	if req.TagKey == "TOTAL" {
		args := db.FetchTotalComplaintsParams{
			Limit:  req.PageSize,
			Offset: (req.PageID - 1) * req.PageSize,
		}
		complaints, err = ser.complaint_repo.FetchTotalComplaints(args)
	} else {
		complaints, err = ser.complaint_repo.FetchAllComplaints(args)
	}

	if err != nil {
		return nil, err
	}

	if len(complaints) == 0 {
		return nil, sql.ErrNoRows
	}

	complaint_info := new(dto.ComplaintInfoDTO).SetComplaintInfoData(complaints...)
	wg.Wait()

	if _, ok := complaint_info.([]dto.ComplaintInfoDTO); ok {
		return complaint_info.([]dto.ComplaintInfoDTO), nil
	}

	s_complaint := complaint_info.(dto.ComplaintInfoDTO)

	return []dto.ComplaintInfoDTO{s_complaint}, nil
}

func (ser *complaint_service) FetchComplaintDetailByComplaint(complaint_id uuid.UUID) (dto.ComplaintInfoByComplaintDTO, error) {
	complaint_info, err := ser.complaint_repo.FetchComplaintDetailByComplaint(complaint_id)

	if err != nil {
		return dto.ComplaintInfoByComplaintDTO{}, err
	}

	wg := sync.WaitGroup{}
	wg.Add(4)
	result := dto.ComplaintInfoByComplaintDTO{}

	// fetch client info by proxy call
	go func(client_id string) {
		defer wg.Done()

		client_info, err := ser.Fetch_client_info(complaint_info.Client.(string))
		if err != nil {
			result.ClientInfo = nil
		} else {
			result.ClientInfo = client_info
		}
	}(complaint_info.Client.(string))

	// fetch all device images and videos for complaint
	go func(complaint_info_id string) {
		defer wg.Done()
		device_images, device_videos, err := ser.FetchDeviceImagesByComplaintId(complaint_info.ComplaintInfoID.String())
		if err != nil {
			result.ComplaintInfo.DeviceImages = []dto.ComplaintDeviceImagesDTO{}
			result.ComplaintInfo.DeviceVideos = []dto.ComplaintDeviceImagesDTO{}
		} else {
			result.ComplaintInfo.DeviceImages = device_images
			result.ComplaintInfo.DeviceVideos = device_videos
		}
	}(complaint_info.ComplaintInfoID.String())

	go func() {
		defer wg.Done()
		user_data, _ := ser.Fetch_user_info(complaint_info.CreatedBy.String(), "EMP", "SALES")
		if user_data.FullName == "" {
			result.ComplaintInfo.CreatedBy = "NOT AVAILABEL"
		} else {
			result.ComplaintInfo.CreatedBy = user_data.FullName

		}
	}()

	// fetch allocation details
	go func() {
		defer wg.Done()
		allocate_data, _ := ser.fetch_allocated_emp_details(complaint_id.String())
		result.ComplaintInfo.AllocatedEmpDetailsDTO = allocate_data
	}()

	available_date := strings.ReplaceAll(complaint_info.ClientAvailableDate.Time.String(), "00:00:00 +0000 UTC", "")

	result.ComplaintInfo = dto.ComplaintFullDetailsDTO{
		Id:                  complaint_info.ComplaintInfoID,
		ComplaintId:         complaint_id,
		Client:              complaint_info.Client.(string),
		DeviceID:            complaint_info.DeviceID,
		DeviceModel:         complaint_info.DeviceModel.String,
		DeviceType:          complaint_info.DeviceType.String,
		ProblemStatement:    complaint_info.ProblemStatement,
		ProblemCategory:     complaint_info.ProblemCategory.String,
		ComplaintStatus:     complaint_info.ComplaintStatus,
		ComplaintRaisedAt:   complaint_info.ComplaintRaisedAt,
		LastModifiedAt:      complaint_info.LastModifiedAt,
		ClientAvailableDate: available_date,
		ClientTimeSlots:     complaint_info.ClientAvailableTimeSlot.String,
		ComplaintAddress:    complaint_info.ComplaintAddress.String,
	}

	wg.Wait()

	return result, nil
}

func (ser *complaint_service) FetchDeviceImagesByComplaintId(complaint_info_id string) ([]dto.ComplaintDeviceImagesDTO, []dto.ComplaintDeviceImagesDTO, error) {
	obj_id, _ := uuid.Parse(complaint_info_id)

	result, err := ser.complaint_repo.FetchDeviceImagesByComplaintId(obj_id)

	if err != nil {
		return nil, nil, err
	}

	if len(result) == 0 {
		return nil, nil, errors.New("device image/video not found")
	}

	device_images := []dto.ComplaintDeviceImagesDTO{}
	device_video := []dto.ComplaintDeviceImagesDTO{}
	for _, device := range result {
		temp := dto.ComplaintDeviceImagesDTO{
			ID:        device.ID,
			File:      "http://" + utils.REQUEST_HOST + "/api/device-file/" + device.DeviceImage,
			CreatedAt: device.CreatedAt,
			FileType:  device.FileType.String,
		}
		if device.FileType.String == "image/png" {
			device_images = append(device_images, temp)
		} else {
			device_video = append(device_video, temp)
		}

	}

	return device_images, device_video, nil
}

func (ser *complaint_service) UploadDeviceImage(file_path string, complaint_info_id string) error {

	complaint_obj_id, err := uuid.Parse(complaint_info_id)

	if err != nil {
		return err
	}

	// chech complaint is not COMPLETE, it should be in INIT or ALLOCATE
	status, err := ser.complaint_repo.FetchComplaintStatus(complaint_obj_id)

	if err != nil {
		return err
	}

	if status == "COMPLETE" {
		// if status is complete then remove the temp file
		ser.remove_local_files(file_path)
		return errors.New("completed complaint will not be updated")
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

	// chech complaint is not COMPLETE, it should be in INIT or ALLOCATE
	status, err := ser.complaint_repo.FetchComplaintStatus(complaint_obj_id)

	if err != nil {
		return err
	}

	if status == "COMPLETE" {
		return errors.New("completed complaint will not be updated")
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
	return err
}

func (ser *complaint_service) UpdateComplaintInfo(req dto.UpdateComplaintRequestDTO) (dto.ComplaintInfoDTO, error) {
	complaint_obj_id, err := uuid.Parse(req.ComplaintInfoId)

	if err != nil {
		return dto.ComplaintInfoDTO{}, err
	}

	// chech complaint is not COMPLETE, it should be in INIT or ALLOCATE
	status, err := ser.complaint_repo.FetchComplaintStatus(complaint_obj_id)

	if err != nil {
		return dto.ComplaintInfoDTO{}, err
	}

	if status == "COMPLETE" {
		return dto.ComplaintInfoDTO{}, errors.New("completed complaint will not be updated")
	}

	available_date, err := time.Parse("2006-01-02", req.ClientAvailableDate)

	if err != nil {
		return dto.ComplaintInfoDTO{}, err
	}

	args := db.UpdateComplaintInfoParams{
		ID:                      complaint_obj_id,
		DeviceID:                req.DeviceID,
		DeviceModel:             sql.NullString{String: req.DeviceModel, Valid: true},
		DeviceType:              sql.NullString{String: req.DeviceType, Valid: true},
		ProblemStatement:        req.ProblemStatement,
		ProblemCategory:         sql.NullString{String: req.ProblemCategory, Valid: true},
		ClientAvailableDate:     sql.NullTime{Time: available_date, Valid: true},
		ClientAvailableTimeSlot: sql.NullString{String: req.ClientTimeSlots.From + "-" + req.ClientTimeSlots.To, Valid: true},
	}

	result, err := ser.complaint_repo.UpdateComplaintInfo(args)
	err = helper.Handle_db_err(err)

	if err != nil {
		return dto.ComplaintInfoDTO{}, err
	}
	return dto.ComplaintInfoDTO{
		ID:                  result.ID,
		ComplaintID:         result.ComplaintID,
		DeviceID:            result.DeviceID,
		ProblemStatement:    result.ProblemStatement,
		ProblemCategory:     result.ProblemCategory.String,
		DeviceType:          result.DeviceType.String,
		DeviceModel:         result.DeviceModel.String,
		ClientAvailableDate: available_date.String(),
		Status:              result.Status,
		ClientTimeSlots:     result.ClientAvailableTimeSlot.String,
		ComplaintAddress:    result.ComplaintAddress.String,
		CreatedAt:           result.CreatedAt,
		UpdatedAt:           result.UpdatedAt,
	}, nil
}

func (ser *complaint_service) ComplaintResolve(complaint_id string) error {
	complanint_obj, err := helper.ValidateUUID(complaint_id)

	if err != nil {
		return err
	}

	args := db.UpdateComplaintStatusParams{
		ComplaintID: complanint_obj,
		Status:      "COMPLETE",
	}

	result, err := ser.complaint_repo.UpdateComplaintStatus(args)

	if err != nil {
		return err
	}

	affected_rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected_rows == 0 {
		return errors.New(utils.UPDATE_FAILED)
	}

	return nil
}

// delete a device files
func (ser *complaint_service) DeleteDeviceFiles(file_id string) error {
	file_obj_id, err := uuid.Parse(file_id)

	if err != nil {
		return err
	}

	// fetch device file details first
	device_file, err := ser.complaint_repo.FetchDeviceFileById(file_obj_id)

	if err != nil {
		return err
	}

	// then delete a refrance from db
	result, err := ser.complaint_repo.DeleteDeviceFiles(file_obj_id)

	if err != nil {
		return err
	}

	affected_rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected_rows == 0 {
		return errors.New("failed to delete device file")
	}

	// delete file from storage
	s3_connection := connections.NewS3Connection()
	err = s3_connection.DeleteFiles(device_file.DeviceImage)
	return err
}

// delete a complaint
func (ser *complaint_service) DeleteComplaint(complaint_id string) error {
	complaint_obj, err := helper.ValidateUUID(complaint_id)

	if err != nil {
		return err
	}

	// check this complaint id is exists or not
	_, err = ser.FetchComplaintById(complaint_id)

	if err != nil {
		return err
	}

	device_images, err := ser.complaint_repo.DeleteComplaint(complaint_obj)

	if err != nil {
		return err
	}

	// make a s3 conncection and delete a file one by one
	s3_connection := connections.NewS3Connection()
	for i := range device_images {
		err := s3_connection.DeleteFiles(device_images[i].DeviceImage)
		if err != nil {
			return err
		}
	}

	return nil
}

// client registartion by automation server
func (ser *complaint_service) ClientRegistration(req dto.ClientRegistration) error {
	// implement a validations for data
	isValid, _ := helper.ValidateInputs(req.Contact)

	if !isValid.(bool) {
		return helper.Err_Invalid_Input
	}

	reqUrl := "userapi/user-registration"
	req_body := proxycalls.ClientRegistration{
		UserName:     req.UserName,
		EmailID:      req.Email,
		MobileNumber: req.Contact,
		Password:     "skroman123",
		Address1:     req.AddressLine,
		City:         req.City,
		State:        req.State,
		PinCode:      req.Pincode,
	}

	var err error

	request_body, err := json.Marshal(&req_body)

	if err != nil {
		return err
	}

	api_call := proxycalls.ProxyCalls{}
	api_call.ReqEndpoint = reqUrl
	api_call.RequestBody = request_body
	api_call.RequestMethod = http.MethodPost

	response, err := api_call.MakeRequestWithBody()

	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		body_data := struct {
			Msg          string `json:"msg"`
			UserID       string `json:"userId"`
			EmailID      string `json:"emailId"`
			MobileNumber string `json:"mobileNumber"`
		}{}

		err = json.NewDecoder(response.Body).Decode(&body_data)
		return err
	default:
		err = errors.New("faild to created user account")
	}

	return err
}

func (ser *complaint_service) remove_local_files(file_path string) {
	_ = os.Remove(file_path)
}

func (ser *complaint_service) count_complaints(status string) (int64, error) {

	if status == "TOTAL" {
		count, err := ser.complaint_repo.TotalComplaints()
		if err != nil {
			return 0, err
		}
		return count, nil
	}
	result, err := ser.complaint_repo.CountComplaints(status)

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
func (ser *complaint_service) Fetch_client_info(client_id string) (interface{}, error) {
	token := ser.jwt_service.GenerateTempToken(utils.TOKEN_ID, "EMP", "SALES")
	api_request := proxycalls.NewAPIRequest("client/"+client_id, http.MethodGet, false, nil, nil, map[string]string{"Authorization": token})

	response, err := api_request.MakeApiRequest()

	if err != nil {
		return proxycalls.ClientByIdResponse{}, err
	}

	if response.StatusCode == http.StatusOK {

		response_data := struct {
			Error    string `json:"error"`
			Message  string `json:"message"`
			Status   bool   `json:"status"`
			UserData struct {
				ID        string    `json:"id"`
				UserName  string    `json:"user_name"`
				Email     string    `json:"email"`
				Contact   string    `json:"contact"`
				Address   string    `json:"address"`
				City      string    `json:"city"`
				State     string    `json:"state"`
				Pincode   string    `json:"pinCode"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
			} `json:"user_data"`
		}{}

		err = json.NewDecoder(response.Body).Decode(&response_data)
		return response_data.UserData, err
	}

	return proxycalls.ClientByIdResponse{}, nil
}

// fetch user data by user id from user-service
func (ser *complaint_service) Fetch_user_info(user_id, user_type, dept string) (dto.AllocatedEmpDetailsDTO, error) {

	// generate auth_token for user id
	token := ser.jwt_service.GenerateTempToken(user_id, user_type, dept)
	requrl := "http://15.207.19.172:8080/api/fetch-user"

	request, _ := http.NewRequest(http.MethodGet, requrl, nil)
	request.Close = true
	request.Header.Set("Authorization", token)

	response, err := http.DefaultClient.Do(request)

	defer func() {
		if err := response.Body.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return dto.AllocatedEmpDetailsDTO{}, err
	}

	if response.StatusCode == http.StatusOK {

		type AutoGenerated struct {
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
		}

		user_data := new(AutoGenerated)
		err = json.NewDecoder(response.Body).Decode(&user_data)

		if err != nil {
			return dto.AllocatedEmpDetailsDTO{}, err
		}
		return dto.AllocatedEmpDetailsDTO{
			AllocateUserID: user_data.UserData.ID,
			FullName:       user_data.UserData.FullName,
			Email:          user_data.UserData.Email,
			Contact:        user_data.UserData.Contact,
			UserType:       user_data.UserData.UserType,
		}, nil
	}

	return dto.AllocatedEmpDetailsDTO{}, nil
}

// fethc complaint allocated emp details
func (ser *complaint_service) fetch_allocated_emp_details(compaint_id string) (dto.AllocatedEmpDetailsDTO, error) {
	// fetch a allocation by complaint id
	allocation_data, err := ser.allocation_serv.FetchAllocationByComplaintId(compaint_id)
	if err != nil {
		return dto.AllocatedEmpDetailsDTO{}, err
	}

	// call user service to fetch a emp data
	user_data, err := ser.Fetch_user_info(allocation_data.AllocatedTo.String(), "EMP", "SALES")
	if err != nil {
		return dto.AllocatedEmpDetailsDTO{}, err
	}

	user_data.ID = allocation_data.ID.String()
	if !allocation_data.CreatedAt.IsZero() && !allocation_data.UpdatedAt.IsZero() {
		user_data.CreatedAt = allocation_data.CreatedAt
		user_data.UpdatedAt = allocation_data.UpdatedAt
	}

	return user_data, nil
}

func (ser *complaint_service) FetchComplaintById(complaint_id string) (db.FetchComplaintByComplaintIdRow, error) {
	complaint_obj, err := helper.ValidateUUID(complaint_id)

	if err != nil {
		return db.FetchComplaintByComplaintIdRow{}, err
	}

	return ser.complaint_repo.FetchComplaintById(complaint_obj)
}

func (ser *complaint_service) FetchAllComplaintCounts() dto.AllComplaintsCount {
	result := dto.AllComplaintsCount{}

	wg := sync.WaitGroup{}
	wg.Add(3)

	// fetch a all counts
	go func() {
		defer wg.Done()
		data, err := ser.fetch_all_complaint_counts()
		if err == nil {
			result.AllComplaints = data.AllComplaints
			result.AllocatedComplaints = data.AllocatedComplaints
			result.CompletedComplaints = data.ComletedComplaints
			result.PendingComplaints = data.PendingComplaints
		}
	}()

	go func() {
		defer wg.Done()
		month_data, _ := ser.fetch_month_wise_count()
		result.MonthWiseCounts = make([]dto.MonthWiseCounts, len(month_data))
		// if month availabel
		for i, data := range month_data {
			month := strings.Trim(data.Month, " ")
			result.MonthWiseCounts[i] = dto.MonthWiseCounts{
				Month: month,
				Count: data.Count,
			}
		}
	}()

	// fetching emp count
	go func() {
		defer wg.Done()
		count, _ := ser.fetch_clinet_count()
		result.EmpCount = count
	}()

	wg.Wait()

	return result
}

func (ser *complaint_service) fetch_all_complaint_counts() (db.CountAllComplaintRow, error) {
	return ser.complaint_repo.AllComplaintsCount()
}

func (ser *complaint_service) fetch_month_wise_count() ([]db.FetchCountByMonthRow, error) {
	return ser.complaint_repo.FetchCountByMonths()
}

func (ser *complaint_service) fetch_clinet_count() (int64, error) {
	token := ser.jwt_service.GenerateTempToken(utils.TOKEN_ID, "EMP", "SALES")
	api_request := proxycalls.NewAPIRequest("count-clients", http.MethodGet, false, nil, nil, map[string]string{"Authorization": token})

	response, err := api_request.MakeApiRequest()

	if err != nil {
		return 0, err
	}

	if response.StatusCode == http.StatusOK {
		response_data := struct {
			ClientCount int64 `json:"client_count"`
		}{}

		err = json.NewDecoder(response.Body).Decode(&response_data)
		return response_data.ClientCount, err
	}
	return 0, nil
}

// fetch complaints by clients
func (ser *complaint_service) FetchComplaintsByClient(req dto.FetchComplaintsByClientRequestDTO) (dto.ComplaintByClientDTO, error) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	err_chan := make(chan error)
	var result dto.ComplaintByClientDTO

	go func() {
		defer wg.Done()
		args := db.FetchComplaintsByClientParams{
			ClientID: req.ClientId,
			Limit:    req.PageSize,
			Offset:   (req.PageID - 1) * req.PageSize,
		}

		complaints, err := ser.complaint_repo.FetchComplaintsByClient(args)

		if err != nil {
			err_chan <- err
			return
		}

		complaint_info := make([]dto.ComplaintInfoDTO, len((complaints)))
		for i, data := range complaints {
			a_date := strings.ReplaceAll(data.ClientAvailableDate.Time.String(), "00:00:00 +0000 UTC", "")

			complaint_info[i] = dto.ComplaintInfoDTO{
				ID:                  data.ID_2,
				ComplaintID:         data.ID,
				DeviceID:            data.DeviceID,
				ProblemStatement:    data.ProblemStatement,
				ProblemCategory:     data.ProblemCategory.String,
				ClientAvailableDate: a_date,
				ClientTimeSlots:     data.ClientAvailableTimeSlot.String,
				ComplaintAddress:    data.ComplaintAddress.String,
				Status:              data.Status,
				CreatedAt:           data.CreatedAt,
				UpdatedAt:           data.UpdatedAt,
				DeviceType:          data.DeviceType.String,
				DeviceModel:         data.DeviceModel.String,
			}
		}

		result.ComplaintInfo = complaint_info
	}()

	go func() {
		defer wg.Done()
		count, _ := ser.complaint_repo.CountComplaintByClient(req.ClientId)
		helper.SetPaginationData(int(req.PageID), count)
	}()

	go func() {
		defer wg.Done()
		client_data, err := ser.Fetch_client_info(req.ClientId)

		if err != nil {
			err_chan <- err
			return
		}
		result.ClientInfo = client_data
	}()

	go func() {
		wg.Wait()
		close(err_chan)
	}()

	for data_err := range err_chan {
		return dto.ComplaintByClientDTO{}, data_err
	}

	return result, nil
}

/* delete a client by proxy call */
func (ser *complaint_service) DeleteClient(client_id string) error {
	client_data := struct {
		UserId string `json:"userId"`
	}{
		UserId: client_id,
	}

	request_body, err := json.Marshal(&client_data)

	if err != nil {
		return err
	}

	proxy := proxycalls.ProxyCalls{}
	proxy.IsRequestBody = true
	proxy.RequestBody = request_body
	proxy.RequestMethod = http.MethodPut
	proxy.ReqEndpoint = "profileapi/deactivate_account"

	response, err := proxy.MakeRequestWithBody()

	if err != nil {
		return err
	}

	defer func(response *http.Response) {
		if err := response.Body.Close(); err != nil {
			return
		}
	}(response)

	response_code := response.StatusCode
	if response_code == http.StatusOK {

		response_body := struct {
			Msg          string `json:"msg"`
			UserId       string `json:"userId"`
			EmailId      string `json:"emailId"`
			MobileNumber string `json:"mobileNumber"`
		}{}

		err := json.NewDecoder(response.Body).Decode(&response_body)
		if err != nil {
			return errors.New("something wents wrong")
		}

		if !strings.Contains(response_body.Msg, "User Account Deactivate Success") {
			return errors.New("failed to delete user please try again")
		}
		return nil

	} else if response_code == http.StatusNotFound {
		return sql.ErrNoRows
	}
	return errors.New("unknow error occured")
}
