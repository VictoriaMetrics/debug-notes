Run with command

```
$ ./bin/vmagent -remoteWrite.url=http://homelab:8429/api/v1/write,http://homelab:8429/api/v1/write -promscrape.config=../debug-notes/gh10607/scrape.yaml -remoteWrite.relabelConfig=../debug-notes/gh10607/remote_write_relabel.yaml -remoteWrite.urlRelabelConfig=../debug-notes/gh10607/remote_write_url_relabel_1.yaml,../debug-notes/gh10607/remote_write_url_relabel_2.yaml
```
