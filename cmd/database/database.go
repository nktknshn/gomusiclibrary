package database

import (
	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/spf13/cobra"
)

var cmd = cobra.Command{
	Use: "database",
}

func init() {
	cmd.AddCommand(&cmdMigrate)
	cmd.AddCommand(&cmdCreate)

	cli.Cmd.AddCommand(&cmd)
}
