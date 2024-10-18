package library

import (
	"github.com/nktknshn/gomusiclibrary/cmd/cli"
	"github.com/spf13/cobra"
)

func init() {
	cli.Cmd.AddCommand(&cmdLibrary)
	cmdLibrary.AddCommand(&cmdLibraryList)
}

var cmdLibrary = cobra.Command{
	Use:  "library",
	Args: cobra.MinimumNArgs(1),
}
