package link

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantErr bool
	}{
		{
			name:    "Generate 8 character ID",
			length:  8,
			wantErr: false,
		},
		{
			name:    "Generate 0 character ID",
			length:  0,
			wantErr: true,
		},
		{
			name:    "Generate negative length ID",
			length:  -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GenerateID(tt.length)

			if tt.wantErr {
				require.Error(t, err, "expected an error but got none")
			} else {
				require.NoError(t, err, "unexpected error occurred")
				assert.Equal(t, tt.length, len(id), "ID length does not match")
			}
		})
	}
}
