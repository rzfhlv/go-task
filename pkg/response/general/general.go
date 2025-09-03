package general

import (
	"math"

	"github.com/rzfhlv/go-task/pkg/param"
)

type Response struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
	Meta    any     `json:"meta,omitempty"`
	Data    any     `json:"data,omitempty"`
	Error   any     `json:"error,omitempty"`
}

type Meta struct {
	Limit     int   `json:"limit"`
	Page      int   `json:"page"`
	PerPage   int   `json:"perPage"`
	PageCount int   `json:"pageCount"`
	Total     int64 `json:"total"`
}

func BuildMeta(param param.Param, data int) Meta {
	pageCount := 0
	if param.Limit > 0 {
		pageCount = int(math.Ceil(float64(param.Total) / float64(param.Limit)))
	}
	return Meta{
		Limit:     param.Limit,
		Page:      param.Page,
		PerPage:   data,
		PageCount: pageCount,
		Total:     param.Total,
	}
}

func Set(success bool, msg *string, meta, data, err any) Response {
	return Response{
		Success: success,
		Message: msg,
		Meta:    meta,
		Data:    data,
		Error:   err,
	}
}
