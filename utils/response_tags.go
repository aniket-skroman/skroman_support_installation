package utils

var EmptyStr string = ""

type EmptyObj struct{}

const USER_REGISTRATION_SUCCESS = "User registration has been successful"
const INVALID_PARAMS = "Invalid params, can't process"

// REQUIRED_PARAMS --------------- DEVELOPER MESSAGES -----------------//
const REQUIRED_PARAMS = "please provide a required params."
const FETCHED_FAILED = "Failed to fetch the data from server."
const FETCHED_SUCCESS = "Data fetched successfully."

// DATA_INSERTED -------------- END USERS MESSAGES ----------------//
const DATA_INSERTED = "Operation successful."
const DATA_INSERTED_FAILED = "Operation failed, please try again."
const UPDATE_FAILED = "update failed please try again"
const UPDATE_SUCCESS = "Update successful."
const PASSWORD_RESET = "Password reset successful."
const FORGOT_PASSWORD = "Forgot password successful."
const LOGIN_SUCCESS = "Login successful."
const DATA_NOT_FOUND = "Data not found."
const DATA_FOUND = "Data found."

// FAILED_PROCESS -------------- COMMON MESSAGES ----------------//
var FAILED_PROCESS = "failed to process the request."

const PAGINATION_INVALID = "Pagination failed, please try again ..."
const DELETE_FAILED = "failed to delete"
const DELETE_SUCCESS = "successfully deleted"

// REQUEST_HOST -------------- BACKEND DEV VARIABLES -----------//S
var REQUEST_HOST = ""

var CURRENT_IDX = 0
var PREVIOUS_IDX = 0
var TOTALCOUNT int64 = 0
var TOKEN_ID = ""
var USER_TYPE = ""

// DATA --------------- Response Key Words ------------//
const DATA = "app_data"

var COMPLAINT_DATA = "complaint_data"

// ---------------------------- COMPLAINTS --------------------------------- //
const COMPLAINT_CREATED = "complaint has been created"
