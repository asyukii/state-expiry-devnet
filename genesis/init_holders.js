const web3 = require("web3")
const init_holders = [
  {
     address: "0xC95b7a10ED53C3efe6Ce03e2A35527Bae0c1A3Cf",
     balance: web3.utils.toBN("500000000000000000000000000").toString("hex") // 500000000e18
  }
];


exports = module.exports = init_holders
