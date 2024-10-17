package database

import (
	"fmt"
	"os"

	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/nktknshn/gomusiclibrary/lib/database/migration"
	"github.com/spf13/cobra"
)

var cmdMigrate = cobra.Command{
	Use:  "migrate <migrations-folder>",
	Run:  handleMigrate,
	Args: cobra.ExactArgs(1),
}

func handleMigrate(cmd *cobra.Command, args []string) {

	dbfile := cli.GetDatabaseFileMust()

	fmt.Println("Migrating", dbfile)

	if err := migration.MigrateFile(cmd.Context(), dbfile, os.DirFS(args[0])); err != nil {
		panic(err)
	}

}
