package main

import (
	"fmt"
	"time"

	"github.com/alanxoc3/ttrack/internal/date"
	"github.com/alanxoc3/ttrack/internal/seconds"
	"github.com/spf13/cobra"
)

var version string = "snapshot"

type params struct {
	beginDate dateArg
	endDate dateArg
	recursive bool
	daily bool
	quote bool
	groups []string
	date dateArg
	duration seconds.Seconds
}

func main() {
	app := &cobra.Command{
	}

    c := params{}

	versionCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ttrack " + version)
		},
	}

	tidyCmd := &cobra.Command{
		Args: cobra.ExactArgs(0),
        Run: func(cmd *cobra.Command, args []string) {
    		// TODO: Implement me.
		},
	}

	listCmd := &cobra.Command{
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			listFunc()
		},
	}

	recCmd := &cobra.Command{
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dur, durerr := time.ParseDuration(args[1])
			if durerr != nil {
				panic(durerr)
			}

			recFunc(args[0], seconds.CreateFromDuration(dur))
		},
	}

	setCmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
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
	}

	addCmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
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
	}

	subCmd := &cobra.Command{
		Args: cobra.ExactArgs(3),
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
	}

	mvCmd := &cobra.Command{
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cpFunc(args[0], args[1], c.beginDate.ToDate(), c.endDate.ToDate())
		},
	}

	cpCmd := &cobra.Command{
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cpFunc(args[0], args[1], c.beginDate.ToDate(), c.endDate.ToDate())
		},
	}

	rmCmd := &cobra.Command{
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			delFunc(args[0])
		},
	}

	aggCmd := &cobra.Command{
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			aggFunc(args[0], c.beginDate.String(), c.endDate.String())
			// viewFunc(args[0], beginDate.ToDate(), endDate.ToDate())
		},
	}

    // Set all the flags.
	aggCmd .Flags().BoolVarP(&c.daily,     "daily",      "d", false, "aggregate per day instead of all together")
	aggCmd .Flags().BoolVarP(&c.recursive, "recursive",  "r", false, "aggregate includes all sub groups recursively")
	aggCmd .Flags().VarPF   (&c.beginDate, "begin-date", "b",        "only aggregate dates after or equal to this")
	aggCmd .Flags().VarPF   (&c.endDate,   "end-date",   "e",        "only aggregate dates before or equal to this")
	cpCmd  .Flags().VarPF   (&c.beginDate, "begin-date", "b",        "only copy dates after or equal to this")
	cpCmd  .Flags().VarPF   (&c.endDate,   "end-date",   "e",        "only copy dates before or equal to this")
	listCmd.Flags().BoolVarP(&c.quote,     "quote",      "q", false, "quote each group according to posix shell quoting rules")
	mvCmd  .Flags().VarPF   (&c.beginDate, "begin-date", "b",        "only move/merge dates after or equal to this")
	mvCmd  .Flags().VarPF   (&c.endDate,   "end-date",   "e",        "only move/merge dates before or equal to this")
	rmCmd  .Flags().BoolVarP(&c.recursive, "recursive",  "r", false, "delete records in subgroups too")
	rmCmd  .Flags().VarPF   (&c.beginDate, "begin-date", "b",        "only delete records after or equal to this")
	rmCmd  .Flags().VarPF   (&c.endDate,   "end-date",   "e",        "only delete records before or equal to this")

	addCmd.Use     = "add <group> <date> <duration>"
	aggCmd.Use     = "agg <group>..."
	app.Use        = "ttrack [command]"
	cpCmd.Use      = "cp <source-group>... <destination-group>"
	listCmd.Use    = "ls [<group>]..."
	mvCmd.Use      = "mv <source-group>... <destination-group>"
	recCmd.Use     = "rec <group> <duration>"
	rmCmd.Use      = "rm <group>..."
	setCmd.Use     = "set <group> <date> <duration>"
	subCmd.Use     = "sub <group> <date> <duration>"
	tidyCmd.Use    = "tidy"
	versionCmd.Use = "version"

	addCmd.Short     = "adds the duration for a group's date"
	aggCmd.Short     = "aggregate dates for range into single duration"
	app.Short        = "A time tracking program."
	cpCmd.Short      = "copy/merge groups"
	listCmd.Short    = "list groups"
	mvCmd.Short      = "rename/merge groups"
	recCmd.Short     = "record current time"
	rmCmd.Short      = "remove a group"
	setCmd.Short     = "sets the duration for a group's date"
	subCmd.Short     = "subtracts the duration for a group's date"
	tidyCmd.Short    = "clean up all the files ttrack uses"
	versionCmd.Short = "print build information"

	app.AddCommand(addCmd)
	app.AddCommand(aggCmd)
	app.AddCommand(cpCmd)
	app.AddCommand(listCmd)
	app.AddCommand(mvCmd)
	app.AddCommand(recCmd)
	app.AddCommand(rmCmd)
	app.AddCommand(setCmd)
	app.AddCommand(subCmd)
	app.AddCommand(tidyCmd)
	app.AddCommand(versionCmd)

	app.Execute()
}
