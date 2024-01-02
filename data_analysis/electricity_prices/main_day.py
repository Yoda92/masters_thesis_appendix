import matplotlib.pyplot as plt
import pandas as pd
import os

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
    df = pd.read_csv(
        os.path.dirname(__file__) + "/out/electricity_price_original_data_day.csv"
    )
    df_only_prices = df.loc[:, "1:00":"0:00"]
    df_only_prices = (df_only_prices / 100)
    df_only_prices = df_only_prices.rename(columns=COLUMN_MAP)
    description = df_only_prices.describe()
    output = description.T.rename_axis('time')
    output = output.rename(columns={"25%": "25_percent", "50%": "50_percent", "75%": "75_percent"})
    output.plot()
    plt.show()
    output.to_csv(os.path.dirname(__file__) + "/out/electricity_prices_day.csv")


if __name__ == "__main__":
    main()
