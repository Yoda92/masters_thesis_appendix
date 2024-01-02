// Code generated by schema tool; DO NOT EDIT.

// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

import * as wasmlib from '../index';
import * as wasmtypes from '../wasmtypes';

export class CoreBlobEventHandlers implements wasmlib.IEventHandlers {
    private myID: u32;
    private coreBlobHandlers: Map<string, (evt: CoreBlobEventHandlers, dec: wasmlib.WasmDecoder) => void> = new Map();

    /* eslint-disable @typescript-eslint/no-empty-function */
    store: (evt: EventStore) => void = () => {};
    /* eslint-enable @typescript-eslint/no-empty-function */

    public constructor() {
        this.myID = wasmlib.eventHandlersGenerateID();
        this.coreBlobHandlers.set('coreblob.store', (evt: CoreBlobEventHandlers, dec: wasmlib.WasmDecoder) => evt.store(new EventStore(dec)));
    }

    public callHandler(topic: string, dec: wasmlib.WasmDecoder): void {
        const handler = this.coreBlobHandlers.get(topic);
        if (handler) {
            handler(this, dec);
        }
    }

    public id(): u32 {
        return this.myID;
    }

    public onCoreBlobStore(handler: (evt: EventStore) => void): void {
        this.store = handler;
    }
}

export class EventStore {
    public readonly timestamp: u64;
    public readonly blobHash: wasmtypes.ScHash;

    public constructor(dec: wasmlib.WasmDecoder) {
        this.timestamp = wasmtypes.uint64Decode(dec);
        this.blobHash = wasmtypes.hashDecode(dec);
        dec.close();
    }
}
