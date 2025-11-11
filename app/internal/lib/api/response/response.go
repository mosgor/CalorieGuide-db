package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "ok"
	StatusError = "error"
)

func Ok() Response {
	return Response{Status: StatusOK}
}

func Error(err string) Response {
	return Response{Status: StatusError, Error: err}
}
