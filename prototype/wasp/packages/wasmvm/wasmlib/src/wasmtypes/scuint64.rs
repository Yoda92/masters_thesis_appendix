// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use crate::*;

pub const SC_UINT64_LENGTH: usize = 8;

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

pub fn uint64_decode(dec: &mut WasmDecoder) -> u64 {
    uint64_from_bytes(&dec.fixed_bytes(SC_UINT64_LENGTH))
}

pub fn uint64_encode(enc: &mut WasmEncoder, value: u64) {
    enc.fixed_bytes(&uint64_to_bytes(value), SC_UINT64_LENGTH);
}

pub fn uint64_from_bytes(buf: &[u8]) -> u64 {
    let len = buf.len();
    if len == 0 {
        return 0;
    }
    if len != SC_UINT64_LENGTH {
        panic("invalid Uint64 length");
    }
    u64::from_le_bytes(buf.try_into().expect("WTF?"))
}

pub fn uint64_to_bytes(value: u64) -> Vec<u8> {
    value.to_le_bytes().to_vec()
}

pub fn uint64_from_string(value: &str) -> u64 {
    value.parse::<u64>().unwrap()
}

pub fn uint64_to_string(value: u64) -> String {
    value.to_string()
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

pub struct ScImmutableUint64 {
    proxy: Proxy,
}

impl ScImmutableUint64 {
    pub fn new(proxy: Proxy) -> ScImmutableUint64 {
        ScImmutableUint64 { proxy }
    }

    pub fn exists(&self) -> bool {
        self.proxy.exists()
    }

    pub fn to_string(&self) -> String {
        uint64_to_string(self.value())
    }

    pub fn value(&self) -> u64 {
        uint64_from_bytes(&self.proxy.get())
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

// value proxy for mutable u64 in host container
pub struct ScMutableUint64 {
    proxy: Proxy,
}

impl ScMutableUint64 {
    pub fn new(proxy: Proxy) -> ScMutableUint64 {
        ScMutableUint64 { proxy }
    }

    pub fn delete(&self) {
        self.proxy.delete();
    }

    pub fn exists(&self) -> bool {
        self.proxy.exists()
    }

    pub fn set_value(&self, value: u64) {
        self.proxy.set(&uint64_to_bytes(value));
    }

    pub fn to_string(&self) -> String {
        uint64_to_string(self.value())
    }

    pub fn value(&self) -> u64 {
        uint64_from_bytes(&self.proxy.get())
    }
}
