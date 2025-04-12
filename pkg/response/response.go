package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type ProResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details `json:"details"`
}

type Details struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ProError(code int, msg string, details Details) ProResponse {
	return ProResponse{
		Code:    code,
		Message: msg,
		Details: details,
	}
}
