package logger_test

import (
	"testing"

	"github.com/rzfhlv/go-task/pkg/logger"
)

func TestLoggerSetDefault(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{
			name: "success logger level debug", level: "DEBUG",
		},
		{
			name: "success logger level error", level: "ERROR",
		},
		{
			name: "success logger level info", level: "INFO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetDefault(tt.level)
		})
	}
}
