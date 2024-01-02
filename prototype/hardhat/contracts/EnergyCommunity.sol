// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "./SmartMeter.sol";
import "./BESS.sol";
import "./ElectricityProvider.sol";

contract EnergyCommunity {
    address private dso_address;
    address private aggregator_address;
    address private electricity_provider_address;
    address[] private smart_meter_addresses;
    BESS private bess;
    ElectricityProvider private electricity_provider;

    struct EnergySettlement {
        uint amount;
        address receiver;
    }

    event SettlementCompleted(uint amount, address payer, address receiver);

    mapping(address => bool) private smart_meter_registered;
    mapping(address => bool) private smart_meter_approved;
    mapping(address => SmartMeter) private smart_meters;
    mapping(uint => bool) private settlement_period_closed;
    mapping(uint => mapping(address => EnergySettlement)) private settlements;

    constructor(
        address _dso_address,
        address _aggregator_address,
        address _electricity_provider_address
    ) {
        dso_address = _dso_address;
        aggregator_address = _aggregator_address;
        electricity_provider_address = _electricity_provider_address;
        bess = new BESS();
        electricity_provider = new ElectricityProvider();
    }

    modifier onlyAggregator() {
        require(
            msg.sender == aggregator_address,
            "Only the aggregator can perform this action."
        );
        _;
    }

    modifier onlySmartMeter() {
        require(
            smart_meter_registered[msg.sender],
            "Only a smart meter can perform this action."
        );
        _;
    }

    modifier canPerformSettlement(uint _period) {
        require(isDataPublished(_period), "Not all data has been published.");
        _;
    }

    function isDataPublished(uint _period) public view returns (bool) {
        for (uint index = 0; index < smart_meter_addresses.length; index++) {
            address smart_meter_address = smart_meter_addresses[index];
            if (
                !smart_meters[smart_meter_address]
                    .getPeriodicData(_period)
                    .is_published
            ) {
                return false;
            }
        }

        if (!bess.getPeriodicData(_period).is_published) {
            return false;
        }

        return true;
    }

    function registerSmartMeter() public {
        require(
            smart_meter_approved[msg.sender],
            "SmartMeter must be approved by aggregator."
        );

        require(
            !smart_meter_registered[msg.sender],
            "SmartMeter is already registered."
        );

        smart_meter_addresses.push(msg.sender);
        smart_meter_registered[msg.sender] = true;
        smart_meters[msg.sender] = new SmartMeter(msg.sender);
    }

    function approveSmartMeter(
        address _smart_meter_address
    ) public onlyAggregator {
        smart_meter_approved[_smart_meter_address] = true;
    }

    function publishSmartMeterPeriodicData(
        uint _consumption_from_grid,
        uint _production_to_grid,
        uint _settlement_period
    ) public onlySmartMeter {
        require(
            !settlement_period_closed[_settlement_period],
            "Periodic data cannot be published after settlement has been closed."
        );

        smart_meters[msg.sender].publishPeriodicData(
            _consumption_from_grid,
            _production_to_grid,
            _settlement_period
        );
    }

    function publishBESSPeriodicData(
        uint _consumption_from_grid,
        uint _production_to_grid,
        uint _soc,
        uint _settlement_period
    ) public onlyAggregator {
        require(
            !settlement_period_closed[_settlement_period],
            "Periodic data cannot be published after settlement has been closed."
        );

        bess.publishPeriodicData(
            _consumption_from_grid,
            _production_to_grid,
            _soc,
            _settlement_period
        );
    }

    function performSettlement(
        uint _settlement_period
    ) public canPerformSettlement(_settlement_period) {
        uint settlement_period_price = electricity_provider.getPrice(
            _settlement_period
        );
        // Aggregator settlement to smart meters.
        for (uint index = 0; index < smart_meter_addresses.length; index++) {
            address smart_meter_address = smart_meter_addresses[index];

            uint smart_meter_consumption = smart_meters[smart_meter_address]
                .getPeriodicData(_settlement_period)
                .consumption_from_grid;

            uint consumption_price = smart_meter_consumption *
                settlement_period_price;

            settlements[_settlement_period][smart_meter_address]
                .amount = consumption_price;
            settlements[_settlement_period][smart_meter_address]
                .receiver = aggregator_address;

            emit SettlementCompleted(
                settlements[_settlement_period][smart_meter_address].amount,
                smart_meter_address,
                settlements[_settlement_period][smart_meter_address].receiver
            );
        }

        // Aggregator settlement to energy provider.
        uint consumption_from_grid_in_period = getCommunityNetEnergyConsumptionInPeriod(
                _settlement_period
            );

        uint net_community_consumption_price = consumption_from_grid_in_period *
            settlement_period_price;

        settlements[_settlement_period][aggregator_address]
            .amount = net_community_consumption_price;
        settlements[_settlement_period][aggregator_address]
            .receiver = electricity_provider_address;

        emit SettlementCompleted(
            settlements[_settlement_period][aggregator_address].amount,
            aggregator_address,
            settlements[_settlement_period][aggregator_address].receiver
        );

        settlement_period_closed[_settlement_period] = true;
    }

    function getCommunityNetEnergyConsumptionInPeriod(
        uint _period
    ) private view returns (uint) {
        uint totalCommunityConsumptionInPeriod = getTotalCommunityConsumptionInPeriod(
                _period
            );
        uint totalCommunityProductionInPeriod = getTotalCommunityProductionInPeriod(
                _period
            );

        uint netCommunityEnergyUsageInPeriod = totalCommunityConsumptionInPeriod -
                min(
                    totalCommunityConsumptionInPeriod,
                    totalCommunityProductionInPeriod
                );

        return netCommunityEnergyUsageInPeriod;
    }

    function min(uint a, uint b) private pure returns (uint) {
        return a < b ? a : b;
    }

    function getTotalCommunityConsumptionInPeriod(
        uint _period
    ) private view returns (uint) {
        uint result = 0;

        for (uint index = 0; index < smart_meter_addresses.length; index++) {
            result += smart_meters[smart_meter_addresses[index]]
                .getPeriodicData(_period)
                .consumption_from_grid;
        }

        BESS.PeriodicData memory bessPeriodicData = bess.getPeriodicData(
            _period
        );
        result += bessPeriodicData.consumption_from_grid;

        return result;
    }

    function getTotalCommunityProductionInPeriod(
        uint _period
    ) private view returns (uint) {
        uint result = 0;

        for (uint index = 0; index < smart_meter_addresses.length; index++) {
            result += smart_meters[smart_meter_addresses[index]]
                .getPeriodicData(_period)
                .production_to_grid;
        }

        BESS.PeriodicData memory bessPeriodicData = bess.getPeriodicData(
            _period
        );
        result += bessPeriodicData.production_to_grid;

        return result;
    }
}
