#!/bin/bash

function build_image {
    # Function to build the docker image
    # Input: Image name

    local img=$1
    echo "Building test image..."
    docker build --target test -t ${img} .
}

function new_container {
    # Function to create a new container and run a command
    # Inputs: image name, container name, command to run

    container_name=${2}
    if docker ps -a --format "{{.Names}}" | grep -q "^${container_name}$"; then
        echo "Container ${container_name} exists. Removing..."
        docker rm -f ${container_name};
    fi

    echo "Running container..."
    local img=$1
    cmd=${3}
    docker run --privileged --name ${container_name} ${img} ${cmd} &

    # Check if the container is running
    while [[ "$(docker inspect -f '{{.State.Running}}' ${container_name} 2>/dev/null)" != "true" ]]; do
        echo "Container ${container_name} is not yet running. Waiting..."
        sleep 1
    done
}

function exec_on_container {
    # Function to execute a command on a running container
    # Inputs: contianer name, command to execute

    container_name=${1}
    cmd=${2}
    docker exec -it ${container_name} ${cmd}
}

function convert_bandwidth {
    # Function to convert bandwidth in bits per second to human-readable format
    # Input: Bandwidth in bits per second

    local bps=$1
    local unit=('bps' 'Kbps' 'Mbps' 'Gbps')
    local idx=0

    while (( bps > 1024 )) && (( idx < ${#unit[@]} )); do
        bps=$(( bps / 1024 ))
        (( idx++ ))
    done

    echo "${bps} ${unit[$idx]}"
}

function wait {
    # Function to wait for a given number of seconds
    # Input: Number of seconds to wait

    local secs=$1
    while (( secs > 0 )); do
        echo -ne "Waiting ${secs} seconds...\033[0K\r"
        sleep 1
        (( secs-- ))
    done
}