package database

import (
	"github.com/spf13/cobra"
)

var cmdCreate = cobra.Command{
	Use: "create",
	Run: create,
}

func create(cmd *cobra.Command, args []string) {

}
