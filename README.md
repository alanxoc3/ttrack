# ttrack

[![Build Status](https://travis-ci.org/alanxoc3/ttrack.svg?branch=master)](https://travis-ci.org/alanxoc3/ttrack)
[![Go Report Card](https://goreportcard.com/badge/github.com/alanxoc3/ttrack)](https://goreportcard.com/report/github.com/alanxoc3/ttrack)
[![Coverage Status](https://coveralls.io/repos/github/alanxoc3/ttrack/badge.svg?branch=main)](https://coveralls.io/github/alanxoc3/ttrack?branch=main)

A cli based app meant to record how long you do things on a daily basis.

TODO: Explain how this compares to other time tracking programs:
- [Watson](https://tailordev.github.io/Watson/)
- [Gtm](https://github.com/laughedelic/gtm)
- [ActivityWatch](https://github.com/ActivityWatch/activitywatch)
- [WakaTime](https://wakatime.com/)
- [SelfSpy](https://github.com/selfspy/selfspy)

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

Bolt db format:
```
{
  "<group-name>": {
    "out": 1m 00s
    "beg": 2021-01-01T07:34:59Z
    "end": 2021-01-01T07:34:59Z
  }
}
```

File format:
```
2021-01-01 10m30s
```
