package response

import (
	"errors"
	"fmt"
)

var (
	ErrBadRequest          = errors.New("Bad request")
	ErrForbiddenResource   = errors.New("Forbidden resource")
	ErrNotFound            = errors.New("Not Found")
	ErrPreConditionFailed  = errors.New("Precondition failed")
	ErrInternalServerError = errors.New("Internal server error")
	ErrTimeoutError        = errors.New("Timeout error")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrConflict            = errors.New("Conflict")
)

const (
	StatusCodeGenericSuccess            = "200000"
	StatusCodeBadRequest                = "400000"
	StatusCodeAlreadyRegistered         = "400001"
	StatusCodeUnauthorized              = "401000"
	StatusCodeForbidden                 = "403000"
	StatusCodeNotFound                  = "404000"
	StatusCodeConflict                  = "409000"
	StatusCodeGenericPreconditionFailed = "412000"
	StatusCodeOTPLimitReached           = "412550"
	StatusCodeNoLinkerExist             = "412553"
	StatusCodeInternalError             = "500000"
	StatusCodeFailedSellBatch           = "500100"
	StatusCodeFailedOTP                 = "503000"
	StatusCodeServiceUnavailable        = "503000"
	StatusCodeTimeoutError              = "504000"
)

func GetErrorCode(err error) string {
	err = getErrType(err)

	switch err {
	case ErrBadRequest:
		return StatusCodeBadRequest
	case ErrForbiddenResource:
		return StatusCodeForbidden
	case ErrNotFound:
		return StatusCodeNotFound
	case ErrConflict:
		return StatusCodeConflict
	case ErrUnauthorized:
		return StatusCodeUnauthorized
	case ErrForbiddenResource:
		return StatusCodeForbidden
	case ErrPreConditionFailed:
		return StatusCodeGenericPreconditionFailed
	case ErrInternalServerError:
		return StatusCodeInternalError
	case ErrTimeoutError:
		return StatusCodeTimeoutError
	case nil:
		return StatusCodeGenericSuccess
	default:
		return StatusCodeInternalError
	}
}

type ErrChain struct {
	Message string
	Cause   error
	Fields  map[string]interface{}
	Type    error
}

func (err ErrChain) Error() string {
	bcoz := ""
	fields := ""
	if err.Cause != nil {
		bcoz = fmt.Sprint(" because {", err.Cause.Error(), "}")
		if len(err.Fields) > 0 {
			fields = fmt.Sprintf(" with Fields {%+v}", err.Fields)
		}
	}
	return fmt.Sprint(err.Message, bcoz, fields)
}

func getErrType(err error) error {
	switch err.(type) {
	case ErrChain:
		errType := err.(ErrChain).Type
		if errType != nil {
			err = errType
		}
	}
	return err
}
