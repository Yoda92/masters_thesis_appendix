import os
import time
import requests
import pandas as pd
from datetime import datetime, timedelta

BASE_DATE_ARG = "2023-10"
URL = "https://nrgi.dk/api/common/pricehistory?region=DK1&date="
INCLUDE_GRID_ARG = "&includeGrid=false"
DEFAULT_TIMESLOTS = [
    "1:00",
    "2:00",
    "3:00",
    "4:00",
    "5:00",
    "6:00",
    "7:00",
    "8:00",
    "9:00",
    "10:00",
    "11:00",
    "12:00",
    "13:00",
    "14:00",
    "15:00",
    "16:00",
    "17:00",
    "18:00",
    "19:00",
    "20:00",
    "21:00",
    "22:00",
    "23:00",
    "0:00",
]
HEADERS = ["date"] + DEFAULT_TIMESLOTS

df = pd.DataFrame(columns=HEADERS)

for day in range(1, 31 + 1):
    date = BASE_DATE_ARG + "-" + str(day)
    response = requests.get(URL + date + INCLUDE_GRID_ARG)
    data = response.json()["prices"]

    timestamps = [entry["localTime"] for entry in data]
    prices_incl_vat = [entry["kwPrice"] for entry in data]
    adjusted_timestamps = [
        (
            datetime.strptime(timestamp, "%Y-%m-%dT%H:%M:%S") + timedelta(hours=1)
        ).strftime("%#H:%M")
        for timestamp in timestamps
    ]

    newRow = pd.DataFrame(
        data=[date] + prices_incl_vat, index=["date"] + adjusted_timestamps
    )
    df = pd.concat([df, newRow.T])
    time.sleep(1)

df = df.rename_axis('index')
df.to_csv(
    os.path.dirname(__file__) + "/out/electricity_price_original_data_month.csv"
)
