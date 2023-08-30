const fs = require("fs");
const readline = require('readline');
const nunjucks = require("nunjucks");

async function processValidatorConf() {
  const fileStream = fs.createReadStream(__dirname + "/validators.conf");

  const rl = readline.createInterface({
    input: fileStream,
    crlfDelay: Infinity
  });
  let validators = [];
  for await (const line of rl) {
    validators.push({
        consensusAddr: line
    });
  }
  return validators
}

processValidatorConf().then(function (validators) {
  const data = {
    validators: validators,
  };
  const templateString = fs.readFileSync(__dirname + '/validators.template').toString();
  const resultString = nunjucks.renderString(templateString, data);
  fs.writeFileSync(__dirname + '/validators.js', resultString);
  console.log("BSCValidatorSet file updated.");
})
