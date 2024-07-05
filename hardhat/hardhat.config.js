require("@nomicfoundation/hardhat-toolbox");
require('dotenv').config();

const { PRIVATE_KEY } = process.env;

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.24",
  networks: {
    local: {
      url: "http://127.0.0.1:8545",
      accounts: [PRIVATE_KEY]
    }
  }
};
