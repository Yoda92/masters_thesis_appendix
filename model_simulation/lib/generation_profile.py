import pandas as pd
import os

DEFAULT_GENERATION_PROFILE_FILENAME = 'generation_profile.csv'
INCLUDE_PV_GENERATION_PROFILE_KEY = "include_pv_generation_profile"

class GenerationProfile:
    def get_default_time_slots(self):
        return pd.read_csv(os.path.dirname(__file__) + '/../../data_analysis/ausgrid/out/' + DEFAULT_GENERATION_PROFILE_FILENAME)

    def __init__(self, config: dict):   
        time_slots = self.get_default_time_slots()     
        if (len(config) == 0 or config[INCLUDE_PV_GENERATION_PROFILE_KEY] == True):
            self.time_slots = time_slots
        else:
            time_slots.loc[:,'mean'] = 0
            time_slots.loc[:,'percentage'] = 0
            self.time_slots = time_slots