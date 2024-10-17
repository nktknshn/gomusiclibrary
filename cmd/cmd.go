package main

import (
	_ "github.com/nktknshn/gomusiclibrary/cmd/all"
	"github.com/nktknshn/gomusiclibrary/cmd/cli"
)

func main() {
	err := cli.Cmd.Execute()

	if err != nil {
		panic(err)
	}
}
