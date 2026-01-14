1. Checkout PR https://github.com/VictoriaMetrics/VictoriaMetrics/pull/10274/files
2. Run `./setup.sh`. It will compile all needed binaries
3. Run storage `./run-vmsingle.sh`
4. Run `go run client.go`
5. Run two copies of agent `./run-vmagent.sh`

Client start failing requests after 5th scrape for one of the agents, while responding OK to the other agent.
In 10 seconds it will stop failing requests respond OK to both agents and cycle repeats.

