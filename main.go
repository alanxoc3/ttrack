package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "snapshot"

func main() {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ttrack " + version)
		},
	}

	app := &cobra.Command{
		Use:   "ttrack [subcommand] [flags] args",
		Short: "A time tracking program.",
	}
	app.AddCommand(versionCmd)
	app.Execute()
}
