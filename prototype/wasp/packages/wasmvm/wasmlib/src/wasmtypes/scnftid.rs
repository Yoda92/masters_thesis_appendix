// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use std::cmp::Ordering;

use crate::*;

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

pub const SC_NFT_ID_LENGTH: usize = 32;

#[derive(PartialEq, Clone, Copy, Eq, Hash, Ord)]
pub struct ScNftID {
    id: [u8; SC_NFT_ID_LENGTH],
}

impl ScNftID {
    pub fn to_bytes(&self) -> Vec<u8> {
        nft_id_to_bytes(self)
    }

    pub fn to_string(&self) -> String {
        nft_id_to_string(self)
    }
}

impl PartialOrd for ScNftID {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.id.cmp(&other.id))
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

pub fn nft_id_decode(dec: &mut WasmDecoder) -> ScNftID {
    let buf = dec.fixed_bytes(SC_NFT_ID_LENGTH);
    ScNftID {
        id: buf.try_into().expect("WTF?"),
    }
}

pub fn nft_id_encode(enc: &mut WasmEncoder, value: &ScNftID) {
    enc.fixed_bytes(&value.id, SC_NFT_ID_LENGTH);
}

pub fn nft_id_from_bytes(buf: &[u8]) -> ScNftID {
    let len = buf.len();
    if len == 0 {
        return ScNftID {
            id: [0; SC_NFT_ID_LENGTH],
        };
    }
    if len != SC_NFT_ID_LENGTH {
        panic("invalid NftID length");
    }
    ScNftID {
        id: buf.try_into().expect("WTF?"),
    }
}

pub fn nft_id_to_bytes(value: &ScNftID) -> Vec<u8> {
    value.id.to_vec()
}

pub fn nft_id_from_string(value: &str) -> ScNftID {
    nft_id_from_bytes(&hex_decode(&value))
}

pub fn nft_id_to_string(value: &ScNftID) -> String {
    hex_encode(&value.id)
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

pub struct ScImmutableNftID {
    proxy: Proxy,
}

impl ScImmutableNftID {
    pub fn new(proxy: Proxy) -> ScImmutableNftID {
        ScImmutableNftID { proxy }
    }

    pub fn exists(&self) -> bool {
        self.proxy.exists()
    }

    pub fn to_string(&self) -> String {
        nft_id_to_string(&self.value())
    }

    pub fn value(&self) -> ScNftID {
        nft_id_from_bytes(&self.proxy.get())
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for mutable ScNftID in host container
pub struct ScMutableNftID {
    proxy: Proxy,
}

impl ScMutableNftID {
    pub fn new(proxy: Proxy) -> ScMutableNftID {
        ScMutableNftID { proxy }
    }

    pub fn delete(&self) {
        self.proxy.delete();
    }

    pub fn exists(&self) -> bool {
        self.proxy.exists()
    }

    pub fn set_value(&self, value: &ScNftID) {
        self.proxy.set(&nft_id_to_bytes(&value));
    }

    pub fn to_string(&self) -> String {
        nft_id_to_string(&self.value())
    }

    pub fn value(&self) -> ScNftID {
        nft_id_from_bytes(&self.proxy.get())
    }
}
