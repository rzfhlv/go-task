package param_test

import (
	"testing"

	"github.com/rzfhlv/go-task/pkg/param"
	"github.com/stretchr/testify/assert"
)

func TestParam(t *testing.T) {
	param := param.Param{
		Page:   2,
		Limit:  10,
		Offset: 0,
		Total:  0,
	}

	nweOffset := param.CalculateOffset()
	assert.Equal(t, nweOffset, 10)
}
