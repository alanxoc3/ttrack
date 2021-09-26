# ttrack

[![Go Report Card](https://goreportcard.com/badge/github.com/alanxoc3/ttrack)](https://goreportcard.com/report/github.com/alanxoc3/ttrack)
[![Coverage Status](https://coveralls.io/repos/github/alanxoc3/ttrack/badge.svg?branch=main)](https://coveralls.io/github/alanxoc3/ttrack?branch=main)

A minimal CLI based app that automates keeping track of you actually do things on a daily basis.

Ttrack is built with these things in mind:
* Daemonless. Time tracking apps don't need to run in the background.
* Efficient. Ttrack usually writes to a single file cache to keep execution time minimal.
* Readable. Time tracked information is ultimately stored as human readable text files that play nice with version control.
* Unix Philosophy. Ttrack's scope is very limited, relying on other programs to take care of non time tracking tasks.
* Relatable. Ttrack is meant for people who want to set timed goals or view their trends on a daily basis according to their own timezone. Time is measured in hours/minutes/seconds, because most people don't operate at millisecond or nanosecond levels.

## Install
Download the latest release from the [release page](https://github.com/alanxoc3/ttrack/releases). At the moment, only Linux and Mac are supported.

If you have golang, you can install the latest updates from source. Keep in mind that the `main` branch is not guaranteed to be stable.
```bash
go install github.com/alanxoc3/ttrack
```

## Usage
Some time tracking utilities require you to manually start your time before you begin a task and end your time after you finish. This is unfortunately error prone. You could forget to start or stop the timer. Or maybe you get distracted or take a break while working on a task.

Ttrack instead works by inserting records and specifying a timeout with your record. If ttrack is called again before the duration of the timeout has been reached, then the time is extended. Otherwise, the time is added to your daily total in a human readable text file.

A common use case of ttrack is to call it from the event system of another program. For example, if you want to keep track of how much time you edit files, you could call ttrack's record command on your editor's key press event:

```bash
ttrack rec text-editor 5s
```

Ttrack records 5 seconds into the "text-editor" group. If the `rec` command is called again within 5 seconds, the time is extended. If not, only 5 seconds will end up getting added to the group.

If you want to see how much time you have used your text editor, use ttrack's aggregate command:

```bash
ttrack agg text-editor
```

Other commands help manage groups, clean the cache, and edit/create entries for other days. Run `ttrack help` to learn more.

## Integrations
Since ttrack is CLI based, it fits nicely into hooks that support shell calls. Here are some snippets that demonstrate how to make ttrack integrate with other applications. If you have another snippet, feel free to open a pull request with your addition.

### Tmux
[Tmux](https://github.com/tmux/tmux) is a popular terminal multiplexer. Add this to your `tmux.conf`:
```
set -g bell-action none
set -g visual-bell off
set -g monitor-activity on
set -g activity-action current
set-hook -g alert-activity "run-shell 'ttrack rec tmux 5s'"
```

### Kakoune
[Kakoune](https://kakoune.org/) is CLI based text editor. Add this to your `kakrc` or `autoload` directory:
```
hook global -group ttrack RawKey . %{
  evaluate-commands %sh{
    {
      [[ $(basename $kak_bufname) =~ '.' ]] && ttrack_name="ext:${kak_bufname##*.}" || ttrack_name="misc"
      [ ! -z "$(command -v ttrack)" ] && ttrack rec "kak/$ttrack_name" 3s
    } > /dev/null 2>&1 < /dev/null &
  }
}

hook global BufCreate .+\.tt %{
    remove-hooks global ttrack
}
```

### Concards
[Concards](https://github.com/alanxoc3/concards) is a CLI based flashcard program. Add this to your `event-startup` and `event-review` hooks:
```
#!/bin/bash
ttrack rec concards 30s
```

### Mpv
[Mpv](https://github.com/mpv-player/mpv) is a CLI based media player. Add this to a `~/.config/mpv/scripts/ttrack.lua` file:
```
mp.add_periodic_timer(1, function()
    os.execute("ttrack rec mpv 3s")
end)
```

## Comparison With Similar Apps
Here are some other time tracking applications and how ttrack relates to them:
* [Watson](https://tailordev.github.io/Watson/): Records time by starting and stopping instead of continually updating.
* [Gtm](https://github.com/laughedelic/gtm): Similar to watson, but meant for working in git.
* [ActivityWatch](https://github.com/ActivityWatch/activitywatch): Heavier & more complex, but supports more features and granularity.
* [WakaTime](https://wakatime.com/): Plugin oriented instead of cli based.
* [SelfSpy](https://github.com/selfspy/selfspy): A daemon that records keystrokes and is X11 specific.

## Dependencies
Ttrack has very few dependencies outside of go's standard library. Here are all the direct dependencies:
* [Cobra](https://github.com/spf13/cobra) for CLI arguments.
* [Testify](https://github.com/stretchr/testify) for unit tests.
* [BoltDb](https://github.com/etcd-io/bbolt) for implementing a cache.

## Internal Details
### Folder Structure
Ttrack uses one directory for data and one directory for a cache, following the XDG standard.

The data directory is calculated by the following order until one succeeds:
1. `$TTRACK_DATA_DIR`
2. `$XDG_DATA_HOME/ttrack`
3. `$HOME/.local/share/ttrack`
4. `./`

The cache directory is calculated by the following order until one succeeds:
1. `$TTRACK_CACHE_DIR`
2. `$XDG_CACHE_HOME/ttrack`
3. `$HOME/.cache/ttrack`
4. `./`

### Data Directory
Within the data directory, files ending in `.tt` are read from and written to. Files not ending in `.tt` are ignored. These data files have a very simple format that looks like this:
```
2021-01-01: 10m33s
2021-01-02: 1h2m3s
2021-01-06: 3h59m49s
2021-02-15: 30s
```

Because groups are stored as files, there are some restrictions as to what a group can be named. Invalid names will automatically be parsed to be valid. Here are the restrictions:
* Groups cannot end with `.tt`.
* Groups cannot start with a `.`.
* Groups cannot contain any whitespace as specified by [golang's unicode package][isspace].
* Groups cannot contain the `/` character. This is used to separate groups into a folder structure.

[isspace]: https://golang.org/pkg/unicode/#IsSpace

### Cache Directory
Within the cache directory, only one file named `db` is read from and written to. This file is a [BoltDB](https://github.com/etcd-io/bbolt) database file. Live time ttracking information is stored here to reduce the amount of writes to files in the data directory. The basic format of this database file is:
```
"<group-name>": {
  "out": 1m00s
  "beg": 2021-01-01T07:34:59Z
  "end": 2021-01-01T07:34:59Z
}
```
