package domain

type ResponseError interface {
	Message() string
	ErrorType() ResponseErrorType
	ResponseError()
}

type responseErrorImp struct {
	errorType ResponseErrorType
	message   string
}

func (r *responseErrorImp) ResponseError() {
	// Noncompliant
}

func (r *responseErrorImp) Message() string {
	return r.message
}

func (r *responseErrorImp) ErrorType() ResponseErrorType {
	return r.errorType
}

func NewReponseError(ErrorType ResponseErrorType, message string) ResponseError {
	if message == "" {
		message = string(ErrorType)
	}

	return &responseErrorImp{
		errorType: ErrorType,
		message:   message,
	}
}

type ResponseErrorType string

const (
	// 400 Bad Request
	ErrorInvalidArgument   ResponseErrorType = "invalid argument"
	ErrorFaildPrecondition ResponseErrorType = "failed precondition"
	ErrorOutOfRange        ResponseErrorType = "out of range"
	// 401 Unauthorized
	ErrorUnauthenticated ResponseErrorType = "unauthenticated"
	// 403 Forbidden
	ErrorPermissionDenied ResponseErrorType = "permission denied"
	// 404 Not Found
	ErrorNotFound ResponseErrorType = "not found"
	// 409 Conflict
	ErrorAlreadyExists ResponseErrorType = "already exists"
	ErrorAborted       ResponseErrorType = "aborted"
	// 429 Too Many Request
	ErrorResourceExhausted ResponseErrorType = "resource exhausted"
	// 499 Client Closed Request
	ErrorCancelled ResponseErrorType = "cancelled"

	// 500 Internal Server Error
	ErrorUnknown  ResponseErrorType = "unknown Erroror"
	ErrorInternal ResponseErrorType = "internal server Erroror"
	ErrorDataLoss ResponseErrorType = "data loss"
	// 501 Not Implemented
	ErrorUnimplemented ResponseErrorType = "unimplemented"
	// 503 Service Unavailable
	ErrorUnavailable ResponseErrorType = "unavailable"
	// 504 Gateway Timeout
	ErrorDeadlineExceeded ResponseErrorType = "deadline exceeded"
)
