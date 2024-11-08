package mapstorage_test

import (
	"testing"

	mapstorage "github.com/MomsEngineer/urlshortener/internal/storage/map_storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinksMap(t *testing.T) {
	lm := mapstorage.NewMapStorage()
	require.NotNil(t, lm)
}

func TestSavedLink(t *testing.T) {
	lm := mapstorage.NewMapStorage()
	require.NotNil(t, lm)

	id := "abc123"
	link := "https://example.com"
	lm.SaveLink(id, link)

	expectedLinks := map[string]string{
		id: link,
	}

	actualLinks := lm.Links

	assert.Equal(t, expectedLinks, actualLinks,
		"The saved link does not match the expected link")
}

func TestGetLink(t *testing.T) {
	lm := mapstorage.NewMapStorage()
	require.NotNil(t, lm)

	id := "abc123"
	link := "https://example.com"
	lm.SaveLink(id, link)

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
			link, exists, _ := lm.GetLink(tt.id)

			require.Equal(t, tt.exists, exists, "expected existence status does not match")

			if tt.exists {
				assert.Equal(t, tt.want, link, "expected link does not match")
			} else {
				assert.Empty(t, link, "expected link to be empty")
			}
		})
	}
}
