package filestorage_test

import (
	"context"
	"os"
	"testing"

	fs "github.com/MomsEngineer/urlshortener/internal/adapters/storage/file_storage"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {
	path := "test.json"
	fileStorage, err := fs.NewFileStorage(path)
	require.NoError(t, err)

	defer func() {
		if fileStorage != nil {
			fileStorage.Close()
		}
		os.Remove(path)
	}()

	short, origin := "example", "https://example.com"

	_, err = fileStorage.SaveLink(context.TODO(), short, origin)
	require.NoError(t, err)

	link, exists, err := fileStorage.GetLink(context.TODO(), short)
	require.NoError(t, err)

	require.True(t, exists, "expected exists to be true")
	require.Equal(t, link, origin)

	_, err = fileStorage.SaveLink(context.TODO(), "test", "https://test.com")
	require.NoError(t, err)

	link, exists, err = fileStorage.GetLink(context.TODO(), short)
	require.NoError(t, err)

	require.True(t, exists, "expected exists to be true")
	require.Equal(t, link, origin)
}
