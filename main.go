package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/alanxoc3/ttrack/internal/cmds"
	"github.com/alanxoc3/ttrack/internal/types"
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
	dur, durerr := time.ParseDuration(args[1])
	if durerr != nil {
		panic(durerr)
	}

    ts := types.CreateDateFromTime(time.Now())
    if len(args) == 3 {
        var tserr error
    	ts, tserr = types.CreateDateFromString(args[2])
    	if tserr != nil {
    		panic(tserr)
    	}
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
		return filepath.Join(usr.HomeDir, homeFallback)
	} else {
		return "."
	}
}

func setCmdMeta(app, cmd *cobra.Command, argMin, argMax int, use, short string) {
    if argMin == argMax {
        cmd.Args = cobra.ExactArgs(argMin)
    } else if argMax < 0 {
        cmd.Args = cobra.MinimumNArgs(argMin)
    } else {
        cmd.Args = cobra.RangeArgs(argMin, argMax)
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
	versionCmd := createCmd(&s, parseNothing, func(s *cmds.State)string { return "ttrack " + version + "\n" })

	aggCmd.Flags().BoolVarP(&s.Cached,    "cached",     "c", false, "only aggregate for cached data")
	aggCmd.Flags().BoolVarP(&s.Daily,     "daily",      "d", false, "aggregate per day instead of all together")
	aggCmd.Flags().BoolVarP(&s.Stored,    "stored",     "s", false, "only aggregate for stored data")
	aggCmd.Flags().VarPF   (&s.BeginDate, "begin-date", "b",        "only aggregate dates after or equal to this")
	aggCmd.Flags().VarPF   (&s.EndDate,   "end-date",   "e",        "only aggregate dates before or equal to this")

	lsCmd .Flags().BoolVarP(&s.Cached,    "cached",     "c", false, "only list cached groups")
	lsCmd .Flags().BoolVarP(&s.Recursive, "recursive",  "r", false, "list subgroups recursively")
	lsCmd .Flags().BoolVarP(&s.Stored,    "stored",     "s", false, "only list stored groups")

	setCmdMeta(app, addCmd    , 2, 3 , "add <group> <duration> [<date>]", "Adds the duration for a group's date (default is today)")
	setCmdMeta(app, aggCmd    , 0, -1, "agg [<group>]..."               , "Aggregate durations in groups or subgroups")
	setCmdMeta(app, lsCmd     , 0, -1, "ls [<group>]..."                , "List groups or subgroups")
	setCmdMeta(app, recCmd    , 2, 2 , "rec <group> <duration>"         , "Record duration for group to today")
	setCmdMeta(app, resetCmd  , 1, 1 , "reset <group>"                  , "Resets a group's recording")
	setCmdMeta(app, setCmd    , 2, 3 , "set <group> <duration> [<date>]", "Sets the duration for a group's date (default is today)")
	setCmdMeta(app, subCmd    , 2, 3 , "sub <group> <duration> [<date>]", "Subtracts the duration for a group's date (default is today)")
	setCmdMeta(app, tidyCmd   , 0, 0 , "tidy"                           , "Clean up all the files ttrack uses")
	setCmdMeta(app, versionCmd, 0, 0 , "version"                        , "Print build information")

	app.Execute()
}
