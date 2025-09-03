package errs

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewErrs(code int, msg string) *HttpError {
	return &HttpError{
		StatusCode: code,
		Message:    msg,
	}
}
