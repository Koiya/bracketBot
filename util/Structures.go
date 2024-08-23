package util

type ErrorWrapper struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Detail string `json:"detail"`
}
