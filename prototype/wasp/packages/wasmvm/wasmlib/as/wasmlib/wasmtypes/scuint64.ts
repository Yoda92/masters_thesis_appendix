// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

export const ScUint64Length = 8;

import {panic} from '../sandbox';
import {uintFromString, WasmDecoder, WasmEncoder} from './codec';
import {Proxy} from './proxy';

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

export function uint64Decode(dec: WasmDecoder): u64 {
    return uint64FromBytes(dec.fixedBytes(ScUint64Length));
}

export function uint64Encode(enc: WasmEncoder, value: u64): void {
    enc.fixedBytes(uint64ToBytes(value), ScUint64Length);
}

export function uint64FromBytes(buf: Uint8Array): u64 {
    if (buf.length == 0) {
        return 0;
    }
    if (buf.length != ScUint64Length) {
        panic('invalid Uint64 length');
    }
    let ret: u64 = buf[7];
    ret = (ret << 8) | buf[6];
    ret = (ret << 8) | buf[5];
    ret = (ret << 8) | buf[4];
    ret = (ret << 8) | buf[3];
    ret = (ret << 8) | buf[2];
    ret = (ret << 8) | buf[1];
    return (ret << 8) | buf[0];
}

export function uint64ToBytes(value: u64): Uint8Array {
    const buf = new Uint8Array(ScUint64Length);
    buf[0] = value as u8;
    buf[1] = (value >> 8) as u8;
    buf[2] = (value >> 16) as u8;
    buf[3] = (value >> 24) as u8;
    buf[4] = (value >> 32) as u8;
    buf[5] = (value >> 40) as u8;
    buf[6] = (value >> 48) as u8;
    buf[7] = (value >> 56) as u8;
    return buf;
}

export function uint64FromString(value: string): u64 {
    return uintFromString(value, 64);
}

export function uint64ToString(value: u64): string {
    return value.toString();
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

export class ScImmutableUint64 {
    proxy: Proxy;

    constructor(proxy: Proxy) {
        this.proxy = proxy;
    }

    exists(): bool {
        return this.proxy.exists();
    }

    toString(): string {
        return uint64ToString(this.value());
    }

    value(): u64 {
        return uint64FromBytes(this.proxy.get());
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

export class ScMutableUint64 extends ScImmutableUint64 {
    delete(): void {
        this.proxy.delete();
    }

    setValue(value: u64): void {
        this.proxy.set(uint64ToBytes(value));
    }
}
