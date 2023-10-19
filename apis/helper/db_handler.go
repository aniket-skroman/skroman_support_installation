package helper

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func Handle_db_err(err error) (err_ error) {
	switch e := err.(type) {
	case *pq.Error:
		fmt.Println("DB Error Code : ", e.Code)
		switch e.Code {
		case "23502":
			// not-null constraint violation
			fmt.Println("Some required data was left out:\n\n", e.Message)
			err_ = errors.New(e.Detail)
			return
		case "23505":
			// unique constraint violation
			if strings.Contains(e.Message, "full_name") {
				err_ = errors.New("user account already exists")
				return
			}
			err_ = errors.New(e.Detail)
			return

		case "23514":
			fmt.Println("Handle_DBError called from constraint check")

			// check constraint violation
			if strings.Contains(e.Message, "contact") {
				err_ = errors.New("contact should not be empty")
				return
			} else if strings.Contains(e.Message, "email") {
				err_ = errors.New("email should not be empty")
				return
			}
			// err_ = validate_err_msg(&e.Message)
			// return
		case "23503":
			err_ = errors.New("invalid id has been provided,please try with valid id's")
			return
		case "2201X":
			// when offset, limit not working, for pagination data
			err_ = errors.New("invalid pagination found")
			return
		default:
			msg := e.Message
			if d := e.Detail; d != "" {
				msg += "\n\n" + d
			}
			if h := e.Hint; h != "" {
				msg += "\n\n" + h
			}
			fmt.Println("Message from default : ", e.Code)
			err_ = errors.New(msg)
			return
		}
	default:
		err_ = nil
		return

	}

	return
}
