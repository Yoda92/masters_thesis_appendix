// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

pragma solidity >=0.8.11;

import "./ISCSandbox.sol";
import "./ISCAccounts.sol";
import "./ISCUtil.sol";
import "./ISCPrivileged.sol";
import "./ERC20BaseTokens.sol";
import "./ERC20NativeTokens.sol";
import "./ERC721NFTs.sol";
import "./ERC721NFTCollection.sol";

library ISC {
    ISCSandbox constant sandbox = __iscSandbox;

    ISCAccounts constant accounts = __iscAccounts;

    ISCUtil constant util = __iscUtil;

    ERC20BaseTokens constant baseTokens = __erc20BaseTokens;

    // Get the ERC20NativeTokens contract for the given foundry serial number
    function nativeTokens(uint32 foundrySN) internal view returns (ERC20NativeTokens) {
        return ERC20NativeTokens(sandbox.erc20NativeTokensAddress(foundrySN));
    }

    ERC721NFTs constant nfts = __erc721NFTs;

    // Get the ERC721NFTCollection contract for the given collection
    function erc721NFTCollection(NFTID collectionID) internal view returns (ERC721NFTCollection) {
        return ERC721NFTCollection(sandbox.erc721NFTCollectionAddress(collectionID));
    }

}
