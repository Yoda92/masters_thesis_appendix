// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

#![allow(dead_code)]
#![allow(unused_imports)]

use crate::*;
use crate::coreroot::*;

#[derive(Clone)]
pub struct MapStringToImmutableBytes {
    pub(crate) proxy: Proxy,
}

impl MapStringToImmutableBytes {
    pub fn get_bytes(&self, key: &str) -> ScImmutableBytes {
        ScImmutableBytes::new(self.proxy.key(&string_to_bytes(key)))
    }
}

#[derive(Clone)]
pub struct ImmutableDeployContractParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableDeployContractParams {
    pub fn new() -> ImmutableDeployContractParams {
        ImmutableDeployContractParams {
            proxy: params_proxy(),
        }
    }

    // additional params for smart contract init function
    pub fn init_params(&self) -> MapStringToImmutableBytes {
        MapStringToImmutableBytes { proxy: self.proxy.clone() }
    }

    // The name of the contract to be deployed, used to calculate the contract's hname.
    // The hname must be unique among all contract hnames in the chain.
    pub fn name(&self) -> ScImmutableString {
        ScImmutableString::new(self.proxy.root(PARAM_NAME))
    }

    // hash of blob that has been previously stored in blob contract
    pub fn program_hash(&self) -> ScImmutableHash {
        ScImmutableHash::new(self.proxy.root(PARAM_PROGRAM_HASH))
    }
}

#[derive(Clone)]
pub struct MapStringToMutableBytes {
    pub(crate) proxy: Proxy,
}

impl MapStringToMutableBytes {
    pub fn clear(&self) {
        self.proxy.clear_map();
    }

    pub fn get_bytes(&self, key: &str) -> ScMutableBytes {
        ScMutableBytes::new(self.proxy.key(&string_to_bytes(key)))
    }
}

#[derive(Clone)]
pub struct MutableDeployContractParams {
    pub(crate) proxy: Proxy,
}

impl MutableDeployContractParams {
    // additional params for smart contract init function
    pub fn init_params(&self) -> MapStringToMutableBytes {
        MapStringToMutableBytes { proxy: self.proxy.clone() }
    }

    // The name of the contract to be deployed, used to calculate the contract's hname.
    // The hname must be unique among all contract hnames in the chain.
    pub fn name(&self) -> ScMutableString {
        ScMutableString::new(self.proxy.root(PARAM_NAME))
    }

    // hash of blob that has been previously stored in blob contract
    pub fn program_hash(&self) -> ScMutableHash {
        ScMutableHash::new(self.proxy.root(PARAM_PROGRAM_HASH))
    }
}

#[derive(Clone)]
pub struct ImmutableGrantDeployPermissionParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableGrantDeployPermissionParams {
    pub fn new() -> ImmutableGrantDeployPermissionParams {
        ImmutableGrantDeployPermissionParams {
            proxy: params_proxy(),
        }
    }

    // agent to grant deploy permission to
    pub fn deployer(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_DEPLOYER))
    }
}

#[derive(Clone)]
pub struct MutableGrantDeployPermissionParams {
    pub(crate) proxy: Proxy,
}

impl MutableGrantDeployPermissionParams {
    // agent to grant deploy permission to
    pub fn deployer(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_DEPLOYER))
    }
}

#[derive(Clone)]
pub struct ImmutableRequireDeployPermissionsParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableRequireDeployPermissionsParams {
    pub fn new() -> ImmutableRequireDeployPermissionsParams {
        ImmutableRequireDeployPermissionsParams {
            proxy: params_proxy(),
        }
    }

    // turns permission check on or off
    pub fn deploy_permissions_enabled(&self) -> ScImmutableBool {
        ScImmutableBool::new(self.proxy.root(PARAM_DEPLOY_PERMISSIONS_ENABLED))
    }
}

#[derive(Clone)]
pub struct MutableRequireDeployPermissionsParams {
    pub(crate) proxy: Proxy,
}

impl MutableRequireDeployPermissionsParams {
    // turns permission check on or off
    pub fn deploy_permissions_enabled(&self) -> ScMutableBool {
        ScMutableBool::new(self.proxy.root(PARAM_DEPLOY_PERMISSIONS_ENABLED))
    }
}

#[derive(Clone)]
pub struct ImmutableRevokeDeployPermissionParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableRevokeDeployPermissionParams {
    pub fn new() -> ImmutableRevokeDeployPermissionParams {
        ImmutableRevokeDeployPermissionParams {
            proxy: params_proxy(),
        }
    }

    // agent to revoke deploy permission for
    pub fn deployer(&self) -> ScImmutableAgentID {
        ScImmutableAgentID::new(self.proxy.root(PARAM_DEPLOYER))
    }
}

#[derive(Clone)]
pub struct MutableRevokeDeployPermissionParams {
    pub(crate) proxy: Proxy,
}

impl MutableRevokeDeployPermissionParams {
    // agent to revoke deploy permission for
    pub fn deployer(&self) -> ScMutableAgentID {
        ScMutableAgentID::new(self.proxy.root(PARAM_DEPLOYER))
    }
}

#[derive(Clone)]
pub struct ImmutableFindContractParams {
    pub(crate) proxy: Proxy,
}

impl ImmutableFindContractParams {
    pub fn new() -> ImmutableFindContractParams {
        ImmutableFindContractParams {
            proxy: params_proxy(),
        }
    }

    // The smart contract’s Hname
    pub fn hname(&self) -> ScImmutableHname {
        ScImmutableHname::new(self.proxy.root(PARAM_HNAME))
    }
}

#[derive(Clone)]
pub struct MutableFindContractParams {
    pub(crate) proxy: Proxy,
}

impl MutableFindContractParams {
    // The smart contract’s Hname
    pub fn hname(&self) -> ScMutableHname {
        ScMutableHname::new(self.proxy.root(PARAM_HNAME))
    }
}
