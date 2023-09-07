require("@nomicfoundation/hardhat-toolbox");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.0",
  defaultNetwork: "devnet",
  networks: {
    hardhat: {
      chainId: 123454321,
    },
    devnet: {
      url: `http://127.0.0.1:8503`,
      accounts: ["190e410a96c56dcc7cbe6ee04ce68fbcf2eb7d86c441e840235373078cf6bb0c"],
      chainId: 123454321,
    },
  }
};
