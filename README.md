# ttrack
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

ttrack ls  [<group>...] --recursive --escape
ttrack rec <group>... 10m30s
ttrack set <group>... 2021-01-01:+20m30s

ttrack mv  <group>... <group> --begin-date=2021-01-01 --end-date=2021-01-05 --recursive
ttrack cp  <group>... <group> --begin-date=2021-01-01 --end-date=2021-01-05 --recursive

ttrack del <group>... --begin-date=2021-01-01 --end-date=2021-01-05 --recursive
ttrack agg <group>... --begin-date=2021-01-01 --end-date=2021-01-05 --recursive --by-day
```

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

Group name restraints:
- No space character allowed in group names.
