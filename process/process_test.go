package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		ProcessName  string
		ExpectErr    error
		ExpectResult bool
	}{
		{
			ProcessName:  "go",
			ExpectResult: true,
			ExpectErr:    nil,
		},

		{
			ProcessName:  "maybe-not-exist",
			ExpectResult: false,
			ExpectErr:    nil,
		},
	}

	manager := NewProcessManager()

	for _, test := range tests {
		isExec, err := manager.IsRunning(test.ProcessName)

		assert.Equal(t, test.ExpectErr, err)
		assert.Equal(t, test.ExpectResult, isExec)
	}
}
