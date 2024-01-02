// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

contract SmartMeter {
    address public owner;
    mapping(uint => PeriodicData) periodicData;

    struct PeriodicData {
        uint consumption_from_grid;
        uint production_to_grid;
        bool is_published;
    }

    constructor(address _owner) {
        owner = _owner;
    }

    function getPeriodicData(
        uint _period
    ) public view returns (PeriodicData memory) {
        return periodicData[_period];
    }

    function publishPeriodicData(
        uint _consumption_from_grid,
        uint _production_to_grid,
        uint _period
    ) public {
        require(
            _consumption_from_grid >= 0,
            "Consumption data should be non-negative."
        );
        require(
            _production_to_grid >= 0,
            "Production data should be non-negative."
        );

        periodicData[_period].consumption_from_grid = _consumption_from_grid;
        periodicData[_period].production_to_grid = _production_to_grid;
        periodicData[_period].is_published = true;
    }
}
