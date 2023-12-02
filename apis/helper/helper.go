package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"regexp"

	"github.com/aniket-skroman/skroman_support_installation/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func SetPaginationData(page int, total int64) {
	if page == 0 {
		utils.PREVIOUS_IDX = 0
	} else {
		utils.PREVIOUS_IDX = page - 1
	}

	utils.CURRENT_IDX = page
	utils.TOTALCOUNT = total
}

type InputValidation interface {
	Validate() (interface{}, error)
}

type UUIDValidate struct {
	InputID string
}

func (uv *UUIDValidate) Validate() (interface{}, error) {
	data, err := uuid.Parse(uv.InputID)

	if err != nil {
		return nil, ERR_INVALID_ID
	}
	return data, nil
}

type ValidateContact struct {
	ContactNumber string
}

func (vm *ValidateContact) Validate() (interface{}, error) {
	regex := `^[0-9]{10}$`
	match, err := regexp.MatchString(regex, vm.ContactNumber)
	return match, err
}

func validate_data(input InputValidation) (interface{}, error) {
	return input.Validate()
}

func ValidateInputs(input any) (interface{}, error) {

	if data, ok := input.(uuid.UUID); ok {
		uuid_input := UUIDValidate{InputID: data.String()}

		var validate InputValidation = &uuid_input
		return validate_data(validate)
	}

	if data, ok := input.(string); ok {
		contact_val := ValidateContact{ContactNumber: data}
		var validate InputValidation = &contact_val
		return validate_data(validate)
	}

	return true, nil
}

func ValidateUUID(input_id string) (uuid.UUID, error) {
	return uuid.Parse(input_id)
}

type ApiError struct {
	Field string
	Msg   string
}

func Error_handler(err error) []ApiError {
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ApiError, len(ve))
			for i, fe := range ve {
				out[i] = ApiError{fe.Field(), msgForTag(fe.Tag())}
			}
			return out
		}
		return nil
	}
	return nil
}

func Handle_required_param_error(err error) string {
	var ve validator.ValidationErrors
	var err_msg string
	if errors.As(err, &ve) {
		for _, fe := range ve {
			err_msg = fmt.Sprintf("%v - %v", fe.Field(), msgForTag(fe.Tag()))
			break
		}
	}

	return err_msg
}

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	case "min":
		return "Invalid length for param"
	case "oneof":
		return "invalid param detected"
	default:
		return Err_Something_Wents_Wrong.Error()
	}
}

var key = "skroman-user-servi-12345"

func EncryptData(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptData(ciphertext string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
}
