#!/bin/bash
PRIVATE_CONFIG=ignore nohup ./build/bin/geth \
--datadir node1 \
--nodiscover \
--verbosity 5 \
--networkid 10 \
--raft \
--raftport 50401 \
--rpc \
--rpcaddr 0.0.0.0 \
--rpcport 22001 \
--rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum,raft \
--emitcheckpoints \
--port 21001 \
>> node1.log 2>&1 &
