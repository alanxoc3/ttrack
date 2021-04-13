# ttrack
Cli based productivity time manager.

WIP.

Api idea:
```
ttrack 
ttrack --help
ttrack help
ttrack version

ttrack rec <group-name> --timeout=1d1h1m1s
ttrack set <group-name> 2021-01-03 10m30s

ttrack mv <group>... <group>
ttrack del <group> --begin-date=2021-01-01 --end-date=2021-01-03

ttrack groups
  <group-name-1>
  <group-name-2>
  ...

ttrack list <group> --begin-date=2021-01-01 --end-date=2021-01-04
  2021-01-01: 10m
  2021-01-02: 13s
  2021-01-03: 30m
  2021-01-04: 1h

ttrack agg <group> --begin-date=2021-01-01 --end-date=2021-01-04
  1h40m13s
  0s
```

Config file:
```
default_timeout=30s
```

Database format idea:
```
{
  "<group-name>": {
    "timeout": 1m 00s
    "start":   2021-01-01T07:34:59Z
    "current": 2021-01-01T07:34:59Z
    "record": {
      "2021-01-01": 10m 35s
      "2021-01-02": 30m 35s
    }
  }
}
```

Group name restraints:
- No space character allowed in group names.
