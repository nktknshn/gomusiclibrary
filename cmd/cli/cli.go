package cli

import (
	"os"

	"github.com/nktknshn/gomusiclibrary/lib/database"
	"github.com/spf13/cobra"
)

const ENV_DATABASE_FILE_KEY = "MLDBFILE"

var (
	flagDatabaseFile = ""
)

var Cmd = cobra.Command{
	Use: "cli",
}

func GetDatabaseFileMust() string {
	if flagDatabaseFile != "" {
		return flagDatabaseFile
	}
	e := os.Getenv(ENV_DATABASE_FILE_KEY)
	if e == "" {
		panic("Database file is not set")
	}
	return e
}

func GetDatabaseMust() *database.Database {
	db, err := database.NewFromFile(GetDatabaseFileMust())
	if err != nil {
		panic(err)
	}
	return db
}

func init() {
	Cmd.PersistentFlags().StringVarP(&flagDatabaseFile, "database", "d", flagDatabaseFile, "sqlite3 database file")
}
