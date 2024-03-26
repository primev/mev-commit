require("@nomicfoundation/hardhat-toolbox");
require("@nomicfoundation/hardhat-foundry");
require("solidity-docgen");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.15",
  networks: {
    localhost: {
      url: "http://127.0.0.1:8545", // Ganache default port
    },
      aws: {
          url: "http://34.213.237.94:8123",
             chainId: 1001,
      accounts: ['0x28b2b0318721be8c8339199172cd7cc8f5e273800a35616ec893083a4b32c02e'] //account private key
      }
  },
  docgen: {
    pages: "files",
    exclude: ["contracts/interfaces", "lib"],
  },
};
