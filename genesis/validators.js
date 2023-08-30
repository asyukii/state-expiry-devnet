const web3 = require("web3")
const RLP = require('rlp');

// Configure
const validators = [
  
   {
     "consensusAddr": "0x0958dE5eE71B51f883cFCb8860866f7298c336B2",
   },
   {
     "consensusAddr": "0xDa7240C984ECB157999eEFef066b317dFD5bafFB",
   },
   {
     "consensusAddr": "0x9F8FAAebf4B467488c9Ab2c703e5FA72f1037DFd",
   },
];

// ===============  Do not edit below ====
function generateExtradata(validators) {
  let extraVanity =Buffer.alloc(32);
  let validatorsBytes = extraDataSerialize(validators);
  let extraSeal =Buffer.alloc(65);
  return Buffer.concat([extraVanity,validatorsBytes,extraSeal]);
}

function extraDataSerialize(validators) {
  let n = validators.length;
  let arr = [];
  for (let i = 0;i<n;i++) {
    let validator = validators[i];
    arr.push(Buffer.from(web3.utils.hexToBytes(validator.consensusAddr)));
  }
  return Buffer.concat(arr);
}

extraValidatorBytes = generateExtradata(validators);

exports = module.exports = {
  extraValidatorBytes: extraValidatorBytes,
}