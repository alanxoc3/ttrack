# ttrack
Cli based productivity time manager.

WIP.

Api idea:
```
ttrack help|version|rec|cp|del|add|output

ttrack 
ttrack --help
ttrack help

ttrack --version
ttrack version

ttrack rec <group-name> --timeout=1s|1m|1h|1d

ttrack cp <group-name> <group-name>

ttrack del <group-name>

ttrack del <group-name> --begin=2021-01-01 --end=2021-01-03

ttrack add <group-name> 2021-01-03T07:34:59Z--2021-01-03T07:34:59Z

ttrack output seconds|minutes|hours|days <group-name> --begin=2021-01-01 --end=2021-01-10
  10

ttrack output groups
  <group-name-1>
  <group-name-2>
  <group-name-3>

ttrack output timestamps <group-name> --begin=2021-01-01 --end=2021-01-10
  2021-01-01--2021-01-09
  2021-01-02--2021-01-09
  2021-01-03--2021-01-09
  2021-01-04--2021-01-09
  2021-01-05--2021-01-09
```

Database format idea:

```
{
  "<group-name>": {
    "timeout": 1m
    "start": 07:34:59
    "current": 07:39:59
    "record": {
      "2021-01-01":
        1000-3000
        4000-5000
        6000-9000
      "2021-01-02":
        1001-3003
        4001-5003
        6001-9003
    }
  }
}
```

No space character allowed in group names.
