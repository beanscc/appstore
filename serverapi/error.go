package serverapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ApiError https://developer.apple.com/documentation/appstoreserverapi/error_codes
type ApiError struct {
	Code    int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

func (e ApiError) Error() string {
	return fmt.Sprintf("errorCode:%d, errorMessage:%s", e.Code, e.Message)
}

func ParseApiError(data []byte) (*ApiError, error) {
	var e ApiError
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func ApiErrorFromError(err error) (*ApiError, bool) {
	if err == nil {
		return nil, true
	}

	var apiErr ApiError
	if errors.As(err, &apiErr) {
		return &apiErr, true
	}

	return nil, false
}

func handleApiErr(statusCode int, payload []byte) error {
	if statusCode != 200 {
		apiErr, err := ParseApiError(payload)
		if err != nil {
			return err
		}
		return apiErr
	}

	return nil
}
