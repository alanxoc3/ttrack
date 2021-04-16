# ttrack
Cli based productivity time manager.

WIP.

Api idea:
```
ttrack 
ttrack --help
ttrack help
ttrack version

ttrack rec <group> 10m30s
ttrack timer <group> 10m30s
ttrack set <group> 2021-01-03 10m30s

ttrack mv <group> <group>
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
```

Database format idea:
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
