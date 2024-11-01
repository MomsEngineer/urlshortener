package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDB(t *testing.T) {
	database := NewDB()
	require.NotNil(t, database)
}

func TestSavedLink(t *testing.T) {
	database := NewDB()
	require.NotNil(t, database)

	id := "abc123"
	link := "https://example.com"
	database.SaveLink(id, link)

	expectedLinks := map[string]string{
		id: link,
	}

	actualLinks := database.Links

	assert.Equal(t, expectedLinks, actualLinks,
		"The saved link does not match the expected link")
}

func TestGetLink(t *testing.T) {
	database := NewDB()
	require.NotNil(t, database)

	id := "abc123"
	link := "https://example.com"
	database.SaveLink(id, link)

	tests := []struct {
		name   string
		id     string
		want   string
		exists bool
	}{
		{
			name:   "Get an existing link",
			id:     "abc123",
			want:   "https://example.com",
			exists: true,
		},
		{
			name:   "Get a non-existing link",
			id:     "abc",
			want:   "",
			exists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, exists := database.GetLink(tt.id)

			require.Equal(t, tt.exists, exists, "expected existence status does not match")

			if tt.exists {
				assert.Equal(t, tt.want, link, "expected link does not match")
			} else {
				assert.Empty(t, link, "expected link to be empty")
			}
		})
	}
}
