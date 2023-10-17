package utils

type Pagination struct {
	CurrentIdx  int   `json:"current_idx"`
	PreviousIdx int   `json:"previous_idx"`
	TotalCount  int64 `json:"total_count"`
}

func PaginationData() Pagination {
	return Pagination{
		CurrentIdx:  CURRENT_IDX,
		PreviousIdx: PREVIOUS_IDX,
		TotalCount:  TOTALCOUNT,
	}
}

func response_builder(status bool, msg, err, data_name *string, data *interface{}, isPagination bool) (response map[string]interface{}) {

	response = map[string]interface{}{}

	response["status"] = status
	response["message"] = msg
	response["error"] = err
	response[*data_name] = data
	if isPagination {
		var paginationData = PaginationData()

		response["pagination"] = paginationData
	}

	return
}

func BuildResponseWithPagination(msg, err, data_name string, data interface{}) map[string]interface{} {
	response := response_builder(true, &msg, &err, &data_name, &data, true)
	return response
}

func BuildSuccessResponse(msg, data_name string, data interface{}) map[string]interface{} {
	return response_builder(true, &msg, &EmptyStr, &data_name, &data, false)
}

func BuildFailedResponse(err string) map[string]interface{} {
	var data interface{}
	data = EmptyObj{}
	return response_builder(false, &FAILED_PROCESS, &err, &COMPLAINT_DATA, &data, false)
}

func SetPaginationData(page int, total int64) {
	if page == 0 {
		PREVIOUS_IDX = 0
	} else {
		PREVIOUS_IDX = page - 1
	}

	CURRENT_IDX = page
	TOTALCOUNT = total
}
