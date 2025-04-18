# docs: improve drop labels documentation

Issue https://github.com/VictoriaMetrics/VictoriaMetrics/issues/8715

1. Checkout VictoriaMetrics repo

```
git clone git@github.com:VictoriaMetrics/VictoriaMetrics.git
git clone git@github.com:makasim/debug-notes.git

cd VictoriaMetrics 
```

2. Build vmagent and storage

```
make vmagent
make victoria-metrics

mkdir firststorage
mkdir secondstorage
```

Start two storages:

```
./bin/victoria-metrics -search.latencyOffset=0 -storageDataPath=firststorage
```

```
./bin/victoria-metrics -search.latencyOffset=0 -httpListenAddr=":8528" -storageDataPath=secondstore
```

4. Run cases (described in issue description) from root dir.

## Cases

Case 1:

OK. One `-dropInputLabels` is dropped to all remote writes

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_1:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_1:1s_max {bar="barVal",baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 1
```

Case 2:

OK. One `-dropInputLabels` with `^^` separated labels. Both applied to all remote writes

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo^^bar" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_2:1s_max {baz="bazVal"}
remote2: a_metric_case_2:1s_max {baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 2
```

Case 3:

OK. One `-dropInputLabels` with `,` separated labels. "foo" is dropped to remote1, bar is dropped to remote2

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo,bar" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_3:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_3:1s_max {baz="bazVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 3
```

Case 4:

OK. One `-dropInputLabels` with empty label before ",". all labels go to remote1, bar is dropped to remote2

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels=",bar" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_4:1s_max {bar="barVal",baz="bazVal",foo="fooVal"}
remote2: a_metric_case_4:1s_max {baz="bazVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 4
```

Case 5:

OK. One `-dropInputLabels` with "," and "^^". "foo, bar" labels are dropped from remote1, "foo, baz" are dropped from remote2

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo^^bar,foo^^baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_5:1s_max {baz="bazVal"}
remote2: a_metric_case_5:1s_max {bar="barVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 5
```

Case 6:

OK. Two `-dropInputLabels` flags, "baz" is dropped from remote1, "foo" is dropped from remote2

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.streamAggr.dropInputLabels="baz" \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_6:1s_max {bar="barVal",foo="fooVal"}
remote2: a_metric_case_6:1s_max {bar="barVal",baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 6
```

Case 7:

Confusing. Two `-dropInputLabels` flags, but values from the first one are used, "foo" is dropped from remote1, "baz" is dropped from remote2

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.streamAggr.dropInputLabels="foo,baz" \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_7:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_7:1s_max {bar="barVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 7
```

Case 8:

Confusing. Two `-dropInputLabels` flags, but values from the first one are used for both remote writes, "foo" is dropped from remote1, "foo" is dropped from remote2

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.streamAggr.dropInputLabels="foo" \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="foo,baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_8:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_8:1s_max {bar="barVal",baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 8
```

Case 10:

OK. Drops "foo" label from both remote writes

```
./bin/vmagent  \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write   \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write   \
  -streamAggr.dropInputLabels="foo"  \
  -streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_10:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_10:1s_max {bar="barVal",baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 10
```

Case 11:

Confusing, expecting "foo" is dropped from remote1, "bar" is dropped from remote2, but it works like foo^^bar

```
./bin/vmagent  \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write   \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write   \
  -streamAggr.dropInputLabels="foo,bar"  \
  -streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_11:1s_max {baz="bazVal"}
remote2: a_metric_case_11:1s_max {baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 11
```

Case 12:

Confusing, "^^" does not work at all for streamAggr.

```
./bin/vmagent  \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write   \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write   \
  -streamAggr.dropInputLabels="foo^^bar"  \
  -streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_12:1s_max {bar="barVal",baz="bazVal",foo="fooVal"}
remote2: a_metric_case_12:1s_max {bar="barVal",baz="bazVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 12
```

Case 13:

OK. Drop labels defined in config takes precedence over command line `-streamAggr.dropInputLabels` flag.

```
./bin/vmagent  \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write   \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write   \
  -streamAggr.dropInputLabels="bar,baz"  \
  -streamAggr.config=../debug-notes/gh8715/aggrcfg_drop_foo.yaml

remote1: a_metric_case_13:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_13:1s_max {bar="barVal",baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 13
```

Case 14:

OK. Drop labels defined in config takes precedence over command line `-remoteWrite.streamAggr.dropInputLabels` flag.

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="bar,baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg_drop_foo.yaml

remote1: a_metric_case_14:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_14:1s_max {bar="barVal",baz="bazVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 14
```

Case 15:

OK. Drop labels defined in config takes precedence over several command line `-remoteWrite.streamAggr.dropInputLabels` flags.

```
./bin/vmagent \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.streamAggr.dropInputLabels="baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg_drop_foo.yaml \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg_drop_bar.yaml

remote1: a_metric_case_15:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_15:1s_max {baz="bazVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 15
```

Case 16:

OK. One defined flag works for all remote writes.

```
./bin/vmagent \
  -remoteWrite.streamAggr.dropInputLabels="baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write

remote1: a_metric_case_16:1s_max {bar="barVal",foo="fooVal"}
remote2: a_metric_case_16:1s_max {bar="barVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 16
```

Case 17:

Confusing. I'd expect baz be dropped from remote2, but instead it is bar (the second flag, which defined before remote write).

```
./bin/vmagent \
  -remoteWrite.streamAggr.dropInputLabels="foo" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml \
  -remoteWrite.url=http://0.0.0.0:8428/api/v1/write  \
  -remoteWrite.streamAggr.dropInputLabels="bar" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml \
  -remoteWrite.url=http://0.0.0.0:8528/api/v1/write \
  -remoteWrite.streamAggr.dropInputLabels="baz" \
  -remoteWrite.streamAggr.config=../debug-notes/gh8715/aggrcfg.yaml

remote1: a_metric_case_17:1s_max {bar="barVal",baz="bazVal"}
remote2: a_metric_case_17:1s_max {baz="bazVal",foo="fooVal"}
```

```
go run ../debug-notes/gh8715/simulate_push_metrics.go 18
```

Open VMUI (http://localhost:8428/vmui/, http://localhost:8528/vmui) and query metric `{__name__="a_metric_.*"}`: