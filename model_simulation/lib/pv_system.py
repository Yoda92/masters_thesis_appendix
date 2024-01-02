import pandas as pd
import os

DEFAULT_GENERATION_PROFILE_FILENAME = 'measurements_winter_gg_normalized_mean.csv'

class PVSystem:
    def init_default_time_slots(self):
        self.time_slots = pd.read_csv(os.path.dirname(__file__) + '/../../data_analysis/ausgrid/out/' + DEFAULT_GENERATION_PROFILE_FILENAME)

    def __init__(self, config: dict):        
        if (len(config) == 0):
            self.init_default_time_slots()
            exit()
        else:
            self.init_default_time_slots()
            exit()