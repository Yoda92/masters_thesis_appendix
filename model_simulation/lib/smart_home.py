from lib import battery
from lib import load_profile
from lib import generation_profile
import uuid

BATTERY_CONFIG_KEY = 'battery'
LOAD_PROFILE_CONFIG_KEY = 'demand'
GENERATION_PROFILE_CONFIG_KEY = 'generation'

class SmartHome:
    def __init__(self, config: dict):
        self.id = uuid.uuid4()
        self.load_profile = load_profile.LoadProfile(config.get(LOAD_PROFILE_CONFIG_KEY, {}))
        self.generation_profile = generation_profile.GenerationProfile(config.get(GENERATION_PROFILE_CONFIG_KEY, {}))
        self.battery = battery.Battery(config.get(BATTERY_CONFIG_KEY, {}))

