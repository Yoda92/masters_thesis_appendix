import pandas as pd
import matplotlib.pyplot as plt
import os

FILENAME = "2012-2013 Solar home electricity data v2.csv"

COLUMN_MAP = {
    "1:00": 0,
    "2:00": 1,
    "3:00": 2,
    "4:00": 3,
    "5:00": 4,
    "6:00": 5,
    "7:00": 6,
    "8:00": 7,
    "9:00": 8,
    "10:00": 9,
    "11:00": 10,
    "12:00": 11,
    "13:00": 12,
    "14:00": 13,
    "15:00": 14,
    "16:00": 15,
    "17:00": 16,
    "18:00": 17,
    "19:00": 18,
    "20:00": 19,
    "21:00": 20,
    "22:00": 21,
    "23:00": 22,
    "0:00": 23,
}


def main():
    start_day = pd.to_datetime("01.08.2012", dayfirst=True)
    end_day = pd.to_datetime("31.08.2012", dayfirst=True)

    # Generation data
    generation_profile = get_measuremets(
        start_day, end_day, "GG", [0, 1, 2, 3, 4, 5, 6]
    )
    generation_profile = generation_profile.loc[:, "1:00":"0:00"].describe().T
    generation_profile = generation_profile.rename_axis("time")
    generation_profile = generation_profile.rename(columns=COLUMN_MAP)
    generation_profile["percentage"] = (
        generation_profile["mean"] / generation_profile["mean"].sum()
    ) * 100
    generation_profile["percentage_sum"] = generation_profile["percentage"].cumsum()
    generation_profile = generation_profile.reset_index()
    generation_profile = generation_profile.rename_axis("index")
    generation_profile = generation_profile.rename(columns={"25%": "25_percent", "50%": "50_percent", "75%": "75_percent"})
    generation_profile.to_csv(
        os.path.dirname(__file__) + "/out/generation_profile.csv"
    )

    # Consumption data
    consumption_profile = get_measuremets(
        start_day, end_day, "GC", [0, 1, 2, 3, 4, 5, 6]
    )
    consumption_profile = consumption_profile.loc[:, "1:00":"0:00"].describe().T
    consumption_profile = consumption_profile.rename_axis("time")
    consumption_profile = consumption_profile.rename(columns=COLUMN_MAP)
    consumption_profile["percentage"] = (
        consumption_profile["mean"] / consumption_profile["mean"].sum()
    ) * 100
    consumption_profile["percentage_sum"] = consumption_profile["percentage"].cumsum()
    consumption_profile = consumption_profile.reset_index()
    consumption_profile = consumption_profile.rename_axis("index")
    consumption_profile = consumption_profile.rename(columns={"25%": "25_percent", "50%": "50_percent", "75%": "75_percent"})

    consumption_profile.to_csv(
        os.path.dirname(__file__) + "/out/consumption_profile.csv"
    )


def get_measuremets(start_date, end_date, category, weekdays):
    df = pd.read_csv(os.path.dirname(__file__) + "/" + FILENAME, skiprows=1)

    df["date"] = pd.to_datetime(df["date"], dayfirst=True)

    df_interval = df.loc[
        (df["date"].between(start_date, end_date))
        & (df["date"].dt.dayofweek.isin(weekdays))
    ]

    df_interval_category = df_interval.loc[
        (df_interval["Consumption Category"] == category)
    ]

    measurements_interval = df_interval_category

    for i in range(0, 24):
        first_time = str(i) + ":30"
        last_time = (str(i + 1) if (i + 1) < 24 else str("0")) + ":00"
        measurements_interval[last_time] = (
            measurements_interval[first_time] + measurements_interval[last_time]
        )
        measurements_interval = measurements_interval.drop(columns=[first_time])

    return measurements_interval


if __name__ == "__main__":
    main()
