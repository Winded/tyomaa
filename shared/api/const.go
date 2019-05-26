package api

// Method constants to avoid http package dependency
const (
	MethodGet    = "GET"
	MethodPost   = "POST"
	MethodDelete = "DELETE"
)

// Status constants to avoid http package dependency
const (
	StatusBadRequest = 400
	StatusNotFound   = 404
)
