#!/bin/bash

# Start the first process
geth --http --http.corsdomain="https://remix.ethereum.org" --http.api web3,eth,debug,personal,net --vmdebug --datadir geth_node --dev --http.port 8545 &

# Start the second process
./main &

# Wait for any process to exit
wait -n

# Exit with status of process that exited first
exit $?