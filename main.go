package main

import (
	"fmt"
	"time"
	"os"
	"os/user"

	"github.com/alanxoc3/ttrack/internal/types"
	"github.com/alanxoc3/ttrack/internal/cmds"
	"github.com/spf13/cobra"
)

var version string = "snapshot"

type parseFunc func([]string, *cmds.State)
type cobraFunc func(*cobra.Command, []string)
type execFunc func(*cmds.State)

func parseRecArgs(args []string, s *cmds.State) {
	dur, durerr := time.ParseDuration(args[1])
	if durerr != nil {
		panic(durerr)
	}
	s.Groups = []types.Group{types.CreateGroupFromString(args[0])}
	s.Duration = types.CreateSecondsFromDuration(dur)
}

func parseGroups(args []string, s *cmds.State) {
    groups := []types.Group{}
    for _, v := range args {
        groups = append(groups, types.CreateGroupFromString(v))
    }
    s.Groups = groups
}

func parseSetArgs(args []string, s *cmds.State) {
	ts, tserr := types.CreateDateFromString(args[1])
	if tserr != nil {
		panic(tserr)
	}

	dur, durerr := time.ParseDuration(args[2])
	if durerr != nil {
		panic(durerr)
	}

    s.Duration = types.CreateSecondsFromDuration(dur)
    s.Groups = []types.Group{types.CreateGroupFromString(args[0])}
    s.Date = *ts
}

func parseNothing(args []string, s *cmds.State) {
    return
}

func createCmd(s *cmds.State, pf parseFunc, ef execFunc) *cobra.Command {
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

func getEnvDir(appVar, xdgVar, homeFallback string) string {
	if val, present := os.LookupEnv(appVar); present {
		return val
	} else if val, present := os.LookupEnv(xdgVar); present {
		return val
	} else if usr, err := user.Current(); err == nil {
		return usr.HomeDir + homeFallback
	} else {
		return ""
	}
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
    s := cmds.State{}
    s.Now = time.Now()
    s.DataDir  = getEnvDir("TTRACK_DATA_DIR" , "XDG_DATA_HOME" , "/.local/share/ttrack")
    s.CacheDir = getEnvDir("TTRACK_CACHE_DIR", "XDG_CACHE_HOME", "/.cache/ttrack")

	app := &cobra.Command{
        Use: "ttrack [command]",
    	Short: "A time tracking program.",
	}

	addCmd     := createCmd(&s, parseSetArgs, cmds.AddFunc)
	aggCmd     := createCmd(&s, parseGroups,  cmds.AggFunc)
	cpCmd      := createCmd(&s, parseGroups,  func(s *cmds.State) { fmt.Println("in implementation") })
	lsCmd    := createCmd(&s,   parseGroups, cmds.LsFunc)
	mvCmd      := createCmd(&s, parseGroups,  func(s *cmds.State) { fmt.Println("in implementation") })
	recCmd     := createCmd(&s, parseRecArgs, cmds.RecFunc)
	rmCmd      := createCmd(&s, parseGroups,  func(s *cmds.State) { fmt.Println("in implementation") })
	setCmd     := createCmd(&s, parseSetArgs, cmds.SetFunc)
	subCmd     := createCmd(&s, parseSetArgs, cmds.SubFunc)
	tidyCmd    := createCmd(&s, parseNothing, func(s *cmds.State) { fmt.Println("tidying") })
	versionCmd := createCmd(&s, parseNothing, func(s *cmds.State) { fmt.Println("ttrack " + version) })

	aggCmd.Flags().BoolVarP(&s.Daily,     "daily",      "d", false, "aggregate per day instead of all together")
	aggCmd.Flags().BoolVarP(&s.Recursive, "recursive",  "r", false, "aggregate includes all sub groups recursively")
	aggCmd.Flags().VarPF   (&s.BeginDate, "begin-date", "b",        "only aggregate dates after or equal to this")
	aggCmd.Flags().VarPF   (&s.EndDate,   "end-date",   "e",        "only aggregate dates before or equal to this")
	cpCmd .Flags().VarPF   (&s.BeginDate, "begin-date", "b",        "only copy dates after or equal to this")
	cpCmd .Flags().VarPF   (&s.EndDate,   "end-date",   "e",        "only copy dates before or equal to this")
	lsCmd .Flags().BoolVarP(&s.Recursive, "recursive",  "r", false, "list subgroups recursively")
	mvCmd .Flags().VarPF   (&s.BeginDate, "begin-date", "b",        "only move/merge dates after or equal to this")
	mvCmd .Flags().VarPF   (&s.EndDate,   "end-date",   "e",        "only move/merge dates before or equal to this")
	rmCmd .Flags().BoolVarP(&s.Recursive, "recursive",  "r", false, "delete records in subgroups too")
	rmCmd .Flags().VarPF   (&s.BeginDate, "begin-date", "b",        "only delete records after or equal to this")
	rmCmd .Flags().VarPF   (&s.EndDate,   "end-date",   "e",        "only delete records before or equal to this")

	setCmdMeta(app, addCmd    , 3, true , "add <group> <date> <duration>", "adds the duration for a group's date")
	setCmdMeta(app, aggCmd    , 1, false, "agg <group>..."               , "aggregate dates for range into single duration")
	setCmdMeta(app, cpCmd     , 2, false, "cp <src-group>... <dst-group>", "copy/merge groups")
	setCmdMeta(app, lsCmd     , 0, false, "ls [<group>]..."              , "list groups")
	setCmdMeta(app, mvCmd     , 2, false, "mv <src-group>... <dst-group>", "rename/merge groups")
	setCmdMeta(app, recCmd    , 2, true , "rec <group> <duration>"       , "record current time")
	setCmdMeta(app, rmCmd     , 1, false, "rm <group>..."                , "remove a group")
	setCmdMeta(app, setCmd    , 3, true , "set <group> <date> <duration>", "sets the duration for a group's date")
	setCmdMeta(app, subCmd    , 3, true , "sub <group> <date> <duration>", "subtracts the duration for a group's date")
	setCmdMeta(app, tidyCmd   , 0, true , "tidy"                         , "clean up all the files ttrack uses")
	setCmdMeta(app, versionCmd, 0, true , "version"                      , "print build information")

	app.Execute()
}
