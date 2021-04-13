package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var version string = "snapshot"
var DEFAULT_TIMEOUT time.Duration = time.Duration(30) * time.Second
var DATE_FORMAT_STRING string = "2006-01-02"

type date time.Time

func newDate(val time.Time, p *time.Time) *date {
	*p = val
	return (*date)(p)
}

func (d *date) Set(s string) error {
	v, err := time.Parse(DATE_FORMAT_STRING, s)
	*d = date(v)
	return err
}

func (d *date) Type() string { return "date" }

func (d *date) String() string {
    if (*time.Time) (d).IsZero() {
        return ""
    } else {
        return (*time.Time)(d).Format(DATE_FORMAT_STRING)
    }
}

func main() {
	app := &cobra.Command{
		Use:   "ttrack [command]",
		Short: "A time tracking program.",
	}

	beginDate := newDate(time.Time{}, &time.Time{})
	endDate := newDate(time.Time{}, &time.Time{})

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
	var timeout time.Duration
	recCmd := &cobra.Command{
		Use:   "rec",
		Short: "record current time",
		Run: func(cmd *cobra.Command, args []string) {
    			fmt.Println("Recording")
    			fmt.Println(timeout.String())
		},
		Args: cobra.ExactArgs(1),
	}

	recCmd.Flags().DurationVarP(&timeout, "timeout", "t", DEFAULT_TIMEOUT, "timeout when to end recording")
	app.AddCommand(recCmd)

	//------ MV COMMAND
	mvCmd := &cobra.Command{
		Use:   "mv [flags] <srcGroup>... <dstGroup>",
		Short: "move/merge/rename groups",
		Run: func(cmd *cobra.Command, args []string) {
    			fmt.Println("moving")
		},
		Args: cobra.MinimumNArgs(2),
	}
	app.AddCommand(mvCmd)

	//------ DEL COMMAND
	delCmd := &cobra.Command{
		Use:   "del [flags] <group>",
		Short: "delete a group or date range in a group",
		Run: func(cmd *cobra.Command, args []string) {
    			fmt.Println("bd: " + beginDate.String())
    			fmt.Println("ed: " + endDate.String())
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
			fmt.Println("listing groups all right")
		},
		Args: cobra.ExactArgs(0),
	}
	app.AddCommand(groupsCmd)

	//------ LIST COMMAND
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "list dates with associated duration",
		Run: func(cmd *cobra.Command, args []string) {
    			fmt.Println("bd: " + beginDate.String())
    			fmt.Println("ed: " + endDate.String())
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
    			fmt.Println("bd: " + beginDate.String())
    			fmt.Println("ed: " + endDate.String())
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
			fmt.Println("set set set.")
		},
		Args: cobra.ExactArgs(3),
	}
	app.AddCommand(setCmd)

	app.Execute()
}
