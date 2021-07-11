# ttrack

[![Build Status](https://travis-ci.com/alanxoc3/ttrack.svg?branch=main)](https://travis-ci.com/alanxoc3/ttrack)
[![Go Report Card](https://goreportcard.com/badge/github.com/alanxoc3/ttrack)](https://goreportcard.com/report/github.com/alanxoc3/ttrack)
[![Coverage Status](https://coveralls.io/repos/github/alanxoc3/ttrack/badge.svg?branch=main)](https://coveralls.io/github/alanxoc3/ttrack?branch=main)

A cli based app that helps you record how long you do things on a daily basis. This app is different from other time tracking programs in that it satisfies these requirements:
* CLI based for scripting.
* Focused on daily time tracking.
* Lightweight, fast, and no daemon.
* No need to worry about stopping your time.
* Store data in text files for version control.

## Install
Download the latest release from the [release
page](https://github.com/alanxoc3/concards/releases). At the moment, only Linux
and Mac are supported.

If you have golang, you can install the latest updates from source. Keep in mind that the `main` branch is not guaranteed to be stable.
```bash
go install github.com/alanxoc3/concards
```

## Usage
Some time tracking utilities require you to manually start your time before you begin a task and end your time after you finish. This is unfortunately error prone, because it's all too common to forget to start or end your time. Ttrack instead primarily works by inserting records and specifying a timeout with your record. If ttrack is called again before the duration of the timeout has been reached, then the time is extended. Otherwise, your time is added to your daily total in a text file.

A common use case of ttrack is to automate calling every time you interact with a program. As an example, you may want to call ttrack every time you input a key in your text editor to see how much time you spend editing files. You could put this command into a hook that gets run on key press:

```bash
ttrack rec text-editor 5s
```

And this command will tell you how much time you have been editing text per day.

```bash
ttrack agg text-editor
```

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
      [ ! -z "$(command -v ttrack)" ] && ttrack rec -- "kak/$kak_bufname" 5s
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
ttrack rec concards 30s
```

## Internal Details
Ttrack uses [BoltDb](https://github.com/etcd-io/bbolt) as cache. The "rec" command stores time tracking information in this file until the timeout has been reached, then the timestamp is added to a text file based on your group name.
The format of the bolt db file looks similar to, not taking serialization/marshaling into account:
```
{
  "<group-name>": {
    "out": 1m00s
    "beg": 2021-01-01T07:34:59Z
    "end": 2021-01-01T07:34:59Z
  }
}, ...
```

The text file format is also very simple. Here is an example of what it might look like:
```
2021-01-01 10m33s
2021-01-02 1h2m3s
2021-01-06 3h59m49s
2021-02-15 30s
...
```

## Comparison Similar Apps
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

## Development Notes
This project is still very much in alpha. Below is how I want the CLI api designed.

```
ttrack [--help] [--version] <command>

ttrack help
ttrack version
ttrack tidy

ttrack ls [<group>...] --recursive

ttrack rec <group> 1h10m30s

ttrack set <group> 2021-01-01 20m30s
ttrack add <group> 2021-01-01 20m30s
ttrack sub <group> 2021-01-01 20m30s

ttrack mv <group>... <group> --begin-date=2021-01-01 --end-date=2021-01-05
ttrack cp <group>... <group> --begin-date=2021-01-01 --end-date=2021-01-05

ttrack rm  <group>... --begin-date=2021-01-01 --end-date=2021-01-05

ttrack agg <group>... --begin-date=2021-01-01 --end-date=2021-01-05 --daily
```

Groups naming rules:
* Groups cannot end with `.tt`.
* Groups cannot start with a `.`.
* Groups cannot contain any whitespace as specified by [this page][isspace].
* Groups cannot contain the `/` character. This is used to separate groups into a folder structure.

[isspace]: https://golang.org/pkg/unicode/#IsSpace
