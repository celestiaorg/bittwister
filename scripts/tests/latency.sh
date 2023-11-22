#!/bin/bash

DOCKER_IMG="bittwister-lt-test"
CONTAINER_NAME="bittwister-lt-test"
CLIENT_CONTAINER_NAME="${CONTAINER_NAME}-client"
NETWORK_INTERFACE="eth0"
NUM_OF_PING_PACKETS=10

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${SCRIPT_DIR}/utils.sh"

function cleanup {
    echo "Cleaning up..."
    docker rm -f ${CONTAINER_NAME};
    docker rm -f ${CLIENT_CONTAINER_NAME};
    echo "Done"
}

function ping_test {
    ping_result=$(exec_on_container ${CLIENT_CONTAINER_NAME} "ping -c ${NUM_OF_PING_PACKETS} $1")
    avg_latency=$(echo "$ping_result" | awk '/min\/avg\/max/ {split($4, a, "/"); print a[2]}')
    echo -e "\nAverage latency: $avg_latency ms \t Expected: $2 ms\n"
}

# ----- main ------ #
echo "Running Latency tests..."
check_required_tools
build_image ${DOCKER_IMG}

# run another contianer to run ping client without any limitations
new_container ${DOCKER_IMG} ${CLIENT_CONTAINER_NAME} "start -d ${NETWORK_INTERFACE}"

allResults=""
# These are the expected latency in miliseconds
# Each iteration perfroms a test with a different latency
# to do so, it creates a new container with the expected latency
# then it runs ping client on the client container
# and finally it compares the expected latency with the actual latency
for EXPECTED_LATENCY in 10 50 100 200 500 700 1000 1500 2000 3000 5000; do
    echo -e "\nTesting with latency: ${EXPECTED_LATENCY} ms"
    new_container ${DOCKER_IMG} ${CONTAINER_NAME} "start -d ${NETWORK_INTERFACE} -l ${EXPECTED_LATENCY}"
    
    IP_ADDRESS=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CONTAINER_NAME})
    echo "Target IP address: ${IP_ADDRESS}"

    test_result=$(ping_test ${IP_ADDRESS} ${EXPECTED_LATENCY})
    allResults="${allResults}${test_result}"
    
    echo -e "\n${test_result}"
done;

echo -e "\n\nResults:"
echo "Number of packets per test: ${NUM_OF_PING_PACKETS}"
echo -e "${allResults}"
echo -e "\n\n"

cleanup

# We wait for the user to see the results before the next test starts
wait 5