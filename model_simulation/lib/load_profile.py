import pandas as pd
import os

DEFAULT_LOAD_PROFILE_FILENAME = 'consumption_profile.csv'

class LoadProfile:
    def init_default_time_slots(self):
        self.time_slots = pd.read_csv(os.path.dirname(__file__) + '/../../data_analysis/ausgrid/out/' + DEFAULT_LOAD_PROFILE_FILENAME)

    def __init__(self, config: dict):        
        if (len(config) == 0):
            self.init_default_time_slots()
        else:
            print("Not supported.")
            exit()