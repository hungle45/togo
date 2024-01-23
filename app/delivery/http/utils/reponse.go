package http_utils

import (
	"net/http"
	"togo/domain"
)

const (
	ReponseStatusSuccess  = "SUCCESSFUL"
	ResponseStatusFail    = "FAILED"
	ResponseStatusProcess = "PROCESSING"
	ResponseStatusPending = "PENDING"
)

func ResponseWithData(status string, data map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"status": status,
		"data":   data,
	}
}

func ResponseWithMessage(status string, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func GetStatusCode(rerr domain.ResponseError) int {
	if rerr == nil {
		return http.StatusOK
	}

	// 400 Bad Request
	if rerr.ErrorType() == domain.ErrorInvalidArgument ||
		rerr.ErrorType() == domain.ErrorFaildPrecondition ||
		rerr.ErrorType() == domain.ErrorOutOfRange {
		return http.StatusBadRequest
	}

	// 401 Unauthorized
	if rerr.ErrorType() == domain.ErrorUnauthenticated {
		return http.StatusUnauthorized
	}

	// 403 Forbidden
	if rerr.ErrorType() == domain.ErrorPermissionDenied {
		return http.StatusForbidden
	}

	// 404 Not Found
	if rerr.ErrorType() == domain.ErrorNotFound {
		return http.StatusNotFound
	}

	// 409 Conflict
	if rerr.ErrorType() == domain.ErrorAlreadyExists ||
		rerr.ErrorType() == domain.ErrorAborted {
		return http.StatusConflict
	}

	// 429 Too Many Request
	if rerr.ErrorType() == domain.ErrorResourceExhausted {
		return http.StatusTooManyRequests
	}

	// 499 Client Closed Request
	if rerr.ErrorType() == domain.ErrorCancelled {
		return 499
	}

	// 500 Internal Server Error
	if rerr.ErrorType() == domain.ErrorUnknown ||
		rerr.ErrorType() == domain.ErrorInternal ||
		rerr.ErrorType() == domain.ErrorDataLoss {
		return http.StatusInternalServerError
	}

	// 501 Not Implemented
	if rerr.ErrorType() == domain.ErrorUnimplemented {
		return http.StatusNotImplemented
	}

	// 503 Service Unavailable
	if rerr.ErrorType() == domain.ErrorUnavailable {
		return http.StatusServiceUnavailable
	}

	// 504 Gateway Timeout
	if rerr.ErrorType() == domain.ErrorDeadlineExceeded {
		return http.StatusGatewayTimeout
	}

	// Default to Internal Server Error
	return http.StatusInternalServerError
}
