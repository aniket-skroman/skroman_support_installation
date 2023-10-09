package utils

func response_builder(status bool, msg, err, data_name *string, data *interface{}) (response map[string]interface{}) {
	response = make(map[string]interface{})

	response["status"] = status
	response["message"] = msg
	response["error"] = err
	response[*data_name] = data

	return
}

func BuildSuccessResponse(msg, data_name string, data interface{}) map[string]interface{} {
	return response_builder(true, &msg, &EmptyStr, &data_name, &data)
}

func BuildFailedResponse(msg, err, data_name string, data interface{}) map[string]interface{} {
	return response_builder(false, &FAILED_PROCESS, &err, &data_name, &data)
}
