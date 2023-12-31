#!/bin/bash

DOCKER_IMG="bittwister-pk-test"
CONTAINER_NAME="bittwister-pk-test"
NUM_OF_PING_PACKETS=50
NETWORK_INTERFACE="eth0"

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${SCRIPT_DIR}/utils.sh"

function cleanup {
    echo "Cleaning up..."
    docker rm -f ${CONTAINER_NAME};
    echo "Done"
}

function ping_test {
    ping_result=$(ping -c ${NUM_OF_PING_PACKETS} $1)
    packet_loss=$(echo "$ping_result" | awk '/packet loss/ {print $6}')

    echo -e "\nPacket loss: $packet_loss \t Expected: $2%\n"
}

# ----- main ------ #
echo "Running packetloss tests..."
check_required_tools
build_image ${DOCKER_IMG}

allResults=""
for EXPECTED_PACKET_LOSS in 0 10 20 50 70 80 100; do
    echo -e "\nTesting with packet loss: ${EXPECTED_PACKET_LOSS}%"
    new_container ${DOCKER_IMG} ${CONTAINER_NAME} "start -d ${NETWORK_INTERFACE} -p ${EXPECTED_PACKET_LOSS}"
    
    IP_ADDRESS=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CONTAINER_NAME})
    echo "Pinging IP address: ${IP_ADDRESS}"
    result=$(ping_test ${IP_ADDRESS} ${EXPECTED_PACKET_LOSS})
    allResults="${allResults}${result}"
    echo ${result}
done;

echo -e "\n\nResults:"
echo "Number of packets per test: ${NUM_OF_PING_PACKETS}"
echo -e "${allResults}"
echo -e "\n\n"

cleanup

# We wait for the user to see the results before the next test starts
wait 5