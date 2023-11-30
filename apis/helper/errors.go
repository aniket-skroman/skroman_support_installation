package helper

import "errors"

var (
	ERR_INVALID_ID            error
	ERR_REQUIRED_PARAMS       error
	Err_Lead_Exists           error
	Err_Data_Not_Found        error
	Err_Update_Failed         error
	Err_Delete_Failed         error
	Err_Something_Wents_Wrong error
	Err_Invalid_Input         error
)

func init() {
	ERR_INVALID_ID = errors.New("invalid id found")
	ERR_REQUIRED_PARAMS = errors.New("please provide a required params")
	Err_Lead_Exists = errors.New("lead info already exists for current lead")
	Err_Data_Not_Found = errors.New("data not found")
	Err_Update_Failed = errors.New("failed to update resources")
	Err_Delete_Failed = errors.New("failed to delete resource")
	Err_Something_Wents_Wrong = errors.New("something wents wrong")
	Err_Invalid_Input = errors.New("invalid input data found")
}
