package library

import (
	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/spf13/cobra"
)

var cmdLibraryList = cobra.Command{
	Use: "list",
	Run: list,
}

func list(cmd *cobra.Command, args []string) {
	db := cli.GetDatabaseMust()

	list, err := db.FilesList()

	if err != nil {
		panic(err)
	}

	for _, f := range list {
		println(f.Path)
	}
}
