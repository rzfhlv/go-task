package general

type Response struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
	Meta    any     `json:"meta,omitempty"`
	Result  any     `json:"result,omitempty"`
	Error   any     `json:"error,omitempty"`
}

func Set(success bool, msg *string, meta, result, err any) Response {
	return Response{
		Success: success,
		Message: msg,
		Meta:    meta,
		Result:  result,
		Error:   err,
	}
}
