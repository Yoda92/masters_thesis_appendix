import csv
import json
import math
from matplotlib import pyplot as plt
import pyomo.environ as pyo
import lib.timeslots
import lib.electricity_prices
import lib.smart_home
import lib.pyomo_compatible_array
import pandas as pd
import os

USE_AVERAGE_PRICE = True
NUMBER_OF_HOUSEHOLDS = 33
INCLUDE_MAX_GRID_CONSUMPTION_CONSTRAINT = False
MAX_GRID_CONSUMPTION = 30
INCLUDE_MAX_BATTERY_CYCLES_CONSTRAINT = False
MAX_BATTERY_CYCLES = 1
HOUSEHOLDS_WITH_PV_RATIO = 0
MIN_BATTERY_CAPACITY = 0
MAX_BATTERY_CAPACITY = 81
CHARGING_POWER = 30
DISCHARGING_POWER = 30
INITIAL_BATTERY_SOC = 0
DEMAND_PER_HOUSEHOLD_PER_DAY = round((4000 / 365), 2)
GENERATION_PER_HOUSEHOLD_PER_DAY = round((4400 / 365), 2)
PEAK_PERIODS = [5, 6, 7, 8, 9, 17, 18, 19, 20, 21]


def battery_soc_constraint(model, i):
    if i == model.period.first():
        return model.battery_soc[i] == INITIAL_BATTERY_SOC
    else:
        return model.battery_soc[i] == (
            model.battery_soc[i - 1]
            + (model.grid_charge_power[i - 1])
            + (model.pv_charge_power[i - 1])
            - (model.discharge_power[i - 1])
        )


def max_battery_soc_charge_constraint(model, i):
    return model.grid_charge_power[i] + model.pv_charge_power[i] <= (
        MAX_BATTERY_CAPACITY - model.battery_soc[i]
    )


def min_battery_soc_discharge_constraint(model, i):
    return model.discharge_power[i] <= model.battery_soc[i]


def consumption_constraint(model, i):
    return (
        model.discharge_power[i] + model.pv_energy_consumption[i]
        <= model.energy_consumption[i]
    )


def charge_power_state_contraint(model, i):
    return (
        model.grid_charge_power[i] + model.pv_charge_power[i]
        <= model.charge_power_state[i] * CHARGING_POWER
    )


def discharge_power_state_contraint(model, i):
    return (
        model.discharge_power[i] <= model.discharge_power_state[i] * DISCHARGING_POWER
    )


def battery_state_constraint(model, i):
    return model.charge_power_state[i] + model.discharge_power_state[i] <= 1


def max_charge_power_constraint(model, i):
    return model.grid_charge_power[i] + model.pv_charge_power[i] <= CHARGING_POWER


def generation_constraint(model, i):
    return (
        model.pv_charge_power[i] + model.pv_energy_consumption[i]
        <= model.energy_generation[i]
    )


def max_grid_peak_constraint(model, i):
    return get_energy_consumption_from_grid_in_period(model, i) <= MAX_GRID_CONSUMPTION


def max_battery_cycles_constraint(model):
    return get_battery_cycles(model) <= MAX_BATTERY_CYCLES


def get_battery_cycles(model: pyo.ConcreteModel):
    return (
        sum(
            [
                model.discharge_power[period]
                + model.grid_charge_power[period]
                + model.pv_charge_power[period]
                for period in model.period
            ]
        )
        / MAX_BATTERY_CAPACITY
    )


def get_energy_consumption_from_grid_in_period(model: pyo.ConcreteModel, period):
    return (
        model.energy_consumption[period]
        + model.grid_charge_power[period]
        - model.discharge_power[period]
        - model.pv_energy_consumption[period]
    )


def get_energy_consumption_from_generation_in_period(model: pyo.ConcreteModel, period):
    return model.pv_charge_power[period] + model.pv_energy_consumption[period]


def get_alternative_cost(model: pyo.ConcreteModel):
    return sum(
        [
            (model.energy_consumption[period] * model.grid_prices[period])
            for period in model.period
        ]
    )


def cost_objective(model: pyo.ConcreteModel):
    total_cost = sum(
        (
            get_energy_consumption_from_grid_in_period(model, period)
            * model.grid_prices[period]
            + get_energy_consumption_from_generation_in_period(model, period)
            * model.pv_prices[period]
        )
        for period in model.period
    )

    return total_cost


def get_config():
    config = {"households": []}
    number_of_households_with_pv = math.ceil(
        NUMBER_OF_HOUSEHOLDS * HOUSEHOLDS_WITH_PV_RATIO
    )
    number_of_households_without_pv = (
        NUMBER_OF_HOUSEHOLDS - number_of_households_with_pv
    )
    for _ in range(number_of_households_with_pv):
        config["households"].append(
            {"generation": {"include_pv_generation_profile": True}}
        )
    for _ in range(number_of_households_without_pv):
        config["households"].append(
            {"generation": {"include_pv_generation_profile": False}}
        )

    return config


def run_simulation():
    config = get_config()

    households = [
        lib.smart_home.SmartHome(household_config)
        for household_config in config["households"]
    ]
    timeslot_count = 24
    energy_consumption_from_grid_aggregated_timelots = [
        0 for _ in range(timeslot_count)
    ]
    energy_generation_to_grid_aggregated_timelots = [0 for _ in range(timeslot_count)]
    for household in households:
        df_demand = household.load_profile.time_slots
        df_generation = household.generation_profile.time_slots
        consumption = (df_demand["percentage"] / 100) * DEMAND_PER_HOUSEHOLD_PER_DAY
        generation = (
            df_generation["percentage"] / 100
        ) * GENERATION_PER_HOUSEHOLD_PER_DAY
        generation_consumption = consumption.copy()
        for index, _ in generation_consumption.items():
            generation_consumption.iloc[index] = min(
                consumption.iloc[index], generation.iloc[index]
            )
        generation_coverage = generation_consumption.sum() / consumption.sum()
        household.generation_coverage = generation_coverage

        consumption_from_grid = consumption - generation_consumption
        household.consumption_from_grid = (
            lib.pyomo_compatible_array.PyomoCompatibleArray(consumption_from_grid)
        )
        generation_to_grid = generation - generation_consumption
        household.generation_to_grid = lib.pyomo_compatible_array.PyomoCompatibleArray(
            generation_to_grid
        )
        household.generation_to_self = generation_consumption

        for i in range(timeslot_count):
            energy_consumption_from_grid_aggregated_timelots[i] = (
                energy_consumption_from_grid_aggregated_timelots[i]
                + consumption_from_grid[i]
            )
            energy_generation_to_grid_aggregated_timelots[i] = (
                energy_generation_to_grid_aggregated_timelots[i] + generation_to_grid[i]
            )

    model = pyo.ConcreteModel()
    opt = pyo.SolverFactory("glpk")

    grid_prices = lib.electricity_prices.get_default_prices(USE_AVERAGE_PRICE)["mean"]
    grid_prices = grid_prices.copy()
    pv_prices = grid_prices.copy() * 0

    model.period = pyo.Set(
        initialize=list(range(1, 25)),
        ordered=True,
    )

    model.households = lib.pyomo_compatible_array.PyomoCompatibleArray(households)
    model.grid_prices = lib.pyomo_compatible_array.PyomoCompatibleArray(grid_prices)
    model.pv_prices = lib.pyomo_compatible_array.PyomoCompatibleArray(pv_prices)
    model.energy_consumption = lib.pyomo_compatible_array.PyomoCompatibleArray(
        energy_consumption_from_grid_aggregated_timelots
    )
    model.energy_generation = lib.pyomo_compatible_array.PyomoCompatibleArray(
        energy_generation_to_grid_aggregated_timelots
    )

    model.battery_soc = pyo.Var(
        model.period, bounds=(MIN_BATTERY_CAPACITY, MAX_BATTERY_CAPACITY)
    )
    model.grid_charge_power = pyo.Var(model.period, bounds=(0, CHARGING_POWER))
    model.pv_charge_power = pyo.Var(model.period, bounds=(0, CHARGING_POWER))
    model.charge_power_state = pyo.Var(model.period, domain=pyo.Binary)
    model.discharge_power = pyo.Var(model.period, bounds=(0, DISCHARGING_POWER))
    model.discharge_power_state = pyo.Var(model.period, domain=pyo.Binary)
    model.pv_energy_consumption = pyo.Var(model.period, domain=pyo.NonNegativeReals)

    model.obj = pyo.Objective(rule=lambda model: cost_objective(model), sense=pyo.minimize)

    model.battery_soc_constraint = pyo.Constraint(
        model.period, rule=battery_soc_constraint
    )
    model.max_battery_soc_charge_constraint = pyo.Constraint(
        model.period, rule=max_battery_soc_charge_constraint
    )
    model.min_battery_soc_discharge_constraint = pyo.Constraint(
        model.period, rule=min_battery_soc_discharge_constraint
    )
    model.consumption_constraint = pyo.Constraint(
        model.period, rule=consumption_constraint
    )
    model.max_charge_power_constraint = pyo.Constraint(
        model.period, rule=max_charge_power_constraint
    )
    model.generation_constraint = pyo.Constraint(
        model.period, rule=generation_constraint
    )
    model.charge_power_state_contraint = pyo.Constraint(
        model.period, rule=charge_power_state_contraint
    )
    model.discharge_power_state_contraint = pyo.Constraint(
        model.period, rule=discharge_power_state_contraint
    )
    model.battery_state_constraint = pyo.Constraint(
        model.period, rule=battery_state_constraint
    )

    if INCLUDE_MAX_GRID_CONSUMPTION_CONSTRAINT:
        model.peak_constraint = pyo.Constraint(
            model.period, rule=max_grid_peak_constraint
        )

    if INCLUDE_MAX_BATTERY_CYCLES_CONSTRAINT:
        model.battery_cycle_constraint = pyo.Constraint(
            rule=max_battery_cycles_constraint
        )

    results = opt.solve(model)

    if results.solver.termination_condition == pyo.TerminationCondition.optimal:
        print("Optimal solution found.")
        print(f"Minimum cost: {pyo.value(model.obj):.2f}")
        print(f"Alternative cost: {get_alternative_cost(model):.2f}")
    else:
        print("No optimal solution found.")
        exit()

    grid_charging_schedule = [
        pyo.value(model.grid_charge_power[period]) for period in model.period
    ]
    generation_charging_schedule = [
        pyo.value(model.pv_charge_power[period]) for period in model.period
    ]
    discharging_schedule = [pyo.value(model.discharge_power[period]) for period in model.period]
    consumption_from_grid = [
        pyo.value(get_energy_consumption_from_grid_in_period(model, period))
        for period in model.period
    ]
    consumption_from_generation = [
        pyo.value(get_energy_consumption_from_generation_in_period(model, period))
        for period in model.period
    ]

    energy_consumption = [pyo.value(model.energy_consumption[period]) for period in model.period]
    energy_generation = [pyo.value(model.energy_generation[period]) for period in model.period]
    battery_soc = [pyo.value(model.battery_soc[period]) for period in model.period]

    df = pd.DataFrame()

    df["timeslot"] = lib.timeslots.DEFAULT_TIMESLOTS
    df["total_demand"] = energy_consumption
    df["total_generation"] = energy_generation
    df["consumption_from_grid"] = consumption_from_grid
    df["consumption_from_generation"] = consumption_from_generation
    df["grid_charge_power"] = grid_charging_schedule
    df["generation_charge_power"] = generation_charging_schedule
    df["discharge_power"] = discharging_schedule
    df["battery_soc"] = battery_soc

    df["battery_soc"].plot()
    plt.show()
    df = df.rename_axis("index")
    df.to_csv(os.path.dirname(__file__) + "/out/simulation_details.csv")

    # Print the charging and discharging schedules
    for i, timeslot in enumerate(lib.timeslots.DEFAULT_TIMESLOTS):
        print(
            f"Timeslot {timeslot}: Consumption = {energy_consumption[i]:.2f}, Charge Power = {grid_charging_schedule[i]:.2f}, PV consumption = {consumption_from_generation[i]:.2f}, PV charge Power = {generation_charging_schedule[i]:.2f}, Discharge Power = {discharging_schedule[i]:.2f}, SoC = {battery_soc[i]:.2f}"
        )

    result_output = {}
    result_output["config"] = {
        "NUMBER_OF_HOUSEHOLDS": NUMBER_OF_HOUSEHOLDS,
        "HOUSEHOLDS_WITH_PV_RATIO": HOUSEHOLDS_WITH_PV_RATIO,
        "MIN_BATTERY_CAPACITY": MIN_BATTERY_CAPACITY,
        "MAX_BATTERY_CAPACITY": MAX_BATTERY_CAPACITY,
        "CHARGING_POWER": CHARGING_POWER,
        "DISCHARGING_POWER": DISCHARGING_POWER,
        "INITIAL_BATTERY_SOC": INITIAL_BATTERY_SOC,
        "DEMAND_PER_HOUSEHOLD_PER_DAY": DEMAND_PER_HOUSEHOLD_PER_DAY,
        "GENERATION_PER_HOUSEHOLD_PER_DAY": GENERATION_PER_HOUSEHOLD_PER_DAY,
    }

    cost = pyo.value(model.obj)
    alternative_cost = get_alternative_cost(model)
    saving = ((alternative_cost - cost) / alternative_cost) * 100
    battery_cycles = pyo.value(get_battery_cycles(model))
    local_generation = sum(energy_generation)
    local_generation_consumption = sum(consumption_from_generation)
    local_generation_usage = (
        (local_generation_consumption / local_generation) * 100
        if local_generation > 0
        else 0
    )

    grid_consumption = sum(consumption_from_grid)

    grid_consumption_in_peak_periods = sum(
        [
            pyo.value(get_energy_consumption_from_grid_in_period(model, period + 1))
            for period in PEAK_PERIODS
        ]
    )

    demand_in_peak_periods = sum(
        [pyo.value(model.energy_consumption[period + 1]) for period in PEAK_PERIODS]
    )

    peak_reduction = (
        (demand_in_peak_periods - grid_consumption_in_peak_periods)
        / demand_in_peak_periods
    ) * 100

    max_peak = max(
        [
            pyo.value(get_energy_consumption_from_grid_in_period(model, period))
            for period in model.period
        ]
    )

    result_output["results"] = {
        "min_cost": round(cost, 2),
        "alternative_cost": round(alternative_cost, 2),
        "saving": round(saving, 2),
        "battery_cycles": round(battery_cycles, 2),
        "local_generation": round(local_generation, 2),
        "local_generation_usage": round(local_generation_usage, 2),
        "grid_consumption": round(grid_consumption, 2),
        "grid_consumption_in_peak_periods": round(grid_consumption_in_peak_periods, 2),
        "demand_in_peak_periods": round(demand_in_peak_periods, 2),
        "peak_reduction": round(peak_reduction, 2),
        "max_peak": round(max_peak, 2),
    }

    with open(
        os.path.dirname(__file__) + "/out/simulation_overview.json", "w", newline=""
    ) as outfile:
        json.dump(result_output, outfile)


def main():
    run_simulation()


if __name__ == "__main__":
    main()
