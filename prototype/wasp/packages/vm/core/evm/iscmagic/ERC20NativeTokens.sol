// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

pragma solidity >=0.8.11;

import "./ISCTypes.sol";
import "./ISCSandbox.sol";
import "./ISCAccounts.sol";
import "./ISCPrivileged.sol";

// The ERC20 contract for ISC L2 native tokens (on-chain foundry).
contract ERC20NativeTokens {
    string _name;
    string _tickerSymbol;
    uint8 _decimals;

    event Approval(
        address indexed tokenOwner,
        address indexed spender,
        uint256 tokens
    );
    event Transfer(address indexed from, address indexed to, uint256 tokens);

    function foundrySerialNumber() internal view returns (uint32) {
        return __iscSandbox.erc20NativeTokensFoundrySerialNumber(address(this));
    }

    function nativeTokenID()
        public
        view
        virtual
        returns (NativeTokenID memory)
    {
        return __iscSandbox.getNativeTokenID(foundrySerialNumber());
    }

    function name() public view returns (string memory) {
        return _name;
    }

    function symbol() public view returns (string memory) {
        return _tickerSymbol;
    }

    function decimals() public view returns (uint8) {
        return _decimals;
    }

    function totalSupply() public view virtual returns (uint256) {
        return
            __iscSandbox
                .getNativeTokenScheme(foundrySerialNumber())
                .maximumSupply;
    }

    function balanceOf(address tokenOwner) public view returns (uint256) {
        ISCChainID chainID = __iscSandbox.getChainID();
        ISCAgentID memory ownerAgentID = ISCTypes.newEthereumAgentID(
            tokenOwner,
            chainID
        );
        return
            __iscAccounts.getL2BalanceNativeTokens(
                nativeTokenID(),
                ownerAgentID
            );
    }

    function transfer(
        address receiver,
        uint256 numTokens
    ) public returns (bool) {
        ISCAssets memory assets;
        assets.nativeTokens = new NativeToken[](1);
        assets.nativeTokens[0].ID = nativeTokenID();
        assets.nativeTokens[0].amount = numTokens;
        __iscPrivileged.moveBetweenAccounts(msg.sender, receiver, assets);
        emit Transfer(msg.sender, receiver, numTokens);
        return true;
    }

    function approve(
        address delegate,
        uint256 numTokens
    ) public returns (bool) {
        __iscPrivileged.setAllowanceNativeTokens(
            msg.sender,
            delegate,
            nativeTokenID(),
            numTokens
        );
        emit Approval(msg.sender, delegate, numTokens);
        return true;
    }

    function allowance(
        address owner,
        address delegate
    ) public view returns (uint256) {
        ISCAssets memory assets = __iscSandbox.getAllowance(owner, delegate);
        NativeTokenID memory myID = nativeTokenID();
        for (uint256 i = 0; i < assets.nativeTokens.length; i++) {
            if (bytesEqual(assets.nativeTokens[i].ID.data, myID.data))
                return assets.nativeTokens[i].amount;
        }
        return 0;
    }

    function bytesEqual(
        bytes memory a,
        bytes memory b
    ) internal pure returns (bool) {
        if (a.length != b.length) return false;
        for (uint256 i = 0; i < a.length; i++) {
            if (a[i] != b[i]) return false;
        }
        return true;
    }

    function transferFrom(
        address owner,
        address buyer,
        uint256 numTokens
    ) public returns (bool) {
        ISCAssets memory assets;
        assets.nativeTokens = new NativeToken[](1);
        assets.nativeTokens[0].ID = nativeTokenID();
        assets.nativeTokens[0].amount = numTokens;
        __iscPrivileged.moveAllowedFunds(owner, msg.sender, assets);
        if (buyer != msg.sender) {
            __iscPrivileged.moveBetweenAccounts(msg.sender, buyer, assets);
        }
        emit Transfer(owner, buyer, numTokens);
        return true;
    }
}
