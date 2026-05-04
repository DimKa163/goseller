package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidEmail(t *testing.T) {
	tests := []struct {
		Name    string
		Email   Email
		IsValid bool
	}{
		{
			Name:    "Valid email",
			Email:   Email("test@example.com"),
			IsValid: true,
		},
		{
			Name:    "Invalid email",
			Email:   Email("invalid-email"),
			IsValid: false,
		},
		{
			Name:    "Invalid email with number forward",
			Email:   Email("1test@example.com"),
			IsValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := tt.Email.Validate()
			if tt.IsValid {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, ErrInvalidEmail, err)
			}
		})
	}
}
