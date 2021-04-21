package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var version string = "snapshot"

func main() {
	app := &cobra.Command{
		Use:   "ttrack [command]",
		Short: "A time tracking program.",
	}

	beginDate := &date{}
	endDate := &date{}

	//------ VERSION COMMAND
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ttrack " + version)
		},
	}
	app.AddCommand(versionCmd)

	//------ REC COMMAND
	recCmd := &cobra.Command{
		Use:   "rec <group> <timeout>",
		Short: "record current time",
		Run: func(cmd *cobra.Command, args []string) {
    			dur, durerr := time.ParseDuration(args[1])
    			if durerr != nil { panic(durerr) }

			recFunc(args[0], (uint32)(dur.Milliseconds()/1000))
		},
		Args: cobra.ExactArgs(2),
	}

	app.AddCommand(recCmd)

	//------ CP COMMAND
	cpCmd := &cobra.Command{
		Use:   "cp <srcGroup> <dstGroup>",
		Short: "copy/merge groups",
		Run: func(cmd *cobra.Command, args []string) {
    			cpFunc(args[0], args[1], beginDate.String(), endDate.String())
		},
		Args: cobra.ExactArgs(2),
	}
	cpCmd.Flags().VarPF(beginDate, "begin-date", "b", "only copy dates after or equal to this")
	cpCmd.Flags().VarPF(endDate, "end-date", "e", "only copy dates before or equal to this")
	app.AddCommand(cpCmd)

	//------ DEL COMMAND
	delCmd := &cobra.Command{
		Use:   "del <group>",
		Short: "delete a group",
		Run: func(cmd *cobra.Command, args []string) {
    			delFunc(args[0])
		},
		Args: cobra.ExactArgs(1),
	}

	app.AddCommand(delCmd)

	//------ LIST COMMAND
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list groups",
		Run: func(cmd *cobra.Command, args []string) {
    			listFunc()
		},
		Args: cobra.ExactArgs(0),
	}
	app.AddCommand(listCmd)

	//------ VIEW COMMAND
	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "view dates with associated duration",
		Run: func(cmd *cobra.Command, args []string) {
    			viewFunc(args[0], beginDate.String(), endDate.String())
		},
		Args: cobra.ExactArgs(1),
	}

	viewCmd.Flags().VarPF(beginDate, "begin-date", "b", "only view dates after or equal to this")
	viewCmd.Flags().VarPF(endDate, "end-date", "e", "only view dates before or equal to this")
	app.AddCommand(viewCmd)

	//------ AGG COMMAND
	aggCmd := &cobra.Command{
		Use:   "agg",
		Short: "aggregate dates for range into single duration",
		Run: func(cmd *cobra.Command, args []string) {
    			aggFunc(args[0], beginDate.String(), endDate.String())
		},
		Args: cobra.ExactArgs(1),
	}

	aggCmd.Flags().VarPF(beginDate, "begin-date", "b", "only aggregate dates after or equal to this")
	aggCmd.Flags().VarPF(endDate, "end-date", "e", "only aggregate dates before or equal to this")
	app.AddCommand(aggCmd)

	//------ SET COMMAND
	setCmd := &cobra.Command{
		Use:   "set <group> <date> <duration>",
		Short: "sets the duration for a group's date",
		Run: func(cmd *cobra.Command, args []string) {
    			g := args[0]
    			ts, tserr := dateStrToTimestamp(args[1])
    			if tserr != nil { panic(tserr) }

    			dur, durerr := time.ParseDuration(args[2])
    			if durerr != nil { panic(durerr) }

    			setFunc(g, ts, (uint32)(dur.Milliseconds()/1000))
		},
		Args: cobra.ExactArgs(3),
	}
	app.AddCommand(setCmd)

	app.Execute()
}
