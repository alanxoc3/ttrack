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
type execFunc func(*cmds.State)string

func parseRecArgs(args []string, s *cmds.State) {
	dur, durerr := time.ParseDuration(args[1])
	if durerr != nil {
		panic(durerr)
	}
	s.Groups = []types.Group{types.CreateGroupFromString(args[0])}
	s.Duration = types.CreateSecondsFromDuration(dur)
}

func parseGroup(args []string, s *cmds.State) {
	s.Groups = []types.Group{types.CreateGroupFromString(args[0])}
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
            if !s.Cached && !s.Stored {
                s.Cached = true
                s.Stored = true
            }

            pf(args, s)
            output := ef(s)
            if output != "" {
                fmt.Print(output)
            }
        },
	}
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
	lsCmd      := createCmd(&s, parseGroups,  cmds.LsFunc)
	recCmd     := createCmd(&s, parseRecArgs, cmds.RecFunc)
	resetCmd   := createCmd(&s, parseGroup,   cmds.ResetFunc)
	setCmd     := createCmd(&s, parseSetArgs, cmds.SetFunc)
	subCmd     := createCmd(&s, parseSetArgs, cmds.SubFunc)
	tidyCmd    := createCmd(&s, parseNothing, cmds.TidyFunc)
	versionCmd := createCmd(&s, parseNothing, func(s *cmds.State)string { return "ttrack " + version })

	aggCmd.Flags().BoolVarP(&s.Cached,    "cached",     "c", false, "only aggregate for cached data")
	aggCmd.Flags().BoolVarP(&s.Daily,     "daily",      "d", false, "aggregate per day instead of all together")
	aggCmd.Flags().BoolVarP(&s.Stored,    "stored",     "s", false, "only aggregate for stored data")
	aggCmd.Flags().VarPF   (&s.BeginDate, "begin-date", "b",        "only aggregate dates after or equal to this")
	aggCmd.Flags().VarPF   (&s.EndDate,   "end-date",   "e",        "only aggregate dates before or equal to this")

	lsCmd .Flags().BoolVarP(&s.Cached,    "cached",     "c", false, "only list cached groups")
	lsCmd .Flags().BoolVarP(&s.Recursive, "recursive",  "r", false, "list subgroups recursively")
	lsCmd .Flags().BoolVarP(&s.Stored,    "stored",     "s", false, "only list stored groups")

	setCmdMeta(app, addCmd    , 3, true , "add <group> <date> <duration>", "adds the duration for a group's date")
	setCmdMeta(app, aggCmd    , 0, false, "agg [<group>]..."             , "aggregate durations for date range")
	setCmdMeta(app, lsCmd     , 0, false, "ls [<group>]..."              , "list groups")
	setCmdMeta(app, recCmd    , 2, true , "rec <group> <duration>"       , "record current time")
	setCmdMeta(app, resetCmd  , 1, true,  "reset <group>"                , "resets a group's recording")
	setCmdMeta(app, setCmd    , 3, true , "set <group> <date> <duration>", "sets the duration for a group's date")
	setCmdMeta(app, subCmd    , 3, true , "sub <group> <date> <duration>", "subtracts the duration for a group's date")
	setCmdMeta(app, tidyCmd   , 0, true , "tidy"                         , "clean up all the files ttrack uses")
	setCmdMeta(app, versionCmd, 0, true , "version"                      , "print build information")

	app.Execute()
}
