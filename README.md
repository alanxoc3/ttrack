# ttrack
Cli based productivity time manager.

WIP.

Api idea:
```
ttrack 
ttrack --help
ttrack help

ttrack --version
ttrack version

ttrack rec <group-name> --timeout=1d1h1m1s
ttrack add <group-name> 2021-01-03 10m30s

ttrack mv <group-name>... <group-name>

ttrack del <group-name> --begin=2021-01-01 --end=2021-01-03
ttrack sub <group-name> --date=2021-01-03 10m30s

ttrack list
  <group-name-1>
  <group-name-2>
  ...

ttrack format <group-name-1> <group-name-2>... --format='%g: %d %h %m %s' --begin=2021-01-01 --end=2021-01-10
```

Config file:
```
--timeout=10m
```

Format meanings:
```
%g = group name
%d = number of days.
%h = number of hours.
%m = number of minutes.
%s = number of seconds.

%D = total time in days.
%H = total time in hours.
%M = total time in minutes.
%S = total time in seconds.
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
