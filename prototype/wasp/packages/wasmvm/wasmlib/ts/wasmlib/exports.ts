// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

// Provide host with details about funcs and views in this smart contract

import {ScFuncContext, ScViewContext} from './context';
import {exportName} from './host';

// Note that we do not use the Wasm export symbol table on purpose
// because Wasm does not allow us to determine whether the symbols
// are meant as view or func, or meant as extra public callbacks
// generated by the compilation of the the Wasm code.
// There are only 2 symbols the ISC host will actually look for
// in the export table:
// on_load (which must be defined by the SC code) and
// on_call (which is defined here as part of WasmLib)

export type ScFuncContextFunc = (f: ScFuncContext) => void;
export type ScViewContextFunc = (v: ScViewContext) => void;

// context for onLoad function to be able to tell host which
// funcs and views are available as entry points to the SC
export class ScExportMap {
    names: string[];
    funcs: ScFuncContextFunc[];
    views: ScViewContextFunc[];

    constructor(names: string[], funcs: ScFuncContextFunc[], views: ScViewContextFunc[]) {
        this.names = names;
        this.funcs = funcs;
        this.views = views;
    }

    // general entrypoint for the host to call any SC function
    // the host will pass the index of one of the entry points
    // that was provided by onLoad during SC initialization
    public dispatch(index: i32): void {
        if (index == -1) {
            // special dispatch for exporting entry points to host
            this.export();
            return;
        }

        if ((index & 0x8000) == 0) {
            // mutable full function, invoke with a WasmLib func call context
            const func = this.funcs[index];
            func(new ScFuncContext());
            return;
        }
        // immutable view function, invoke with a WasmLib view call context
        const view = this.views[index & 0x7fff];
        view(new ScViewContext());
    }

    // constructs the symbol export context for the onLoad function
    public export(): void {
        exportName(-1, 'WASM::TYPESCRIPT');

        for (let i = 0; i < this.funcs.length; i++) {
            exportName(i as i32, this.names[i]);
        }

        const offset = this.funcs.length;
        for (let i = 0; i < this.views.length; i++) {
            exportName((i as i32) | 0x8000, this.names[offset + i]);
        }
    }
}

