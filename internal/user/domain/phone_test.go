package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhone(t *testing.T) {
	tests := []struct {
		Name    string
		Phone   Phone
		IsValid bool
	}{
		{
			Name:    "Valid phone",
			Phone:   Phone("1234567890"),
			IsValid: true,
		},
		{
			Name:    "Invalid phone",
			Phone:   Phone("+1234567890"),
			IsValid: false,
		},
		{
			Name:    "Invalid phone",
			Phone:   Phone("invalid-phone"),
			IsValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Phone.Validate()
			if tt.IsValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, ErrInvalidPhone, err)
			}
		})
	}
}
