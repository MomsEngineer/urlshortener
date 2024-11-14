package mapstorage_test

import (
	"context"
	"errors"
	"testing"

	ms "github.com/MomsEngineer/urlshortener/internal/adapters/storage/map_storage"
	"github.com/MomsEngineer/urlshortener/internal/entities/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinksMap(t *testing.T) {
	lm := ms.NewMapStorage()
	require.NotNil(t, lm)
}

func TestSavedLink(t *testing.T) {
	lm := ms.NewMapStorage()
	require.NotNil(t, lm)

	id := "abc123"
	original := "https://example.com"
	l, err := link.NewLink("userID", id, original)
	require.NoError(t, err)

	lm.SaveLink(context.TODO(), l)

	expectedLinks := map[string]string{
		id: original,
	}

	actualLinks := lm.Links

	assert.Equal(t, expectedLinks, actualLinks,
		"The saved link does not match the expected link")
}

func TestGetLink(t *testing.T) {
	lm := ms.NewMapStorage()
	require.NotNil(t, lm)

	id := "abc123"
	original := "https://example.com"

	l, err := link.NewLink("userID", id, original)
	require.NoError(t, err)

	lm.SaveLink(context.TODO(), l)

	tests := []struct {
		name string
		id   string
		want string
		err  error
	}{
		{
			name: "Get an existing link",
			id:   "abc123",
			want: "https://example.com",
			err:  nil,
		},
		{
			name: "Get a non-existing link",
			id:   "abc",
			want: "",
			err:  errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := link.NewLink("userID", tt.id, "")
			require.NoError(t, err)

			err = lm.GetLink(context.TODO(), l)

			require.Equal(t, tt.err, err, "expected err does not match")
			assert.Equal(t, tt.want, l.OriginalURL, "expected link does not match")
		})
	}
}
