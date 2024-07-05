const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying contracts with the account:", deployer.address);

    const Faucet = await ethers.getContractFactory("Faucet");
    const faucet = await Faucet.deploy();

    await faucet.waitForDeployment();

    const contractAddress = await faucet.getAddress()

    console.log("Contract deployed to address:", contractAddress);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });
