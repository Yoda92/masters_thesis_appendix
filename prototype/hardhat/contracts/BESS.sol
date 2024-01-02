// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

contract BESS {
    mapping(uint => PeriodicData) periodic_data;

    struct PeriodicData {
        uint consumption_from_grid;
        uint production_to_grid;
        uint soc;
        bool is_published;
    }

    constructor() {}

    function getPeriodicData(
        uint _period
    ) public view returns (PeriodicData memory) {
        return periodic_data[_period];
    }

    function publishPeriodicData(
        uint _consumption_from_grid,
        uint _production_to_grid,
        uint _soc,
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
        require(
            _consumption_from_grid == 0 || _production_to_grid == 0,
            "BESS cannot be charging and discharging at the same time."
        );

        periodic_data[_period].soc = _soc;
        periodic_data[_period].consumption_from_grid = _consumption_from_grid;
        periodic_data[_period].production_to_grid = _production_to_grid;
        periodic_data[_period].is_published = true;
    }
}
