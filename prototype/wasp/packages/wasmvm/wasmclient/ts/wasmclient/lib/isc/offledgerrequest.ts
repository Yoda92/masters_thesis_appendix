// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import {Blake2b} from '@iota/crypto.js';
import * as wasmlib from 'wasmlib';
import {KeyPair} from './keypair';
import {chainIDEncode, hnameEncode} from "wasmlib";

export class OffLedgerSignature {
    publicKey: Uint8Array;
    signature: Uint8Array;

    public constructor(publicKey: Uint8Array) {
        this.publicKey = publicKey;
        this.signature = new Uint8Array(0);
    }
}

export class OffLedgerRequest {
    chainID: wasmlib.ScChainID;
    contract: wasmlib.ScHname;
    entryPoint: wasmlib.ScHname;
    params: Uint8Array;
    signature: OffLedgerSignature = new OffLedgerSignature(new KeyPair(new Uint8Array(0)).publicKey);
    nonce: u64;
    allowance: wasmlib.ScAssets = new wasmlib.ScAssets(new Uint8Array(0));
    gasBudget: u64 = 2n ** 64n - 1n;

    public constructor(chainID: wasmlib.ScChainID, contract: wasmlib.ScHname, entryPoint: wasmlib.ScHname, params: Uint8Array, nonce: u64) {
        this.chainID = chainID;
        this.contract = contract;
        this.entryPoint = entryPoint;
        this.params = params;
        this.nonce = nonce;
    }

    public bytes(): Uint8Array {
        let enc = this.essenceEncode();
        enc.fixedBytes(this.signature.publicKey, 32);
        enc.bytes(this.signature.signature);
        return enc.buf();
    }

    public essence(): Uint8Array {
        return this.essenceEncode().buf();
    }

    private essenceEncode() : wasmlib.WasmEncoder {
        const enc = new wasmlib.WasmEncoder();
        enc.byte(1); // requestKindOffLedgerISC
        chainIDEncode(enc, this.chainID);
        hnameEncode(enc, this.contract);
        hnameEncode(enc, this.entryPoint);
        enc.fixedBytes(this.params, this.params.length);
        enc.vliEncode64(this.nonce);
        enc.vluEncode64((this.gasBudget < 0xffffffffffffffffn) ? this.gasBudget + 1n : 0n);
        const allowance = this.allowance.toBytes();
        enc.fixedBytes(allowance, allowance.length);
        return enc;
    }

    public ID(): wasmlib.ScRequestID {
        // req id is hash of req bytes with output index zero
        const hash = Blake2b.sum256(this.bytes());
        const reqId = new wasmlib.ScRequestID();
        reqId.id.set(hash, 0);
        return reqId;
    }

    public sign(keyPair: KeyPair): OffLedgerRequest {
        const req = new OffLedgerRequest(this.chainID, this.contract, this.entryPoint, this.params, this.nonce);
        req.signature = new OffLedgerSignature(keyPair.publicKey);
        req.signature.signature = keyPair.sign(req.essence());
        return req;
    }

    public withAllowance(allowance: wasmlib.ScAssets): void {
        this.allowance = allowance;
    }
}
