import pandas as pd
import os
from lib import timeslots

def get_default_prices(use_average_prices = True):
    if use_average_prices:
        prices = pd.read_csv(os.path.dirname(__file__) + '/../../data_analysis/electricity_prices/out/electricity_prices_month.csv')
    else:
        prices = pd.read_csv(os.path.dirname(__file__) + '/../../data_analysis/electricity_prices/out/electricity_prices_day.csv')
    prices["timeslot"] = timeslots.DEFAULT_TIMESLOTS
    return prices