// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Import the necessary interfaces and libraries
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ERC20Token is ERC20, Ownable {
    // Constructor to initialize the token with a name and symbol
    constructor() ERC20("ERC20", "ERC20") {
        // Mint an initial supply of tokens to the contract creator
        _mint(msg.sender, 1000000000 * 10 ** uint256(decimals())); // 1,000,000,000 tokens with decimals = 18
    }

    // Function to mint additional tokens, only callable by the owner
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }

    // Function to burn tokens, only callable by the owner
    function burn(uint256 amount) public onlyOwner {
        _burn(msg.sender, amount);
    }
}