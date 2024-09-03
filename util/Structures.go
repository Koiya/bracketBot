package util

type ErrorWrapper struct {
	Errors []ErrorDetail `json:"errors"`
}

type ErrorDetail struct {
	Detail string `json:"detail"`
}

type ErrorWrapperArray struct {
	Errors []ErrorDetailArray `json:"errors"`
}

type ErrorDetailArray struct {
	Detail []string `json:"detail"`
}
