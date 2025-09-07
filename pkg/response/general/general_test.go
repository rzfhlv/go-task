package general_test

import (
	"testing"

	"github.com/rzfhlv/go-task/pkg/param"
	"github.com/rzfhlv/go-task/pkg/response/general"
	"github.com/stretchr/testify/assert"
)

func TestResponseBuildMeta(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		expect := general.Meta{
			Limit:     10,
			Page:      2,
			PerPage:   10,
			PageCount: 2,
			Total:     20,
		}
		param := param.Param{
			Page:   2,
			Limit:  10,
			Offset: 10,
			Total:  20,
		}

		meta := general.BuildMeta(param, 10)

		assert.Equal(t, expect, meta)
	})
}

func TestResponseSet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		msg := "Success"
		expect := general.Response{
			Success: true,
			Message: &msg,
			Meta:    nil,
			Data:    nil,
			Error:   nil,
		}

		resp := general.Set(true, &msg, nil, nil, nil)

		assert.Equal(t, expect, resp)
	})
}
