package main

import (
	"fmt"
	"time"

	"github.com/alanxoc3/ttrack/internal/date"
	"github.com/alanxoc3/ttrack/internal/seconds"
	"github.com/spf13/cobra"
)

var version string = "snapshot"

func main() {
	app := &cobra.Command{
		Use:   "ttrack [command]",
		Short: "A time tracking program.",
	}

	beginDate := dateArg{}
	endDate := dateArg{}
	recursive := false
	daily := false
	escape := false

	//------ VERSION COMMAND
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print build information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ttrack " + version)
		},
	}
	app.AddCommand(versionCmd)

	//------ TIDY COMMAND
	tidyCmd := &cobra.Command{
		Use:   "tidy",
		Short: "clean up all the files ttrack uses",
        Run: func(cmd *cobra.Command, args []string) {
    		// TODO: Implement me.
		},
		Args: cobra.ExactArgs(0),
	}

	app.AddCommand(tidyCmd)

	//------ LS COMMAND
	listCmd := &cobra.Command{
		Use:   "ls [group]...",
		Short: "list groups",
		Run: func(cmd *cobra.Command, args []string) {
			listFunc()
		},
		Args: cobra.ExactArgs(0),
	}
	listCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "list all subgroups too")
	listCmd.Flags().BoolVarP(&escape, "quote", "q", false, "quote each group according to posix shell quoting rules")
	app.AddCommand(listCmd)

	//------ REC COMMAND
	recCmd := &cobra.Command{
		Use:   "rec group... <duration>",
		Short: "record current time",
		Run: func(cmd *cobra.Command, args []string) {
			dur, durerr := time.ParseDuration(args[1])
			if durerr != nil {
				panic(durerr)
			}

			recFunc(args[0], seconds.CreateFromDuration(dur))
		},
		Args: cobra.ExactArgs(2),
	}

	app.AddCommand(recCmd)

	//------ SET COMMAND
	setCmd := &cobra.Command{
		Use:   "set group... <date>:[+|-]<duration>",
		Short: "sets the duration for a group's date",
		Run: func(cmd *cobra.Command, args []string) {
			g := args[0]
			ts, tserr := date.CreateFromString(args[1])
			if tserr != nil {
				panic(tserr)
			}

			dur, durerr := time.ParseDuration(args[2])
			if durerr != nil {
				panic(durerr)
			}

			setFunc(g, *ts, seconds.Seconds(dur.Milliseconds()/1000))
		},
		Args: cobra.MinimumNArgs(2),
	}
	app.AddCommand(setCmd)

	//------ MV COMMAND
	mvCmd := &cobra.Command{
		Use:   "mv <source-group>... <destination-group>",
		Short: "rename/merge groups",
		Run: func(cmd *cobra.Command, args []string) {
			cpFunc(args[0], args[1], beginDate.ToDate(), endDate.ToDate())
		},
		Args: cobra.MinimumNArgs(2),
	}
	mvCmd.Flags().VarPF(&beginDate, "begin-date", "b", "only move/merge dates after or equal to this")
	mvCmd.Flags().VarPF(&endDate,   "end-date",   "e", "only move/merge dates before or equal to this")
	app.AddCommand(mvCmd)


	//------ CP COMMAND
	cpCmd := &cobra.Command{
		Use:   "cp <src-group>... <dst-group>",
		Short: "copy/merge groups",
		Run: func(cmd *cobra.Command, args []string) {
			cpFunc(args[0], args[1], beginDate.ToDate(), endDate.ToDate())
		},
		Args: cobra.MinimumNArgs(2),
	}
	cpCmd.Flags().VarPF(&beginDate, "begin-date", "b", "only copy dates after or equal to this")
	cpCmd.Flags().VarPF(&endDate,   "end-date",   "e", "only copy dates before or equal to this")
	app.AddCommand(cpCmd)

	//------ DEL COMMAND
	delCmd := &cobra.Command{
		Use:   "del <group>...",
		Short: "delete a group",
		Run: func(cmd *cobra.Command, args []string) {
			delFunc(args[0])
		},
		Args: cobra.MinimumNArgs(1),
	}
	delCmd.Flags()   .VarPF(&beginDate, "begin-date", "b", "only delete records after or equal to this")
	delCmd.Flags()   .VarPF(&endDate,   "end-date",   "e", "only delete records before or equal to this")
	delCmd.Flags().BoolVarP(&recursive, "recursive",  "r", false, "delete records in subgroups too")

	app.AddCommand(delCmd)

	//------ AGG COMMAND
	aggCmd := &cobra.Command{
		Use:   "agg [<group>...]",
		Short: "aggregate dates for range into single duration",
		Run: func(cmd *cobra.Command, args []string) {
			aggFunc(args[0], beginDate.String(), endDate.String())
			// viewFunc(args[0], beginDate.ToDate(), endDate.ToDate())
		},
		Args: cobra.ExactArgs(1),
	}

	aggCmd.Flags().VarPF   (&beginDate, "begin-date", "b", "only aggregate dates after or equal to this")
	aggCmd.Flags().VarPF   (&endDate,   "end-date",   "e", "only aggregate dates before or equal to this")
	aggCmd.Flags().BoolVarP(&recursive, "recursive",  "r", false, "aggregate includes all sub groups recursively")
	aggCmd.Flags().BoolVarP(&daily,     "daily",      "d", false, "aggregate per day instead of all together")
	app.AddCommand(aggCmd)

	app.Execute()
}
