// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
const hre = require("hardhat");
const fs = require("fs");

async function main() {
  const ERC20Token = await hre.ethers.getContractFactory("ERC20Token");
  const erc20Token = await ERC20Token.deploy();

  await erc20Token.deployed();

  console.log(`ERC20Token deployed to ${erc20Token.address}`, "tx hash", erc20Token.deployTransaction.hash);

  // Check if the JSON file exists
  jsonPath = "../deployed_contracts.json";
  let data = {};
  if (fs.existsSync(jsonPath)) {
    // Read the existing deployed contracts from the JSON file
    const fileContents = fs.readFileSync(jsonPath, "utf8");
    data = JSON.parse(fileContents);
  }

  // Update the deployed address of ABCToken in the data object
  data.ERC20Token = erc20Token.address;

  // Write the updated data object to the JSON file
  fs.writeFileSync(jsonPath, JSON.stringify(data, null, 2));
}


// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
