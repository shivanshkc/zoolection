#!/bin/bash

# Get a reference to either podman or docker.
DOCKER=$(which podman || which docker || echo 'docker')

# Cleanup the previous deployment if present.
$DOCKER rm --force zoo1 zoo2 zoo3 || true
$DOCKER network rm --force zoonet || true

# All Zookeeper instances will run on this network.
$DOCKER network create zoonet --driver bridge

# First container
$DOCKER run \
    --detach \
    --name zoo1 \
    --hostname zoo1 \
    --restart always \
    --network zoonet \
    --publish 2181:2181 \
    --env 'ZOO_MY_ID=1' \
    --env 'ZOO_SERVERS=server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181' \
    zookeeper:3.9

# Second container
$DOCKER run \
    --detach \
    --name zoo2 \
    --hostname zoo2 \
    --restart always \
    --network zoonet \
    --publish 2182:2181 \
    --env 'ZOO_MY_ID=2' \
    --env 'ZOO_SERVERS=server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181' \
    zookeeper:3.9

# Third container
$DOCKER run \
    --detach \
    --name zoo3 \
    --hostname zoo3 \
    --restart always \
    --network zoonet \
    --publish 2183:2181 \
    --env 'ZOO_MY_ID=3' \
    --env 'ZOO_SERVERS=server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181' \
    zookeeper:3.9