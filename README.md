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
Api idea:
```
ttrack 
ttrack --help
ttrack help
ttrack version

ttrack groups
  <group-1>
  <group-2>

ttrack rec <group> 10m30s
ttrack set <group> 2021-01-03 10m30s
ttrack timer <group> 10m30s
ttrack del <group>
ttrack cp <group> <group> --begin-date=2021-01-01 --end-date=2021-01-03

ttrack list <group> --begin-date=2021-01-01 --end-date=2021-01-04
  2021-01-01: 10m
  2021-01-02: 13s
  2021-01-03: 30m
  2021-01-04: 1h

ttrack agg <group> --begin-date=2021-01-01 --end-date=2021-01-04
  1h40m13s
```

Bolt db format:
```
{
  "<group-name>": {
    "out": 1m 00s
    "beg": 2021-01-01T07:34:59Z
    "end": 2021-01-01T07:34:59Z
    "rec": {
      "2021-01-01": 10m 35s
      "2021-01-02": 30m 35s
    }
  }
}
```

Group name restraints:
- No space character allowed in group names.
