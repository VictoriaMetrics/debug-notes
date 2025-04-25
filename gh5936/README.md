# Query stats playground

Initially used in https://github.com/VictoriaMetrics/VictoriaMetrics-enterprise/pull/865

To run the playground execute the following command:
```bash
docker compose -f compose.yml up
```

It will spin-up the following components:
* VictoriaMetrics cluster with custom image of vmselect that has the query stats 
  enabled. _Once this feature is released - just use the latest image version._
* vmalert that generates the query load, so vmselect could print query stats logs
* vector to collect the generated stats logs from vmselects
* victorialogs to store logs collected by vector
* grafana to display `provisioning/query-stats.json` dashboard

Once compose started, visit [http://localhost:3000](http://localhost:3000) (admin:admin)
and open Query Stats dashboard.