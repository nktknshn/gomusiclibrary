package scan

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/nktknshn/gomusiclibrary/lib/database"
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

	_, err := database.NewFromFile(cli.GetDatabaseFileMust())

	if err != nil {
		panic(err)
	}

	folder := args[0]

	var d fs.FS = os.DirFS(folder)

	err = fs.WalkDir(d, ".", func(path string, d fs.DirEntry, err error) error {
		fmt.Println(path)

		return err
	})

	if err != nil {
		panic(err)
	}

}
