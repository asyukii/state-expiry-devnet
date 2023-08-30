const { spawn } = require("child_process")
const program = require("commander")
const nunjucks = require("nunjucks")
const fs = require("fs")
const web3 = require("web3")

const validators = require("./validators")
const init_holders = require("./init_holders")

program.option("-c, --chainid <chainid>", "chain id", "123454321")

program.option(
    "-t, --template <template>",
    "Genesis template json",
    "./genesis-template.json"
)

program.option(
    "-o, --output <output-file>",
    "Genesis json file",
    "./genesis.json"
)

program.parse(process.argv)

const data = {
    chainId: program.chainid,
    initHolders: init_holders,
    extraData: web3.utils.bytesToHex(validators.extraValidatorBytes)
}

const templateString = fs.readFileSync(program.template).toString()

const resultString = nunjucks.renderString(templateString, data)

fs.writeFileSync(program.output, resultString)