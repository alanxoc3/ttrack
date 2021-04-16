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

	//------ MV COMMAND
	mvCmd := &cobra.Command{
		Use:   "mv <srcGroup> <dstGroup>",
		Short: "move/merge/rename groups",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("moving")
		},
		Args: cobra.ExactArgs(2),
	}
	app.AddCommand(mvCmd)

	//------ DEL COMMAND
	delCmd := &cobra.Command{
		Use:   "del [flags] <group>",
		Short: "delete a group or date range in a group",
		Run: func(cmd *cobra.Command, args []string) {
    			delFunc(args[0], time.Time(*beginDate), time.Time(*endDate))
		},
		Args: cobra.ExactArgs(1),
	}

	delCmd.Flags().VarPF(beginDate, "begin-date", "b", "begin date for deletion range")
	delCmd.Flags().VarPF(endDate, "end-date", "e", "end date for deletion range")
	app.AddCommand(delCmd)

	//------ GROUPS COMMAND
	groupsCmd := &cobra.Command{
		Use:   "groups",
		Short: "show available groups",
		Run: func(cmd *cobra.Command, args []string) {
    			groupsFunc()
		},
		Args: cobra.ExactArgs(0),
	}
	app.AddCommand(groupsCmd)

	//------ LIST COMMAND
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list dates with associated duration",
		Run: func(cmd *cobra.Command, args []string) {
    			listFunc(args[0])
		},
		Args: cobra.ExactArgs(1),
	}

	listCmd.Flags().VarPF(beginDate, "begin-date", "b", "only list dates after this")
	listCmd.Flags().VarPF(endDate, "end-date", "e", "only list dates before this")
	app.AddCommand(listCmd)

	//------ AGG COMMAND
	aggCmd := &cobra.Command{
		Use:   "agg",
		Short: "aggregate dates for range into single duration",
		Run: func(cmd *cobra.Command, args []string) {
    			aggFunc(args[0])
		},
		Args: cobra.ExactArgs(1),
	}

	aggCmd.Flags().VarPF(beginDate, "begin-date", "b", "only aggregate dates after this")
	aggCmd.Flags().VarPF(endDate, "end-date", "e", "only aggregate dates before this")
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
