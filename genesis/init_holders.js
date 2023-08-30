const web3 = require("web3")
const init_holders = [
  {
     address: "0x0958dE5eE71B51f883cFCb8860866f7298c336B2",
     balance: web3.utils.toBN("500000000000000000000000000").toString("hex") // 500000000e18
  }
];


exports = module.exports = init_holders
