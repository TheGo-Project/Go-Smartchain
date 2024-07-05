#!/bin/bash

#geth --http --http.corsdomain="https://remix.ethereum.org" --http.api web3,eth,debug,personal,net --vmdebug --datadir geth_node --dev --http.port 8545 &
#GETH_PID=$!
#
#./main &
#MAIN_PID=$!
#
#wait_for_process() {
#    local pid=$1
#    while kill -0 $pid > /dev/null 2>&1; do
#        sleep 1
#    done
#}
#
#wait_for_process $GETH_PID
#wait_for_process $MAIN_PID
#
#exit 0


geth --http --http.corsdomain="https://remix.ethereum.org" --http.api web3,eth,debug,personal,net --vmdebug --datadir geth_node --dev --http.port 8545 &
GETH_PID=$!

./main &
MAIN_PID=$!

wait $GETH_PID || wait $MAIN_PID

# Exit with the status of the process that exited first
exit $?
