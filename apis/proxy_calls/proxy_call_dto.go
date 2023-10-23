package proxycalls

// response dto for fetch user by id

type ClientInfoDTO struct {
	UserID       string `json:"userId,omitempty"`
	EmailID      string `json:"emailId,omitempty"`
	Address1     string `json:"address1,omitempty"`
	Address2     string `json:"address2,omitempty"`
	City         string `json:"city,omitempty"`
	ImageUser    string `json:"imageUser,omitempty"`
	MobileNumber string `json:"mobileNumber,omitempty"`
	PinCode      string `json:"pinCode,omitempty"`
	State        string `json:"state,omitempty"`
	UserName     string `json:"userName,omitempty"`
}
type ClientByIdResponse struct {
	Msg    string        `json:"msg"`
	Result ClientInfoDTO `json:"result"`
}

type ClientRegistration struct {
	UserName     string `json:"userName"`
	EmailID      string `json:"emailId"`
	MobileNumber string `json:"mobileNumber"`
	Password     string `json:"password"`
	Address1     string `json:"address1"`
	City         string `json:"city"`
	State        string `json:"state"`
	PinCode      string `json:"pinCode"`
}
