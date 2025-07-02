package constants

import "net/http"

const (
	StatusCtxKey                = 0
	StatusSuccess               = http.StatusOK
	StatusErrorForm             = http.StatusBadRequest
	StatusErrorUnknown          = http.StatusBadGateway
	StatusInternalError         = http.StatusInternalServerError
	StatusUnauthorized          = http.StatusUnauthorized
	StatusCreated               = http.StatusCreated
	StatusAccepted              = http.StatusAccepted
	StatusForbidden             = http.StatusForbidden
	StatusInvalidAuthentication = http.StatusProxyAuthRequired
	StatusNotFound              = http.StatusNotFound
	StatusConflict              = http.StatusConflict
	StatusUnprocessableEntity   = http.StatusUnprocessableEntity
	StatusTooManyRequests       = http.StatusTooManyRequests
	StatusServiceUnavailable    = http.StatusServiceUnavailable
	StatusGatewayTimeout        = http.StatusGatewayTimeout
	StatusNotImplemented        = http.StatusNotImplemented
)

var statusMap = map[int][]string{
	StatusSuccess:               {"STATUS_OK", "Success"},
	StatusErrorForm:             {"STATUS_BAD_REQUEST", "Invalid data request"},
	StatusErrorUnknown:          {"STATUS_BAD_GATEWAY", "Oops something went wrong"},
	StatusInternalError:         {"INTERNAL_SERVER_ERROR", "Oops something went wrong"},
	StatusUnauthorized:          {"STATUS_UNAUTHORIZED", "Not authorized to access the service"},
	StatusCreated:               {"STATUS_CREATED", "Resource has been created"},
	StatusAccepted:              {"STATUS_ACCEPTED", "Resource has been accepted"},
	StatusForbidden:             {"STATUS_FORBIDDEN", "Forbidden access the resource "},
	StatusInvalidAuthentication: {"STATUS_INVALID_AUTHENTICATION", "The resource owner or authorization server denied the request"},
	StatusNotFound:              {"STATUS_NOT_FOUND", "Not Found"},
	StatusConflict:              {"STATUS_CONFLICT", "Conflict"},
	StatusUnprocessableEntity:   {"STATUS_UNPROCESSABLE_ENTITY", "Unprocessable Entity"},
	StatusTooManyRequests:       {"STATUS_TOO_MANY_REQUESTS", "Too Many Requests"},
	StatusServiceUnavailable:    {"STATUS_SERVICE_UNAVAILABLE", "Oops something went wrong"},
	StatusGatewayTimeout:        {"STATUS_GATEWAY_TIMEOUT", "Oops something went wrong"},
	StatusNotImplemented:        {"STATUS_NOT_IMPLEMENTED", "Oops something went wrong"},
}

func StatusCode(code int) string {
	return statusMap[code][0]
}

func StatusText(code int) string {
	return statusMap[code][1]
}

const (
	RepositoryPostgreSQL = "repository.postgresql"
	Usecase              = "usecase"
	Controller           = "controller"
)

const (
	RoleAdmin = 1
	RoleUser  = 2
)

const (
	TypeAttendance = 1
	TypeOvertime   = 2
)

const OvertimeFee = 100000
