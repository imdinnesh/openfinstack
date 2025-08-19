package apierror


type APIError struct {
Code string `json:"code"`
Message string `json:"message"`
}


func (e APIError) Error() string { return e.Message }


func BadRequest(msg string) APIError { return APIError{Code: "BAD_REQUEST", Message: msg} }
func NotFound(msg string) APIError { return APIError{Code: "NOT_FOUND", Message: msg} }
func Conflict(msg string) APIError { return APIError{Code: "CONFLICT", Message: msg} }
func Internal(msg string) APIError { return APIError{Code: "INTERNAL", Message: msg} }