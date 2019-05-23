package api

type ApiError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func Error(status int, message string) error {
	return ApiError{
		Status:  status,
		Message: message,
	}
}

func (this ApiError) Error() string {
	return "API Error: " + this.Message
}
