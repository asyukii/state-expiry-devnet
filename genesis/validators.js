const web3 = require("web3")
const RLP = require('rlp');

// Configure
const validators = [
  
   {
     "consensusAddr": "0xC95b7a10ED53C3efe6Ce03e2A35527Bae0c1A3Cf",
   },
   {
     "consensusAddr": "0x2932eCBE405256d9C3913760b02206444CD3A7B0",
   },
   {
     "consensusAddr": "0xFD82F4dF7C5624Cf2F5F4692506013682a16BaBC",
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