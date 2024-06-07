#!/bin/bash

# Get a reference to either podman or docker.
DOCKER=$(which podman || which docker || echo 'docker')

echo ">> Tearing down existing deployment..."
# Cleanup the previous deployment if present.
$DOCKER rm --force zoo1 zoo2 zoo3 || true
$DOCKER network rm --force zoonet || true


echo ">> Creating network..."
# All Zookeeper instances will run on this network.
$DOCKER network create zoonet --driver bridge
sleep 1

for n in {1..3};
do
    echo ">> Running zoo$n"
    $DOCKER run \
        --detach \
        --name "zoo$n" \
        --hostname "zoo$n" \
        --restart always \
        --network zoonet \
        --publish "218$n:2181" \
        --env "ZOO_MY_ID=$n" \
        --env 'ZOO_SERVERS=server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181' \
        zookeeper:3.9
    # Wait for a litle bit before running the next container.
    sleep 1
done
