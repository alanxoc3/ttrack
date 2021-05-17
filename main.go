package main

import (
	"fmt"
	"time"

	"github.com/alanxoc3/ttrack/internal/date"
	"github.com/alanxoc3/ttrack/internal/seconds"
	"github.com/spf13/cobra"
)

var version string = "snapshot"

type state struct {
	beginDate dateArg
	endDate dateArg
	recursive bool
	daily bool
	quote bool
	groups []string
	date date.Date
	duration seconds.Seconds
}

type parseFunc func([]string, *state)
type cobraFunc func(*cobra.Command, []string)
type execFunc func(*state)

func parseRecArgs(args []string, s *state) {
	dur, durerr := time.ParseDuration(args[1])
	if durerr != nil {
		panic(durerr)
	}
	s.groups = []string{args[0]}
	s.duration = seconds.CreateFromDuration(dur)
}

func parseGroups(args []string, s *state) {
    s.groups = args
}

func parseSetArgs(args []string, s *state) {
	ts, tserr := date.CreateFromString(args[1])
	if tserr != nil {
		panic(tserr)
	}

	dur, durerr := time.ParseDuration(args[2])
	if durerr != nil {
		panic(durerr)
	}

    s.duration = seconds.CreateFromDuration(dur)
    s.groups = []string{args[0]}
    s.date = *ts
}

func parseNothing(args []string, s *state) {
    return
}

func createCmd(s *state, pf parseFunc, ef execFunc) *cobra.Command {
	return &cobra.Command{
        Run: func(cmd *cobra.Command, args []string) {
            pf(args, s)
            ef(s)
        },
	}
}

func setStateMaxArg(argCount int, pfargs []string) {
	fmt.Println("ttrack " + version)
}

func setCmdMeta(app, cmd *cobra.Command, argCount int, exact bool, use, short string) {
    if exact {
        cmd.Args = cobra.ExactArgs(argCount)
    } else {
        cmd.Args = cobra.MinimumNArgs(argCount)
    }

    cmd.Use = use
    cmd.Short = short
	app.AddCommand(cmd)
}

func main() {
    s := state{}
	app := &cobra.Command{
        Use: "ttrack [command]",
    	Short: "A time tracking program.",
	}

	addCmd     := createCmd(&s, parseSetArgs, func(s *state) { setFunc(s.groups[0], s.date, s.duration) })
	aggCmd     := createCmd(&s, parseGroups,  func(s *state) { aggFunc(s.groups[0], s.beginDate.String(), s.endDate.String()) })
	cpCmd      := createCmd(&s, parseGroups,  func(s *state) { cpFunc(s.groups[0], s.groups[1], s.beginDate.ToDate(), s.endDate.ToDate()) })
	listCmd    := createCmd(&s, parseNothing, func(s *state) { listFunc() })
	mvCmd      := createCmd(&s, parseGroups,  func(s *state) { cpFunc(s.groups[0], s.groups[1], s.beginDate.ToDate(), s.endDate.ToDate()) })
	recCmd     := createCmd(&s, parseRecArgs, func(s *state) { recFunc(s.groups[0], s.duration) })
	rmCmd      := createCmd(&s, parseGroups,  func(s *state) { delFunc(s.groups[0]) })
	setCmd     := createCmd(&s, parseSetArgs, func(s *state) { setFunc(s.groups[0], s.date, s.duration) })
	subCmd     := createCmd(&s, parseSetArgs, func(s *state) { setFunc(s.groups[0], s.date, s.duration) })
	tidyCmd    := createCmd(&s, parseNothing, func(s *state) { fmt.Println("tidying") })
	versionCmd := createCmd(&s, parseNothing, func(s *state) { fmt.Println("ttrack " + version) })

	aggCmd .Flags().BoolVarP(&s.daily,     "daily",      "d", false, "aggregate per day instead of all together")
	aggCmd .Flags().BoolVarP(&s.recursive, "recursive",  "r", false, "aggregate includes all sub groups recursively")
	aggCmd .Flags().VarPF   (&s.beginDate, "begin-date", "b",        "only aggregate dates after or equal to this")
	aggCmd .Flags().VarPF   (&s.endDate,   "end-date",   "e",        "only aggregate dates before or equal to this")
	cpCmd  .Flags().VarPF   (&s.beginDate, "begin-date", "b",        "only copy dates after or equal to this")
	cpCmd  .Flags().VarPF   (&s.endDate,   "end-date",   "e",        "only copy dates before or equal to this")
	listCmd.Flags().BoolVarP(&s.quote,     "quote",      "q", false, "quote each group according to posix shell quoting rules")
	mvCmd  .Flags().VarPF   (&s.beginDate, "begin-date", "b",        "only move/merge dates after or equal to this")
	mvCmd  .Flags().VarPF   (&s.endDate,   "end-date",   "e",        "only move/merge dates before or equal to this")
	rmCmd  .Flags().BoolVarP(&s.recursive, "recursive",  "r", false, "delete records in subgroups too")
	rmCmd  .Flags().VarPF   (&s.beginDate, "begin-date", "b",        "only delete records after or equal to this")
	rmCmd  .Flags().VarPF   (&s.endDate,   "end-date",   "e",        "only delete records before or equal to this")

	setCmdMeta(app, addCmd    , 3, true , "add <group> <date> <duration>", "adds the duration for a group's date")
	setCmdMeta(app, aggCmd    , 1, true , "agg <group>..."               , "aggregate dates for range into single duration")
	setCmdMeta(app, cpCmd     , 2, false, "cp <src-group>... <dst-group>", "copy/merge groups")
	setCmdMeta(app, listCmd   , 0, true , "ls [<group>]..."              , "list groups")
	setCmdMeta(app, mvCmd     , 2, false, "mv <src-group>... <dst-group>", "rename/merge groups")
	setCmdMeta(app, recCmd    , 2, true , "rec <group> <duration>"       , "record current time")
	setCmdMeta(app, rmCmd     , 1, false, "rm <group>..."                , "remove a group")
	setCmdMeta(app, setCmd    , 3, true , "set <group> <date> <duration>", "sets the duration for a group's date")
	setCmdMeta(app, subCmd    , 3, true , "sub <group> <date> <duration>", "subtracts the duration for a group's date")
	setCmdMeta(app, tidyCmd   , 0, true , "tidy"                         , "clean up all the files ttrack uses")
	setCmdMeta(app, versionCmd, 0, true , "version"                      , "print build information")

	app.Execute()
}
