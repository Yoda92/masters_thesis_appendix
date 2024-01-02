import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import jsonData from '../client/out/accounts.json'

const private_keys = jsonData.map(account => account.private_key);

const config: HardhatUserConfig = {
  solidity: "0.8.19",
  defaultNetwork: "local",
  networks: {
    local: {
      url: "http://localhost/wasp/api/v1/chains/tst1pzsgz84lkfk0savr4unumlxayrfyvgyssq6p420ke6gatepplpa0vtngafn/evm",
      chainId: 1074,
      accounts: private_keys,
      timeout: 60000,
    },
  },
};

export default config;
