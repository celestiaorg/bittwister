#!/bin/bash

DOCKER_IMG="bittwister-pk-test"
CONTAINER_NAME="bittwister-pk-test"
NUM_OF_PING_PACKETS=50

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

build_image

allResults=""
for EXPECTED_PACKET_LOSS in 0 10 20 50 70 80 100; do
    echo -e "\nTesting with packet loss: ${EXPECTED_PACKET_LOSS}%"
    run_container "-p ${EXPECTED_PACKET_LOSS}"
    
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