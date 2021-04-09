# ttrack
Cli based productivity time manager.

WIP.

Api idea:
```
ttrack rec <group-name> --timeout=1m
ttrack cp <group-name>... <group-name> # Timeout not copied.
ttrack aggregate --begin=2021-01-01 --end=2021-01-03 --group=<group-name>

ttrack del <group-name>
ttrack del <group-name> --begin=2021-01-01 --end=2021-01-03

ttrack read-file file.json
ttrack read-file file.yaml

ttrack add-timestamp <group-name> 2021-01-03T07:34:59Z--2021-01-03T07:34:59Z

ttrack output <group-name> --format=json
ttrack read <group-name> --format=yaml

ttrack filter-timestamps <group-name> --begin=2021-01-01 --end=2021-01-10
  2021-01-01--2021-01-09
  2021-01-02--2021-01-09
  2021-01-03--2021-01-09
  2021-01-04--2021-01-09
  2021-01-05--2021-01-09

ttrack agg-time <group-name> --begin=2021-01-01 --end=2021-01-10
  2021-01-01--2021-01-09
  2021-01-02--2021-01-09
  2021-01-03--2021-01-09
  2021-01-04--2021-01-09
  2021-01-05--2021-01-09


ttrack groups
  <group-name-1>
  <group-name-2>
  <group-name-3>

ttrack --help | ttrack help

ttrack rename
ttrack combine

ttrack 
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
