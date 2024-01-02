from web3 import Web3
import subprocess
import pandas as pd
import os

NUMBER_OF_ACCOUNTS = 36
WASP_CLI_DIR = r"C:\Users\ander\Documents\Github\masters_thesis\prototype\wasp-cli\wasp-cli_1.0.1-rc.10_Windows_x86_64"

w3 = Web3()

accounts = []

for index in range(NUMBER_OF_ACCOUNTS):
    account = w3.eth.account.create()
    command = ["wasp-cli", "chain", "deposit", account.address, "base:1000000"]

    result = subprocess.run(
        command,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        cwd=WASP_CLI_DIR,
    )

    if result.returncode == 0:
        print("Created new account.")
        accounts.append(account)
    else:
        print("Command failed.")
        print("Error output:")
        print(result.stderr)
        exit()

df = pd.DataFrame()
df["address"] = [account.address for account in accounts]
df["private_key"] = [Web3.to_hex(account.key) for account in accounts]

df = df.rename_axis("index")
df.to_json(os.path.dirname(__file__) + "/out/accounts.json", orient="records")
