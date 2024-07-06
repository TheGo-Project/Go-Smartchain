const { ethers } = require('ethers');
const fs = require('fs');
const path = require('path');
require('dotenv').config();

// Connect to the local Geth node
const provider = new ethers.JsonRpcProvider('http://127.0.0.1:8545');
const privateKey = process.env.PRIVATE_KEY;
const wallet = new ethers.Wallet(privateKey, provider);

async function main() {
    const network = await provider.getNetwork();
    console.log("Current network:", network.name);

    // Read the compiled contract
    const contractPath = path.join(__dirname, '..', 'artifacts', 'contracts', 'faucet.sol', 'Faucet.json');
    const contractJSON = JSON.parse(fs.readFileSync(contractPath));
    const abi = contractJSON.abi;
    const bytecode = contractJSON.bytecode;

    // Create a Contract Factory
    const ContractFactory = new ethers.ContractFactory(abi, bytecode, wallet);

    const feeData = await provider.getFeeData()
    console.log("feeData:", ethers.formatUnits(feeData.maxFeePerGas, "ether"))

    // Deploy the contract
    console.log('Deploying contract...');
    const contract = await ContractFactory.deploy();

    // Wait for the contract to be mined
    // await contract;
    console.log('Contract deployed to address:', await contract.getAddress());
}

main().catch(error => {
    console.error('Error deploying contract:', error);
    process.exit(1);
});

