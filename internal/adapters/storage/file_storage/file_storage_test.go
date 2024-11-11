package filestorage_test

import (
	"context"
	"errors"
	"os"
	"testing"

	fs "github.com/MomsEngineer/urlshortener/internal/adapters/storage/file_storage"
	"github.com/MomsEngineer/urlshortener/internal/entities/link"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {
	path := "test.json"
	store, err := fs.NewFileStorage(path)
	require.NoError(t, err)

	defer func() {
		if store != nil {
			store.Close()
		}
		os.Remove(path)
	}()

	l, err := link.NewLink("example", "https://example.com")
	require.NoError(t, err)

	err = store.SaveLink(context.TODO(), l)
	require.NoError(t, err)

	l, err = link.NewLink("test", "https://test.com")
	require.NoError(t, err)

	err = store.SaveLink(context.TODO(), l)
	require.NoError(t, err)

	tests := []struct {
		name  string
		short string
		want  string
		err   error
	}{
		{
			name:  "Get an existing first link",
			short: "example",
			want:  "https://example.com",
			err:   nil,
		},
		{
			name:  "Get a non-existing link",
			short: "abc",
			want:  "",
			err:   errors.New("not found"),
		},
		{
			name:  "Get an existing second link",
			short: "test",
			want:  "https://test.com",
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := link.NewLink(tt.short, "")
			require.NoError(t, err)

			err = store.GetLink(context.TODO(), l)

			require.Equal(t, tt.err, err, "expected err does not match")
			assert.Equal(t, tt.want, l.OriginalURL, "expected link does not match")
		})
	}
}
