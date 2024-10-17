package database_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/nktknshn/gomusiclibrary/lib/database"
	"github.com/nktknshn/gomusiclibrary/lib/database/migration"
	"github.com/nktknshn/gomusiclibrary/lib/database/models"
	"github.com/nktknshn/gomusiclibrary/lib/tests"
	"github.com/stretchr/testify/require"
	"zombiezen.com/go/sqlite"
)

func conn() (*sqlite.Conn, error) {
	c, err := sqlite.OpenConn(":memory:")
	if err != nil {
		panic(err)
	}

	err = migration.MigrateConn(context.Background(), c, os.DirFS(tests.CurrentTestPath("migrations")))

	if err != nil {
		return nil, err
	}

	return c, err
}

var (
	files = models.FileSlice{
		{Path: "/test/1.mp3", Size: 667, Sha256Hash: "", Ctime: 0, Mtime: 0, CreatedAt: time.Now()},
		{Path: "/test/2.mp3", Size: 777, Sha256Hash: "", Ctime: 0, Mtime: 0, CreatedAt: time.Now()},
	}
)

func TestInsertAndList(t *testing.T) {

	c, err := conn()
	require.NoError(t, err)
	defer c.Close()

	d := database.New(c)
	err = d.FilesInsert(files)
	require.NoError(t, err)

	f, err := d.FilesList()
	require.NoError(t, err)

	require.Equal(t, []models.FileID{1, 2}, f.IDs())

}
