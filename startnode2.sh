#!/bin/bash
PRIVATE_CONFIG=ignore nohup ./build/bin/geth \
--datadir node2 \
--nodiscover \
--verbosity 5 \
--networkid 10 \
--raft \
--raftport 50402 \
--rpc \
--rpcaddr 0.0.0.0 \
--rpcport 22002 \
--rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum,raft \
--emitcheckpoints \
--port 21002 \
>> node2.log 2>&1 &
