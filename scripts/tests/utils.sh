#!/bin/bash

function build_image {
    echo "Building test image..."
    docker build --target test -t ${DOCKER_IMG} .
}

function run_container {
    if docker ps -a --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"; then
        echo "Container ${CONTAINER_NAME} exists. Removing..."
        docker rm -f ${CONTAINER_NAME};
    fi

    echo "Running container..."
    cmd=${1}
    docker run --privileged --name ${CONTAINER_NAME} ${DOCKER_IMG} ${cmd} &

    # Check if the container is running
    while [[ "$(docker inspect -f '{{.State.Running}}' ${CONTAINER_NAME} 2>/dev/null)" != "true" ]]; do
        echo "Container ${CONTAINER_NAME} is not yet running. Waiting..."
        sleep 1
    done
}