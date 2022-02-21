package validator

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateInput(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedOutput error
	}{
		{
			name:           "Should return an error when file path is empty",
			input:          "",
			expectedOutput: fmt.Errorf("file path can't be empty"),
		},
		{
			name:           "Should return a nil error when file path is not empty",
			input:          "valid file path",
			expectedOutput: nil,
		},
	}

	userRetentionValidator := NewUserRetentionValidatorImpl()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			// Operation
			output := userRetentionValidator.ValidateInput(c.input)

			// Validation
			assert.EqualValues(t, c.expectedOutput, output)
		})
	}
}
