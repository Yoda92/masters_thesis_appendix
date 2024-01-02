// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

export const ScUint32Length = 4;

import {panic} from '../sandbox';
import {uintFromString, WasmDecoder, WasmEncoder} from './codec';
import {Proxy} from './proxy';

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

export function uint32Decode(dec: WasmDecoder): u32 {
    return uint32FromBytes(dec.fixedBytes(ScUint32Length));
}

export function uint32Encode(enc: WasmEncoder, value: u32): void {
    enc.fixedBytes(uint32ToBytes(value), ScUint32Length);
}

export function uint32FromBytes(buf: Uint8Array): u32 {
    if (buf.length == 0) {
        return 0;
    }
    if (buf.length != ScUint32Length) {
        panic('invalid Uint32 length');
    }
    // note: stupidly << 8 can result in a negative number, so use *256 instead
    let ret: u32 = buf[3];
    ret = ret * 256 + buf[2];
    ret = ret * 256 + buf[1];
    ret = ret * 256 + buf[0];
    return ret;
}

export function uint32ToBytes(value: u32): Uint8Array {
    const buf = new Uint8Array(ScUint32Length);
    buf[0] = value as u8;
    buf[1] = (value >> 8) as u8;
    buf[2] = (value >> 16) as u8;
    buf[3] = (value >> 24) as u8;
    return buf;
}

export function uint32FromString(value: string): u32 {
    return uintFromString(value, 32) as u32;
}

export function uint32ToString(value: u32): string {
    return value.toString();
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

export class ScImmutableUint32 {
    proxy: Proxy;

    constructor(proxy: Proxy) {
        this.proxy = proxy;
    }

    exists(): bool {
        return this.proxy.exists();
    }

    toString(): string {
        return uint32ToString(this.value());
    }

    value(): u32 {
        return uint32FromBytes(this.proxy.get());
    }
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

export class ScMutableUint32 extends ScImmutableUint32 {
    delete(): void {
        this.proxy.delete();
    }

    setValue(value: u32): void {
        this.proxy.set(uint32ToBytes(value));
    }
}
