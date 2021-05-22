# ttrack

[![Build Status](https://travis-ci.com/alanxoc3/ttrack.svg?branch=main)](https://travis-ci.com/alanxoc3/ttrack)
[![Go Report Card](https://goreportcard.com/badge/github.com/alanxoc3/ttrack)](https://goreportcard.com/report/github.com/alanxoc3/ttrack)
[![Coverage Status](https://coveralls.io/repos/github/alanxoc3/ttrack/badge.svg?branch=main)](https://coveralls.io/github/alanxoc3/ttrack?branch=main)

A cli based app that helps you record how long you do things on a daily basis. This app is different from other time tracking programs in that it satisfies these requirements:
- CLI based for scripting.
- Focused on daily time tracking.
- Lightweight, fast, and no daemon.
- No need to worry about stopping your time.
- Store data in text files for version control.



## Integrations
### Kakoune
Add this to your `.kakrc`:
```
hook global RawKey . %{
  evaluate-commands %sh{
    {
      [ ! -z "$(command -v ttrack)" ] && ttrack rec -- "kak:$kak_bufname" 5s
    } > /dev/null 2>&1 < /dev/null &
  }
}
```

## Notes
Groups naming rules:
- `.tt` will be removed if at end of group/subgroup.
- Name cannot start with a `.`.
- Slashes denote sub groups.
- Timeout, begin ts, & end ts all go in cache.

Api idea:
```
ttrack [--help] [--version] <command>

ttrack help
ttrack version
ttrack tidy

ttrack ls [<group>...] --recursive --quote

ttrack rec <group> 1h10m30s
ttrack set <group> 2021-01-01 20m30s
ttrack add <group> 2021-01-01 20m30s
ttrack sub <group> 2021-01-01 20m30s

ttrack mv  <group>... <group> --begin-date=2021-01-01 --end-date=2021-01-05 --recursive
ttrack cp  <group>... <group> --begin-date=2021-01-01 --end-date=2021-01-05 --recursive
ttrack rm  <group>... --begin-date=2021-01-01 --end-date=2021-01-05 --recursive
ttrack agg <group>... --begin-date=2021-01-01 --end-date=2021-01-05 --recursive --daily
```

What about timers? Two types of timers:
- Wait until you hit a certain amount of time, then quit. So really like a wait/watch command.
- Set a timer and update the db/storage with the amount of time in that timer. If you ctrl-c out of it, it will do some of the timer.


## Internal Details
Ttrack uses [BoltDb](https://github.com/etcd-io/bbolt) as cache. The "rec" command stores time tracking information in this file until the timeout has been reached. At that pointa

stores time tracking instores timestamp upUpdates are stored in this file until the timeout has been reached.
Bolt db format:
```
{
  "<group-name>": {
    "out": 1m 00s
    "beg": 2021-01-01T07:34:59Z
    "end": 2021-01-01T07:34:59Z
  }
}, ...
```

File format:
```
2021-01-01 10m30s
```

## More Time Tracking Software
If this program doesn't satisfy your needs, there are other options you may be interested in that have a different focus from concards:
* [Watson](https://tailordev.github.io/Watson/): Records time by starting and stopping instead of continually updating.
* [Gtm](https://github.com/laughedelic/gtm): Similar to watson, but meant for working in git.
* [ActivityWatch](https://github.com/ActivityWatch/activitywatch): Heavier & more complex, but supports more features and granularity.
* [WakaTime](https://wakatime.com/): Plugin oriented instead of cli based.
* [SelfSpy](https://github.com/selfspy/selfspy): A daemon that records keystrokes and is X11 specific.
