// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

contract ElectricityProvider {
    // This is a stub. It is assumed that this function would use an oracle in a real-world scenario.
    function getPrice(uint _period) public pure returns (uint) {        
        return _period % 100;
    }
}
