#!/bin/bash

DOCKER_IMG="bittwister-bw-test"
CONTAINER_NAME="bittwister-bw-test"
CLIENT_CONTAINER_NAME="${CONTAINER_NAME}-client"
NETWORK_INTERFACE="eth0"

PARALLEL_CONNECTIONS=100
# in seconds
DURATION_PER_TEST=60

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${SCRIPT_DIR}/utils.sh"

function cleanup {
    echo "Cleaning up..."
    docker rm -f ${CONTAINER_NAME};
    docker rm -f ${CLIENT_CONTAINER_NAME};
    echo "Done"
}

# ----- main ------ #
echo "Running Bandwidth tests..."
check_required_tools
build_image ${DOCKER_IMG}

# run another contianer to run iperf3 client without any limitations
new_container ${DOCKER_IMG} ${CLIENT_CONTAINER_NAME} "start -d ${NETWORK_INTERFACE}"

allResults=""
for EXPECTED_BANDWIDTH in 65536 131072 262144 524288 1048576 2097152 4194304 8388608 16777216 33554432 67108864 134217728 268435456 536870912 1073741824; do
    echo -e "\nTesting with bandwidth: ${EXPECTED_BANDWIDTH} bps"
    new_container ${DOCKER_IMG} ${CONTAINER_NAME} "start -d ${NETWORK_INTERFACE} -b ${EXPECTED_BANDWIDTH}"
    
    # Running iperf3 server in daemon mode
    exec_on_container ${CONTAINER_NAME} "iperf3 -s -D"

    IP_ADDRESS=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CONTAINER_NAME})
    echo "Target IP address: ${IP_ADDRESS}"

    test_result=$(exec_on_container ${CLIENT_CONTAINER_NAME} "iperf3 -c ${IP_ADDRESS} -t ${DURATION_PER_TEST} -P ${PARALLEL_CONNECTIONS} --json")
    receiver=$(echo -e "$test_result" | jq '.end.sum_received.bits_per_second')
    
    converted_expected_bandwidth=$(convert_bandwidth $EXPECTED_BANDWIDTH)
    converted_receiver_bandwidth=$(convert_bandwidth $receiver)
    txt_output="expected bandwidth: ${converted_expected_bandwidth} \tactual bandwidth: ${converted_receiver_bandwidth}\n"

    allResults="${allResults}${txt_output}"
    echo -e "\n${txt_output}"
done;

echo -e "\n\nResults:"
echo "Number of parallel connections per test: ${PARALLEL_CONNECTIONS}"
echo "Test duration per test: ${DURATION_PER_TEST} seconds"
echo -e "${allResults}"
echo -e "\n\n"

cleanup

# We wait for the user to see the results before the next test starts
wait 5