package scan

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/nktknshn/gomusiclibrary/lib/database"
	"github.com/nktknshn/gomusiclibrary/lib/database/models"
	"github.com/spf13/cobra"
)

func init() {
	cli.Cmd.AddCommand(&cmdScan)
}

var cmdScan = cobra.Command{
	Use:  "scan <folder>",
	Args: cobra.ExactArgs(1),
	Run:  scan,
}

func scan(cmd *cobra.Command, args []string) {

	db, err := database.NewFromFile(cli.GetDatabaseFileMust())

	if err != nil {
		panic(err)
	}

	folder := args[0]

	var d fs.FS = os.DirFS(folder)
	var files []models.File

	err = fs.WalkDir(d, ".", func(p string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()

		if err != nil {
			return err
		}

		f := models.File{
			Path:  path.Join(folder, p),
			Size:  info.Size(),
			Mtime: info.ModTime(),
		}

		files = append(files, f)

		return err
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("Insert %d files.\n", len(files))

	err = db.FilesInsert(files)

	if err != nil {
		panic(err)
	}

}
