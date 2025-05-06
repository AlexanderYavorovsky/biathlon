# Biathlon
## Task
Task is described in `testdata/README.md`.

## Prerequisites
- Go (1.20+)
- Optional: `make` `golangci-lint(v2)`

## How to run
### 1. Run using make
Default: (runs program with config1.json and events1)
```
make run
```

Customizable:
```
make run CONFIG_PATH=<path_to_config> EVENTS=<path_to_events>
```
Other make commands:
- `make build`
- `make test`
- `make cover` (show test coverage)
- `make clean`
- `make lint` (run linters)
- `make fmt` (run formatters)



### 2. Run 'with your hands'
```
CONFIG_PATH=<path_to_config.json> go run . < <path_to_events_file>
```

Example:
```
CONFIG_PATH=testdata/config1.json go run . < testdata/events1
```

Example output:
```
[09:05:59.867] The competitor (1) registered
[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000
[09:29:45.734] The competitor(1) is on the start line
[09:30:01.005] The competitor(1) has started
[09:49:31.659] The competitor(1) is on the firing range(1)
[09:49:33.123] The target(1) has been hit by competitor(1)
[09:49:34.650] The target(2) has been hit by competitor(1)
[09:49:35.937] The target(4) has been hit by competitor(1)
[09:49:37.364] The target(5) has been hit by competitor(1)
[09:49:38.339] The competitor(1) left the firing range
[09:49:55.915] The competitor(1) entered the penalty laps
[09:51:48.391] The competitor(1) left the penalty laps
[09:59:03.872] The competitor(1) ended the main lap
[09:59:03.872] The competitor(1) can`t continue: Lost in the forest

Final Report:
[NotFinished] 1 [{00:29:03.872, 2.094}, {,}] {00:01:52.476, 0.445} 4/5
```