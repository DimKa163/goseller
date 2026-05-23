package domain

import "testing"

func TestCategoryID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      "123e4567-e89b-12d3-a456-426614174000",
			wantErr: false,
		},
		{
			name:    "invalid id",
			id:      "invalid-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCategoryID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCategoryID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
