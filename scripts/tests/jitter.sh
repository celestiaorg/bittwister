#!/bin/bash

DOCKER_IMG="bittwister-jt-test"
CONTAINER_NAME="bittwister-jt-test"
CLIENT_CONTAINER_NAME="${CONTAINER_NAME}-client"
NETWORK_INTERFACE="eth0"
NUM_OF_PING_PACKETS=50

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "${SCRIPT_DIR}/utils.sh"

function cleanup {
    echo "Cleaning up..."
    docker rm -f ${CONTAINER_NAME};
    docker rm -f ${CLIENT_CONTAINER_NAME};
    echo "Done"
}

function ping_test {
    ping_result=$(exec_on_container_noit ${CLIENT_CONTAINER_NAME} "ping -c ${NUM_OF_PING_PACKETS} $1")
    rtt_times=($(echo "$ping_result" | grep -oE 'time=[0-9.]+ ms' | grep -oE '[0-9.]+'))

    # Calculate jitter
    jitter_sum=0
    prev_rtt=${rtt_times[0]}

    for ((i = 1; i < ${#rtt_times[@]}; i++))
    do
        current_rtt=${rtt_times[$i]}
        diff=$(echo "$current_rtt - $prev_rtt" | bc -l | awk '{print sqrt($1^2)}')
        jitter_sum=$(echo "$jitter_sum + $diff" | bc -l)
        prev_rtt=$current_rtt
    done

    # Calculate average jitter
    num_packets=${#rtt_times[@]}
    average_jitter=$(echo "scale=3; $jitter_sum / ($num_packets - 1)" | bc -l)

    echo -e "\nAverage jitter: $average_jitter ms \t Max expected jitter: $2 ms\n"
}

# ----- main ------ #
echo "Running jitter tests..."
check_required_tools
build_image ${DOCKER_IMG}

# run another contianer to run ping client without any limitations
new_container ${DOCKER_IMG} ${CLIENT_CONTAINER_NAME} "start -d ${NETWORK_INTERFACE}"

allResults=""
# These are the expected max jitter in miliseconds
# Each iteration perfroms a test with a different jitter
# to do so, it creates a new container with the expected jitter
# then it runs ping client on the client container
# and finally it compares the expected jitter with the actual jitter
for MAX_JITTER in 10 50 100 200 500 700 1000 1500 2000 3000 5000; do
    echo -e "\nTesting with jitter: ${MAX_JITTER} ms"
    new_container ${DOCKER_IMG} ${CONTAINER_NAME} "start -d ${NETWORK_INTERFACE} -j ${MAX_JITTER}"
    
    IP_ADDRESS=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ${CONTAINER_NAME})
    echo "Target IP address: ${IP_ADDRESS}"

    test_result=$(ping_test ${IP_ADDRESS} ${MAX_JITTER})
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