# Zoolection

Zoolection is a minimal implementation of leader-election in Go with the help of Apache Zookeeper.

## Instructions

- Running the zookeeper cluster:
    - With docker compose: `docker compose -f zookeeper.yml up -d`
    - With docker or podman: `bash zookeeper.sh`
- Run a few instances of the application using `go run .` across different terminals. One of them should get elected as the leader.
- Close the instance that is the leader and see another instance get elected.

