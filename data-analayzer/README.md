# Data analyzer

Helps to analyze data exported via /api/v1/export API.
Example of the data for [#8807](https://github.com/VictoriaMetrics/VictoriaMetrics/issues/8807#issuecomment-2833295291)
is in `gh8807.json`.

To run:
```sh
go run .

......

2025-04-18 07:50:38 +0200 CEST 7617979 (step: 15000, value change: 1102)
2025-04-18 07:50:53 +0200 CEST 7619296 (step: 15000, value change: 1317)
2025-04-18 07:51:08 +0200 CEST 7620647 (step: 15000, value change: 1351)
2025-04-18 07:51:23 +0200 CEST 7622120 (step: 15000, value change: 1473)
2025/05/08 16:50:03 total samples 44
2025/05/08 16:50:03 duplicates 0
2025/05/08 16:50:03 avg step 14 sec
2025/05/08 16:50:03 min date 2025-04-18 07:40:38 +0200 CEST
2025/05/08 16:50:03 min date 2025-04-18 07:51:23 +0200 CEST
```