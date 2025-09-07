package errs_test

import (
	"net/http"
	"testing"

	"github.com/rzfhlv/go-task/pkg/errs"
	"github.com/stretchr/testify/assert"
)

func TestErrs(t *testing.T) {
	httpErr := errs.HttpError{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
	}

	msg := httpErr.Error()
	assert.Equal(t, "internal server error", msg)
}

func TestErrsNewErrs(t *testing.T) {
	er := errs.NewErrs(http.StatusInternalServerError, "internal server error")

	assert.Equal(t, http.StatusInternalServerError, er.StatusCode)
	assert.Equal(t, "internal server error", er.Message)
}
