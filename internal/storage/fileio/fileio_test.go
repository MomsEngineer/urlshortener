package fileio_test

import (
	"os"
	"testing"

	"github.com/MomsEngineer/urlshortener/internal/storage/fileio"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileIO(t *testing.T) {
	fileName := "test.json"
	file, err := fileio.NewFileIO(fileName)
	require.NoError(t, err)

	defer func() {
		file.Close()
		os.Remove(fileName)
	}()

	short, origin := "example", "https://example.com"

	err = file.Write(short, origin)
	require.NoError(t, err)

	m, err := file.Read()
	require.NoError(t, err)

	require.Contains(t, m, short)
	assert.Equal(t, m[short], origin)
}
