#!/usr/bin/env bash

INSTANCES=${1}
CONTAINERS_PER_INSTANCE=${2}

for i in $(seq $INSTANCES); do
    go run . -n ${CONTAINERS_PER_INSTANCE} &
done

wait
