# Zoolection

Zoolection is a minimal implementation of leader-election in Go with the help of Apache Zookeeper.

## Instructions

- Run `docker compose -f zookeeper.yml up -d` to run a local Zookeeper cluster.
- Run a few instances of the application using `go run .` across different terminals. One of them should get elected as the leader.
- Close the instance that is the leader and see another instance get elected.

## Navigating the code

The whole Zookeeper implementation lives in the `pkg/election/zookeeper.go` file. See the `main.go` file for its usage.
