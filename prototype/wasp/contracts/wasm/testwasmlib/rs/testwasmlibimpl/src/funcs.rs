// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

use crate::*;
use testwasmlib::*;
use wasmlib::*;

pub fn func_param_types(ctx: &ScFuncContext, f: &ParamTypesContext) {
    if f.params.address().exists() {
        ctx.require(
            f.params.address().value() == ctx.account_id().address(),
            "mismatch: Address",
        );
    }
    if f.params.agent_id().exists() {
        ctx.require(
            f.params.agent_id().value() == ctx.account_id(),
            "mismatch: AgentID",
        );
    }
    if f.params.big_int().exists() {
        let big_int_data = big_int_from_string("100000000000000000000");
        ctx.require(
            f.params.big_int().value().cmp(&big_int_data) == 0,
            "mismatch: BigInt",
        );
    }
    if f.params.bool().exists() {
        ctx.require(f.params.bool().value(), "mismatch: Bool");
    }
    if f.params.bytes().exists() {
        let byte_data = "these are bytes".as_bytes();
        ctx.require(f.params.bytes().value() == byte_data, "mismatch: Bytes");
    }
    if f.params.chain_id().exists() {
        ctx.require(
            f.params.chain_id().value() == ctx.current_chain_id(),
            "mismatch: ChainID",
        );
    }
    if f.params.hash().exists() {
        let hash = hash_from_bytes("0123456789abcdeffedcba9876543210".as_bytes());
        ctx.require(f.params.hash().value() == hash, "mismatch: Hash");
    }
    if f.params.hname().exists() {
        ctx.require(
            f.params.hname().value() == ctx.account_id().hname(),
            "mismatch: Hname",
        );
    }
    if f.params.int8().exists() {
        ctx.require(f.params.int8().value() == -123, "mismatch: Int8");
    }
    if f.params.int16().exists() {
        ctx.require(f.params.int16().value() == -12345, "mismatch: Int16");
    }
    if f.params.int32().exists() {
        ctx.require(f.params.int32().value() == -1234567890, "mismatch: Int32");
    }
    if f.params.int64().exists() {
        ctx.require(
            f.params.int64().value() == -1234567890123456789,
            "mismatch: Int64",
        );
    }
    if f.params.nft_id().exists() {
        let nft_id = nft_id_from_bytes("abcdefghijklmnopqrstuvwxyz123456".as_bytes());
        ctx.require(f.params.nft_id().value() == nft_id, "mismatch: NftID");
    }
    if f.params.request_id().exists() {
        let request_id =
            request_id_from_bytes("abcdefghijklmnopqrstuvwxyz123456\x00\x00".as_bytes());
        ctx.require(
            f.params.request_id().value() == request_id,
            "mismatch: RequestID",
        );
    }
    if f.params.string().exists() {
        ctx.require(
            f.params.string().value() == "this is a string",
            "mismatch: String",
        );
    }
    if f.params.token_id().exists() {
        let token_id = token_id_from_bytes("abcdefghijklmnopqrstuvwxyz1234567890AB".as_bytes());
        ctx.require(f.params.token_id().value() == token_id, "mismatch: TokenID");
    }
    if f.params.uint8().exists() {
        ctx.require(f.params.uint8().value() == 123, "mismatch: Uint8");
    }
    if f.params.uint16().exists() {
        ctx.require(f.params.uint16().value() == 12345, "mismatch: Uint16");
    }
    if f.params.uint32().exists() {
        ctx.require(f.params.uint32().value() == 1234567890, "mismatch: Uint32");
    }
    if f.params.uint64().exists() {
        ctx.require(
            f.params.uint64().value() == 1234567890123456789,
            "mismatch: Uint64",
        );
    }
}

pub fn func_random(ctx: &ScFuncContext, f: &RandomContext) {
    f.state.random().set_value(ctx.random(1000));
}

pub fn func_take_allowance(ctx: &ScFuncContext, _f: &TakeAllowanceContext) {
    ctx.transfer_allowed(
        &ctx.account_id(),
        &ScTransfer::from_balances(&ctx.allowance()),
    );
}

pub fn func_take_balance(ctx: &ScFuncContext, f: &TakeBalanceContext) {
    f.results.tokens().set_value(ctx.balances().base_tokens());
}

pub fn func_trigger_event(_ctx: &ScFuncContext, f: &TriggerEventContext) {
    f.events
        .test(&f.params.address().value(), &f.params.name().value());
}

pub fn view_get_random(_ctx: &ScViewContext, f: &GetRandomContext) {
    f.results.random().set_value(f.state.random().value());
}

pub fn view_token_balance(ctx: &ScViewContext, f: &TokenBalanceContext) {
    f.results.tokens().set_value(ctx.balances().base_tokens());
}

//////////////////// array of StringArray \\\\\\\\\\\\\\\\\\\\

pub fn func_array_of_string_array_append(
    _ctx: &ScFuncContext,
    f: &ArrayOfStringArrayAppendContext,
) {
    let index = f.params.index().value();
    let val_len = f.params.value().length();

    let sa: ArrayOfMutableString;
    if f.state.array_of_string_array().length() <= index {
        sa = f.state.array_of_string_array().append_string_array();
    } else {
        sa = f.state.array_of_string_array().get_string_array(index);
    }

    for i in 0..val_len {
        let elt = f.params.value().get_string(i).value();
        sa.append_string().set_value(&elt);
    }
}

pub fn func_array_of_string_array_clear(_ctx: &ScFuncContext, f: &ArrayOfStringArrayClearContext) {
    let length = f.state.array_of_string_array().length();
    for i in 0..length {
        let array = f.state.array_of_string_array().get_string_array(i);
        array.clear();
    }
    f.state.array_of_string_array().clear();
}

pub fn func_array_of_string_array_set(_ctx: &ScFuncContext, f: &ArrayOfStringArraySetContext) {
    let index0 = f.params.index0().value();
    let index1 = f.params.index1().value();
    let array = f.state.array_of_string_array().get_string_array(index0);
    let value = f.params.value().value();
    array.get_string(index1).set_value(&value);
}

pub fn view_array_of_string_array_length(
    _ctx: &ScViewContext,
    f: &ArrayOfStringArrayLengthContext,
) {
    let length = f.state.array_of_string_array().length();
    f.results.length().set_value(length);
}

pub fn view_array_of_string_array_value(_ctx: &ScViewContext, f: &ArrayOfStringArrayValueContext) {
    let index0 = f.params.index0().value();
    let index1 = f.params.index1().value();

    let elt = f
        .state
        .array_of_string_array()
        .get_string_array(index0)
        .get_string(index1)
        .value();
    f.results.value().set_value(&elt);
}

//////////////////// array of StringMap \\\\\\\\\\\\\\\\\\\\

pub fn func_array_of_string_map_clear(_ctx: &ScFuncContext, f: &ArrayOfStringMapClearContext) {
    let length = f.state.array_of_string_array().length();
    for i in 0..length {
        let mmap = f.state.array_of_string_map().get_string_map(i);
        mmap.clear();
    }
    f.state.array_of_string_map().clear();
}

pub fn func_array_of_string_map_set(_ctx: &ScFuncContext, f: &ArrayOfStringMapSetContext) {
    let index = f.params.index().value();
    let value = f.params.value().value();
    let key = f.params.key().value();
    if f.state.array_of_string_map().length() <= index {
        let mmap = f.state.array_of_string_map().append_string_map();
        mmap.get_string(&key).set_value(&value);
        return;
    }
    let mmap = f.state.array_of_string_map().get_string_map(index);
    mmap.get_string(&key).set_value(&value);
}

pub fn view_array_of_string_map_value(_ctx: &ScViewContext, f: &ArrayOfStringMapValueContext) {
    let index = f.params.index().value();
    let key = f.params.key().value();
    let mmap = f.state.array_of_string_map().get_string_map(index);
    f.results.value().set_value(&mmap.get_string(&key).value());
}

//////////////////// StringMap of StringArray \\\\\\\\\\\\\\\\\\\\

pub fn func_string_map_of_string_array_append(
    _ctx: &ScFuncContext,
    f: &StringMapOfStringArrayAppendContext,
) {
    let name = f.params.name().value();
    let array = f.state.string_map_of_string_array().get_string_array(&name);
    let value = f.params.value().value();
    array.append_string().set_value(&value);
}

pub fn func_string_map_of_string_array_clear(
    _ctx: &ScFuncContext,
    f: &StringMapOfStringArrayClearContext,
) {
    let name = f.params.name().value();
    let array = f.state.string_map_of_string_array().get_string_array(&name);
    array.clear();
}

pub fn func_string_map_of_string_array_set(
    _ctx: &ScFuncContext,
    f: &StringMapOfStringArraySetContext,
) {
    let name = f.params.name().value();
    let array = f.state.string_map_of_string_array().get_string_array(&name);
    let index = f.params.index().value();
    let value = f.params.value().value();
    array.get_string(index).set_value(&value);
}

pub fn view_string_map_of_string_array_length(
    _ctx: &ScViewContext,
    f: &StringMapOfStringArrayLengthContext,
) {
    let name = f.params.name().value();
    let array = f.state.string_map_of_string_array().get_string_array(&name);
    let length = array.length();
    f.results.length().set_value(length);
}

pub fn view_string_map_of_string_array_value(
    _ctx: &ScViewContext,
    f: &StringMapOfStringArrayValueContext,
) {
    let name = f.params.name().value();
    let array = f.state.string_map_of_string_array().get_string_array(&name);
    let index = f.params.index().value();
    let value = array.get_string(index).value();
    f.results.value().set_value(&value);
}

//////////////////// StringMap of StringMap \\\\\\\\\\\\\\\\\\\\

pub fn func_string_map_of_string_map_clear(
    _ctx: &ScFuncContext,
    f: &StringMapOfStringMapClearContext,
) {
    let name = f.params.name().value();
    let mmap = f.state.string_map_of_string_map().get_string_map(&name);
    mmap.clear();
}

pub fn func_string_map_of_string_map_set(_ctx: &ScFuncContext, f: &StringMapOfStringMapSetContext) {
    let name = f.params.name().value();
    let mmap = f.state.string_map_of_string_map().get_string_map(&name);
    let key = f.params.key().value();
    let value = f.params.value().value();
    mmap.get_string(&key).set_value(&value);
}

pub fn view_string_map_of_string_map_value(
    _ctx: &ScViewContext,
    f: &StringMapOfStringMapValueContext,
) {
    let name = f.params.name().value();
    let mmap = f.state.string_map_of_string_map().get_string_map(&name);
    let key = f.params.key().value();
    f.results.value().set_value(&mmap.get_string(&key).value());
}

//////////////////// array of AddressArray \\\\\\\\\\\\\\\\\\\\

pub fn func_array_of_address_array_append(
    _ctx: &ScFuncContext,
    f: &ArrayOfAddressArrayAppendContext,
) {
    let index = f.params.index().value();
    let val_len = f.params.value_addr().length();

    let sa: ArrayOfMutableAddress;
    if f.state.array_of_string_array().length() <= index {
        sa = f.state.array_of_address_array().append_address_array();
    } else {
        sa = f.state.array_of_address_array().get_address_array(index);
    }

    for i in 0..val_len {
        let elt = f.params.value_addr().get_address(i).value();
        sa.append_address().set_value(&elt);
    }
}

pub fn func_array_of_address_array_clear(
    _ctx: &ScFuncContext,
    f: &ArrayOfAddressArrayClearContext,
) {
    let length = f.state.array_of_address_array().length();
    for i in 0..length {
        let array = f.state.array_of_address_array().get_address_array(i);
        array.clear();
    }
    f.state.array_of_address_array().clear();
}

pub fn func_array_of_address_array_set(_ctx: &ScFuncContext, f: &ArrayOfAddressArraySetContext) {
    let index0 = f.params.index0().value();
    let index1 = f.params.index1().value();
    let array = f.state.array_of_address_array().get_address_array(index0);
    let value = f.params.value_addr().value();
    array.get_address(index1).set_value(&value);
}

pub fn view_array_of_address_array_length(
    _ctx: &ScViewContext,
    f: &ArrayOfAddressArrayLengthContext,
) {
    let length = f.state.array_of_address_array().length();
    f.results.length().set_value(length);
}

pub fn view_array_of_address_array_value(
    _ctx: &ScViewContext,
    f: &ArrayOfAddressArrayValueContext,
) {
    let index0 = f.params.index0().value();
    let index1 = f.params.index1().value();

    let elt = f
        .state
        .array_of_address_array()
        .get_address_array(index0)
        .get_address(index1)
        .value();
    f.results.value_addr().set_value(&elt);
}

//////////////////// array of AddressMap \\\\\\\\\\\\\\\\\\\\

pub fn func_array_of_address_map_clear(_ctx: &ScFuncContext, f: &ArrayOfAddressMapClearContext) {
    let length = f.state.array_of_address_array().length();
    for i in 0..length {
        let mmap = f.state.array_of_address_map().get_address_map(i);
        mmap.clear();
    }
    f.state.array_of_address_map().clear();
}

pub fn func_array_of_address_map_set(_ctx: &ScFuncContext, f: &ArrayOfAddressMapSetContext) {
    let index = f.params.index().value();
    let value = f.params.value_addr().value();
    let key = f.params.key_addr().value();
    if f.state.array_of_address_map().length() <= index {
        let mmap = f.state.array_of_address_map().append_address_map();
        mmap.get_address(&key).set_value(&value);
        return;
    }
    let mmap = f.state.array_of_address_map().get_address_map(index);
    mmap.get_address(&key).set_value(&value);
}

pub fn view_array_of_address_map_value(_ctx: &ScViewContext, f: &ArrayOfAddressMapValueContext) {
    let index = f.params.index().value();
    let key = f.params.key_addr().value();
    let mmap = f.state.array_of_address_map().get_address_map(index);
    f.results
        .value_addr()
        .set_value(&mmap.get_address(&key).value());
}

//////////////////// AddressMap of AddressArray \\\\\\\\\\\\\\\\\\\\

pub fn func_address_map_of_address_array_append(
    _ctx: &ScFuncContext,
    f: &AddressMapOfAddressArrayAppendContext,
) {
    let addr = f.params.name_addr().value();
    let array = f
        .state
        .address_map_of_address_array()
        .get_address_array(&addr);
    let value = f.params.value_addr().value();
    array.append_address().set_value(&value);
}

pub fn func_address_map_of_address_array_clear(
    _ctx: &ScFuncContext,
    f: &AddressMapOfAddressArrayClearContext,
) {
    let addr = f.params.name_addr().value();
    let array = f
        .state
        .address_map_of_address_array()
        .get_address_array(&addr);
    array.clear();
}

pub fn func_address_map_of_address_array_set(
    _ctx: &ScFuncContext,
    f: &AddressMapOfAddressArraySetContext,
) {
    let addr = f.params.name_addr().value();
    let array = f
        .state
        .address_map_of_address_array()
        .get_address_array(&addr);
    let index = f.params.index().value();
    let value = f.params.value_addr().value();
    array.get_address(index).set_value(&value);
}

pub fn view_address_map_of_address_array_length(
    _ctx: &ScViewContext,
    f: &AddressMapOfAddressArrayLengthContext,
) {
    let addr = f.params.name_addr().value();
    let array = f
        .state
        .address_map_of_address_array()
        .get_address_array(&addr);
    let length = array.length();
    f.results.length().set_value(length);
}

pub fn view_address_map_of_address_array_value(
    _ctx: &ScViewContext,
    f: &AddressMapOfAddressArrayValueContext,
) {
    let addr = f.params.name_addr().value();
    let array = f
        .state
        .address_map_of_address_array()
        .get_address_array(&addr);
    let index = f.params.index().value();
    let value = array.get_address(index).value();
    f.results.value_addr().set_value(&value);
}

//////////////////// AddressMap of AddressMap \\\\\\\\\\\\\\\\\\\\

pub fn func_address_map_of_address_map_clear(
    _ctx: &ScFuncContext,
    f: &AddressMapOfAddressMapClearContext,
) {
    let name = f.params.name_addr().value();
    let my_map = f.state.address_map_of_address_map().get_address_map(&name);
    my_map.clear();
}

pub fn func_address_map_of_address_map_set(
    _ctx: &ScFuncContext,
    f: &AddressMapOfAddressMapSetContext,
) {
    let name = f.params.name_addr().value();
    let my_map = f.state.address_map_of_address_map().get_address_map(&name);
    let key = f.params.key_addr().value();
    let value = f.params.value_addr().value();
    my_map.get_address(&key).set_value(&value);
}

pub fn view_address_map_of_address_map_value(
    _ctx: &ScViewContext,
    f: &AddressMapOfAddressMapValueContext,
) {
    let name = f.params.name_addr().value();
    let my_map = f.state.address_map_of_address_map().get_address_map(&name);
    let key = f.params.key_addr().value();
    f.results
        .value_addr()
        .set_value(&my_map.get_address(&key).value());
}

pub fn view_big_int_add(_ctx: &ScViewContext, f: &BigIntAddContext) {
    let lhs = f.params.lhs().value();
    let rhs = f.params.rhs().value();
    let res = lhs.add(&rhs);
    f.results.res().set_value(&res);
}

pub fn view_big_int_div(_ctx: &ScViewContext, f: &BigIntDivContext) {
    let lhs = f.params.lhs().value();
    let rhs = f.params.rhs().value();
    let res = lhs.div(&rhs);
    f.results.res().set_value(&res);
}

pub fn view_big_int_div_mod(_ctx: &ScViewContext, f: &BigIntDivModContext) {
    let lhs = f.params.lhs().value();
    let rhs = f.params.rhs().value();
    let (quo, remainder) = lhs.div_mod(&rhs);
    f.results.quo().set_value(&quo);
    f.results.remainder().set_value(&remainder);
}

pub fn view_big_int_mod(_ctx: &ScViewContext, f: &BigIntModContext) {
    let lhs = f.params.lhs().value();
    let rhs = f.params.rhs().value();
    let res = lhs.modulo(&rhs);
    f.results.res().set_value(&res);
}

pub fn view_big_int_mul(_ctx: &ScViewContext, f: &BigIntMulContext) {
    let lhs = f.params.lhs().value();
    let rhs = f.params.rhs().value();
    let res = lhs.mul(&rhs);
    f.results.res().set_value(&res);
}

pub fn view_big_int_sub(_ctx: &ScViewContext, f: &BigIntSubContext) {
    let lhs = f.params.lhs().value();
    let rhs = f.params.rhs().value();
    let res = lhs.sub(&rhs);
    f.results.res().set_value(&res);
}

pub fn view_big_int_shl(_ctx: &ScViewContext, f: &BigIntShlContext) {
    let lhs = f.params.lhs().value();
    let shift = f.params.shift().value();
    let res = lhs.shl(shift);
    f.results.res().set_value(&res);
}

pub fn view_big_int_shr(_ctx: &ScViewContext, f: &BigIntShrContext) {
    let lhs = f.params.lhs().value();
    let shift = f.params.shift().value();
    let res = lhs.shr(shift);
    f.results.res().set_value(&res);
}

pub fn view_check_agent_id(ctx: &ScViewContext, f: &CheckAgentIDContext) {
    let sc_agent_id = f.params.sc_agent_id().value();
    let agent_bytes = f.params.agent_bytes().value();
    let agent_string = f.params.agent_string().value();
    ctx.require(
        sc_agent_id == agent_id_from_bytes(&agent_id_to_bytes(&sc_agent_id)),
        "agentID bytes conversion failed",
    );
    ctx.require(
        sc_agent_id == agent_id_from_string(&agent_id_to_string(&sc_agent_id)),
        "agentID string conversion failed",
    );
    ctx.require(
        sc_agent_id.to_bytes() == agent_bytes,
        "agentID bytes mismatch",
    );
    ctx.require(
        sc_agent_id.to_string() == agent_string,
        "agentID string mismatch",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    agent_id_encode(&mut enc, &sc_agent_id);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_agent_id == agent_id_decode(&mut dec),
        "agentID encode/decode conversion failed",
    );
}

pub fn view_check_address(ctx: &ScViewContext, f: &CheckAddressContext) {
    let sc_address = f.params.sc_address().value();
    let address_bytes = f.params.address_bytes().value();
    let address_string = f.params.address_string().value();
    ctx.require(
        sc_address == address_from_bytes(&address_to_bytes(&sc_address)),
        "address bytes conversion failed",
    );
    ctx.require(
        sc_address == address_from_string(&address_to_string(&sc_address)),
        "address string conversion failed",
    );
    ctx.require(
        sc_address.to_bytes() == address_bytes,
        "address bytes mismatch",
    );
    ctx.require(
        sc_address.to_string() == address_string,
        "address string mismatch",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    address_encode(&mut enc, &sc_address);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_address == address_decode(&mut dec),
        "address encode/decode conversion failed",
    );
}

pub fn view_check_eth_address_and_agent_id(
    ctx: &ScViewContext,
    f: &CheckEthAddressAndAgentIDContext,
) {
    let address = f.params.eth_address().value();
    let address_string = f.params.eth_address_string().value();
    ctx.require(
        address.to_string() == address_string,
        "eth address string encoding failed",
    );
    ctx.require(
        address_from_string(&address_string) == address,
        "eth address string decoding failed",
    );
    ctx.require(
        address == address_from_bytes(&address_to_bytes(&address)),
        "eth address bytes conversion failed",
    );
    ctx.require(
        address == address_from_string(&address_to_string(&address)),
        "eth address to/from string conversion failed",
    );
    ctx.require(
        address_string == address_to_string(&address_from_string(&address_string)),
        "eth address from/to string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    address_encode(&mut enc, &address);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        address == address_decode(&mut dec),
        "address encode/decode conversion failed",
    );

    let agent_id = f.params.eth_agent_id().value();
    let agent_id_string = f.params.eth_agent_id_string().value();
    ctx.require(
        agent_id.to_string() == agent_id_string,
        "eth agentID string encoding failed",
    );
    ctx.require(
        agent_id_from_string(&agent_id_string) == agent_id,
        "eth agentID string decoding failed",
    );
    ctx.require(
        agent_id == agent_id_from_bytes(&agent_id_to_bytes(&agent_id)),
        "eth agentID bytes conversion failed",
    );
    ctx.require(
        agent_id == agent_id_from_string(&agent_id_to_string(&agent_id)),
        "eth agentID to/from string conversion failed",
    );
    ctx.require(
        agent_id_string == agent_id_to_string(&agent_id_from_string(&agent_id_string)),
        "eth agentID from/to string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    agent_id_encode(&mut enc, &agent_id);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        agent_id == agent_id_decode(&mut dec),
        "agent_id encode/decode conversion failed",
    );

    let agent_id_from_address = ScAgentID::for_ethereum(&agent_id.address(), &address);
    ctx.require(
        agent_id_from_address == agent_id_from_bytes(&agent_id_to_bytes(&agent_id_from_address)),
        "eth agentID bytes conversion failed",
    );
    ctx.require(
        agent_id_from_address == agent_id_from_string(&agent_id_to_string(&agent_id_from_address)),
        "eth agentID string conversion failed",
    );

    let address_from_agent_id = agent_id_from_address.address();
    ctx.require(
        address_from_agent_id == address_from_bytes(&address_to_bytes(&address_from_agent_id)),
        "eth raw agentID bytes conversion failed",
    );
    ctx.require(
        address_from_agent_id == address_from_string(&address_to_string(&address_from_agent_id)),
        "eth raw agentID string conversion failed",
    );

    let eth_address_from_agent_id = agent_id_from_address.eth_address();
    ctx.require(
        eth_address_from_agent_id == address_from_bytes(&address_to_bytes(&eth_address_from_agent_id)),
        "eth raw agentID bytes conversion failed",
    );
    ctx.require(
        eth_address_from_agent_id == address_from_string(&address_to_string(&eth_address_from_agent_id)),
        "eth raw agentID string conversion failed",
    );
}

pub fn view_check_hash(ctx: &ScViewContext, f: &CheckHashContext) {
    let sc_hash = f.params.sc_hash().value();
    let hash_bytes = f.params.hash_bytes().value();
    let hash_string = f.params.hash_string().value();
    ctx.require(
        sc_hash == hash_from_bytes(&hash_to_bytes(&sc_hash)),
        "bytes conversion failed",
    );
    ctx.require(
        sc_hash == hash_from_string(&hash_to_string(&sc_hash)),
        "string conversion failed",
    );
    ctx.require(sc_hash.to_bytes() == hash_bytes, "bytes mismatch");
    ctx.require(sc_hash.to_string() == hash_string, "string mismatch");

    let mut enc = wasmtypes::WasmEncoder::new();
    hash_encode(&mut enc, &sc_hash);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_hash == hash_decode(&mut dec),
        "hash encode/decode conversion failed",
    );
}

pub fn view_check_nft_id(ctx: &ScViewContext, f: &CheckNftIDContext) {
    let sc_nft_id = f.params.sc_nft_id().value();
    let nft_id_bytes = f.params.nft_id_bytes().value();
    let nft_id_string = f.params.nft_id_string().value();
    ctx.require(
        sc_nft_id == nft_id_from_bytes(&nft_id_to_bytes(&sc_nft_id)),
        "bytes conversion failed",
    );
    ctx.require(
        sc_nft_id == nft_id_from_string(&nft_id_to_string(&sc_nft_id)),
        "string conversion failed",
    );
    ctx.require(sc_nft_id.to_bytes() == nft_id_bytes, "bytes mismatch");
    ctx.require(sc_nft_id.to_string() == nft_id_string, "string mismatch");

    let mut enc = wasmtypes::WasmEncoder::new();
    nft_id_encode(&mut enc, &sc_nft_id);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_nft_id == nft_id_decode(&mut dec),
        "nft_id encode/decode conversion failed",
    );
}

pub fn view_check_request_id(ctx: &ScViewContext, f: &CheckRequestIDContext) {
    let sc_request_id = f.params.sc_request_id().value();
    let request_id_bytes = f.params.request_id_bytes().value();
    let request_id_string = f.params.request_id_string().value();
    ctx.require(
        sc_request_id == request_id_from_bytes(&request_id_to_bytes(&sc_request_id)),
        "bytes conversion failed",
    );
    ctx.require(
        sc_request_id == request_id_from_string(&request_id_to_string(&sc_request_id)),
        "string conversion failed",
    );
    ctx.require(
        sc_request_id.to_bytes() == request_id_bytes,
        "bytes mismatch",
    );
    ctx.require(
        sc_request_id.to_string() == request_id_string,
        "string mismatch",
    );

    let mut enc = wasmtypes::WasmEncoder::new();
    request_id_encode(&mut enc, &sc_request_id);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_request_id == request_id_decode(&mut dec),
        "request_id encode/decode conversion failed",
    );
}

pub fn view_check_token_id(ctx: &ScViewContext, f: &CheckTokenIDContext) {
    let sc_token_id = f.params.sc_token_id().value();
    let token_id_bytes = f.params.token_id_bytes().value();
    let token_id_string = f.params.token_id_string().value();
    ctx.require(
        sc_token_id == token_id_from_bytes(&token_id_to_bytes(&sc_token_id)),
        "bytes conversion failed",
    );
    ctx.require(
        sc_token_id == token_id_from_string(&token_id_to_string(&sc_token_id)),
        "string conversion failed",
    );
    ctx.require(sc_token_id.to_bytes() == token_id_bytes, "bytes mismatch");
    ctx.require(
        sc_token_id.to_string() == token_id_string,
        "string mismatch",
    );

    let mut enc = wasmtypes::WasmEncoder::new();
    token_id_encode(&mut enc, &sc_token_id);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_token_id == token_id_decode(&mut dec),
        "token_id encode/decode conversion failed",
    );
}

pub fn view_check_big_int(ctx: &ScViewContext, f: &CheckBigIntContext) {
    let sc_big_int = f.params.sc_big_int().value();
    let big_int_bytes = f.params.big_int_bytes().value();
    let big_int_string = f.params.big_int_string().value();
    ctx.require(
        sc_big_int == big_int_from_bytes(&big_int_to_bytes(&sc_big_int)),
        "bytes conversion failed",
    );
    ctx.require(
        sc_big_int == big_int_from_string(&big_int_to_string(&sc_big_int)),
        "string conversion failed",
    );
    ctx.require(sc_big_int.to_bytes() == big_int_bytes, "bytes mismatch");
    ctx.require(sc_big_int.to_string() == big_int_string, "string mismatch");

    let mut enc = wasmtypes::WasmEncoder::new();
    big_int_encode(&mut enc, &sc_big_int);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_big_int == big_int_decode(&mut dec),
        "big_int encode/decode conversion failed",
    );
}

pub fn view_check_int_and_uint(ctx: &ScViewContext, _f: &CheckIntAndUintContext) {
    let mut int8 = std::i8::MAX;
    ctx.require(
        int8 == int8_from_bytes(&int8_to_bytes(int8)),
        "bytes conversion failed",
    );
    ctx.require(
        int8 == int8_from_string(&int8_to_string(int8)),
        "string conversion failed",
    );
    int8 = std::i8::MIN;
    ctx.require(
        int8 == int8_from_bytes(&int8_to_bytes(int8)),
        "bytes conversion failed",
    );
    ctx.require(
        int8 == int8_from_string(&int8_to_string(int8)),
        "string conversion failed",
    );
    int8 = 1;
    ctx.require(
        int8 == int8_from_bytes(&int8_to_bytes(int8)),
        "bytes conversion failed",
    );
    ctx.require(
        int8 == int8_from_string(&int8_to_string(int8)),
        "string conversion failed",
    );
    int8 = 0;
    ctx.require(
        int8 == int8_from_bytes(&int8_to_bytes(int8)),
        "bytes conversion failed",
    );
    ctx.require(
        int8 == int8_from_string(&int8_to_string(int8)),
        "string conversion failed",
    );
    int8 = -1;
    ctx.require(
        int8 == int8_from_bytes(&int8_to_bytes(int8)),
        "bytes conversion failed",
    );
    ctx.require(
        int8 == int8_from_string(&int8_to_string(int8)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    int8_encode(&mut enc, int8);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        int8 == int8_decode(&mut dec),
        "int8 encode/decode conversion failed",
    );
    let mut uint8 = std::u8::MAX;
    ctx.require(
        uint8 == uint8_from_bytes(&uint8_to_bytes(uint8)),
        "bytes conversion failed",
    );
    ctx.require(
        uint8 == uint8_from_string(&uint8_to_string(uint8)),
        "string conversion failed",
    );
    uint8 = std::u8::MIN;
    ctx.require(
        uint8 == uint8_from_bytes(&uint8_to_bytes(uint8)),
        "bytes conversion failed",
    );
    ctx.require(
        uint8 == uint8_from_string(&uint8_to_string(uint8)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    uint8_encode(&mut enc, uint8);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        uint8 == uint8_decode(&mut dec),
        "int8 encode/decode conversion failed",
    );

    let mut int16 = std::i16::MAX;
    ctx.require(
        int16 == int16_from_bytes(&int16_to_bytes(int16)),
        "bytes conversion failed",
    );
    ctx.require(
        int16 == int16_from_string(&int16_to_string(int16)),
        "string conversion failed",
    );
    int16 = std::i16::MIN;
    ctx.require(
        int16 == int16_from_bytes(&int16_to_bytes(int16)),
        "bytes conversion failed",
    );
    ctx.require(
        int16 == int16_from_string(&int16_to_string(int16)),
        "string conversion failed",
    );
    int16 = 1;
    ctx.require(
        int16 == int16_from_bytes(&int16_to_bytes(int16)),
        "bytes conversion failed",
    );
    ctx.require(
        int16 == int16_from_string(&int16_to_string(int16)),
        "string conversion failed",
    );
    int16 = 0;
    ctx.require(
        int16 == int16_from_bytes(&int16_to_bytes(int16)),
        "bytes conversion failed",
    );
    ctx.require(
        int16 == int16_from_string(&int16_to_string(int16)),
        "string conversion failed",
    );
    int16 = -1;
    ctx.require(
        int16 == int16_from_bytes(&int16_to_bytes(int16)),
        "bytes conversion failed",
    );
    ctx.require(
        int16 == int16_from_string(&int16_to_string(int16)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    int16_encode(&mut enc, int16);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        int16 == int16_decode(&mut dec),
        "int16 encode/decode conversion failed",
    );
    let mut uint16 = std::u16::MAX;
    ctx.require(
        uint16 == uint16_from_bytes(&uint16_to_bytes(uint16)),
        "bytes conversion failed",
    );
    ctx.require(
        uint16 == uint16_from_string(&uint16_to_string(uint16)),
        "string conversion failed",
    );
    uint16 = std::u16::MIN;
    ctx.require(
        uint16 == uint16_from_bytes(&uint16_to_bytes(uint16)),
        "bytes conversion failed",
    );
    ctx.require(
        uint16 == uint16_from_string(&uint16_to_string(uint16)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    uint16_encode(&mut enc, uint16);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        uint16 == uint16_decode(&mut dec),
        "uint16 encode/decode conversion failed",
    );

    let mut int32 = std::i32::MAX;
    ctx.require(
        int32 == int32_from_bytes(&int32_to_bytes(int32)),
        "bytes conversion failed",
    );
    ctx.require(
        int32 == int32_from_string(&int32_to_string(int32)),
        "string conversion failed",
    );
    int32 = std::i32::MIN;
    ctx.require(
        int32 == int32_from_bytes(&int32_to_bytes(int32)),
        "bytes conversion failed",
    );
    ctx.require(
        int32 == int32_from_string(&int32_to_string(int32)),
        "string conversion failed",
    );
    int32 = 1;
    ctx.require(
        int32 == int32_from_bytes(&int32_to_bytes(int32)),
        "bytes conversion failed",
    );
    ctx.require(
        int32 == int32_from_string(&int32_to_string(int32)),
        "string conversion failed",
    );
    int32 = 0;
    ctx.require(
        int32 == int32_from_bytes(&int32_to_bytes(int32)),
        "bytes conversion failed",
    );
    ctx.require(
        int32 == int32_from_string(&int32_to_string(int32)),
        "string conversion failed",
    );
    int32 = -1;
    ctx.require(
        int32 == int32_from_bytes(&int32_to_bytes(int32)),
        "bytes conversion failed",
    );
    ctx.require(
        int32 == int32_from_string(&int32_to_string(int32)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    int32_encode(&mut enc, int32);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        int32 == int32_decode(&mut dec),
        "int32 encode/decode conversion failed",
    );
    let mut uint32 = std::u32::MAX;
    ctx.require(
        uint32 == uint32_from_bytes(&uint32_to_bytes(uint32)),
        "bytes conversion failed",
    );
    ctx.require(
        uint32 == uint32_from_string(&uint32_to_string(uint32)),
        "string conversion failed",
    );
    uint32 = std::u32::MIN;
    ctx.require(
        uint32 == uint32_from_bytes(&uint32_to_bytes(uint32)),
        "bytes conversion failed",
    );
    ctx.require(
        uint32 == uint32_from_string(&uint32_to_string(uint32)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    uint32_encode(&mut enc, uint32);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        uint32 == uint32_decode(&mut dec),
        "uint32 encode/decode conversion failed",
    );

    let mut int64 = std::i64::MAX;
    ctx.require(
        int64 == int64_from_bytes(&int64_to_bytes(int64)),
        "bytes conversion failed",
    );
    ctx.require(
        int64 == int64_from_string(&int64_to_string(int64)),
        "string conversion failed",
    );
    int64 = std::i64::MIN;
    ctx.require(
        int64 == int64_from_bytes(&int64_to_bytes(int64)),
        "bytes conversion failed",
    );
    ctx.require(
        int64 == int64_from_string(&int64_to_string(int64)),
        "string conversion failed",
    );
    int64 = 1;
    ctx.require(
        int64 == int64_from_bytes(&int64_to_bytes(int64)),
        "bytes conversion failed",
    );
    ctx.require(
        int64 == int64_from_string(&int64_to_string(int64)),
        "string conversion failed",
    );
    int64 = 0;
    ctx.require(
        int64 == int64_from_bytes(&int64_to_bytes(int64)),
        "bytes conversion failed",
    );
    ctx.require(
        int64 == int64_from_string(&int64_to_string(int64)),
        "string conversion failed",
    );
    int64 = -1;
    ctx.require(
        int64 == int64_from_bytes(&int64_to_bytes(int64)),
        "bytes conversion failed",
    );
    ctx.require(
        int64 == int64_from_string(&int64_to_string(int64)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    int64_encode(&mut enc, int64);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        int64 == int64_decode(&mut dec),
        "int64 encode/decode conversion failed",
    );
    let mut uint64 = std::u64::MAX;
    ctx.require(
        uint64 == uint64_from_bytes(&uint64_to_bytes(uint64)),
        "bytes conversion failed",
    );
    ctx.require(
        uint64 == uint64_from_string(&uint64_to_string(uint64)),
        "string conversion failed",
    );
    uint64 = std::u64::MIN;
    ctx.require(
        uint64 == uint64_from_bytes(&uint64_to_bytes(uint64)),
        "bytes conversion failed",
    );
    ctx.require(
        uint64 == uint64_from_string(&uint64_to_string(uint64)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    uint64_encode(&mut enc, uint64);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        uint64 == uint64_decode(&mut dec),
        "uint64 encode/decode conversion failed",
    );
}

pub fn view_check_bool(ctx: &ScViewContext, _f: &CheckBoolContext) {
    ctx.require(
        bool_from_bytes(&bool_to_bytes(true)),
        "bytes conversion failed",
    );
    ctx.require(
        bool_from_string(&bool_to_string(true)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    bool_encode(&mut enc, true);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        bool_decode(&mut dec),
        "bool encode/decode conversion failed",
    );
    ctx.require(
        !bool_from_bytes(&bool_to_bytes(false)),
        "bytes conversion failed",
    );
    ctx.require(
        !bool_from_string(&bool_to_string(false)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    bool_encode(&mut enc, false);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        !bool_decode(&mut dec),
        "bool encode/decode conversion failed",
    );
}

pub fn view_check_bytes(ctx: &ScViewContext, f: &CheckBytesContext) {
    let byte_data: Vec<u8> = f.params.bytes().value();
    ctx.require(
        byte_data == bytes_from_bytes(&bytes_to_bytes(&byte_data)),
        "bytes conversion failed",
    );
    ctx.require(
        byte_data == bytes_from_string(&bytes_to_string(&byte_data)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    bytes_encode(&mut enc, &byte_data);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        byte_data == bytes_decode(&mut dec),
        "bytes encode/decode conversion failed",
    );
}

pub fn view_check_hname(ctx: &ScViewContext, f: &CheckHnameContext) {
    let sc_hname: ScHname = f.params.sc_hname().value();
    let hname_bytes: Vec<u8> = f.params.hname_bytes().value();
    let hname_string: String = f.params.hname_string().value();
    ctx.require(
        sc_hname == hname_from_bytes(&hname_to_bytes(sc_hname)),
        "bytes conversion failed",
    );
    ctx.require(
        sc_hname == hname_from_string(&hname_to_string(sc_hname)),
        "string conversion failed",
    );
    ctx.require(
        hname_bytes == hname_to_bytes(sc_hname),
        "bytes conversion failed",
    );
    ctx.require(
        hname_string == hname_to_string(sc_hname),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    hname_encode(&mut enc, sc_hname);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        sc_hname == hname_decode(&mut dec),
        "hname encode/decode conversion failed",
    );
}

pub fn view_check_string(ctx: &ScViewContext, f: &CheckStringContext) {
    let string_data: String = f.params.string().value();
    ctx.require(
        string_data == string_from_bytes(&string_to_bytes(&string_data)),
        "bytes conversion failed",
    );
    ctx.require(
        string_data == string_to_string(&string_from_string(&string_data)),
        "string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    string_encode(&mut enc, &string_data);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        string_data == string_decode(&mut dec),
        "string encode/decode conversion failed",
    );
}

pub fn view_check_eth_empty_address_and_agent_id(
    ctx: &ScViewContext,
    f: &CheckEthEmptyAddressAndAgentIDContext,
) {
    let address = f.params.eth_address().value();
    let address_string = "0x0";
    let address_string_long = "0x0000000000000000000000000000000000000000";
    ctx.require(
        address == address_from_bytes(&address_to_bytes(&address)),
        "eth address bytes conversion failed",
    );
    ctx.require(
        address == address_from_string(address_string),
        "eth address to/from string conversion failed",
    );
    ctx.require(
        address == address_from_string(address_string_long),
        "eth address to/from string conversion failed",
    );
    ctx.require(
        address_string_long == address_to_string(&address),
        "eth address to/from string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    address_encode(&mut enc, &address);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        address == address_decode(&mut dec),
        "eth address encode/decode conversion failed",
    );

    let agent_id = f.params.eth_agent_id().value();
    let agent_id_string = f.params.eth_agent_id_string().value();
    ctx.require(
        agent_id.to_string() == agent_id_string,
        "eth agentID string encoding failed",
    );
    ctx.require(
        agent_id == agent_id_from_bytes(&agent_id_to_bytes(&agent_id)),
        "eth agentID bytes conversion failed",
    );
    ctx.require(
        agent_id_string == agent_id_to_string(&agent_id),
        "eth agentID from/to string conversion failed",
    );
    let mut enc = wasmtypes::WasmEncoder::new();
    agent_id_encode(&mut enc, &agent_id);
    let buf = enc.buf();
    let mut dec = wasmtypes::WasmDecoder::new(&buf);
    ctx.require(
        agent_id == agent_id_decode(&mut dec),
        "eth agentID encode/decode conversion failed",
    );

    let agent_id_from_address = ScAgentID::for_ethereum(&agent_id.address(), &address);
    ctx.require(
        agent_id_from_address == agent_id_from_bytes(&agent_id_to_bytes(&agent_id_from_address)),
        "eth agentID bytes conversion failed",
    );
    ctx.require(
        agent_id_from_address == agent_id_from_string(&agent_id_to_string(&agent_id_from_address)),
        "eth agentID string conversion failed",
    );

    let address_from_agent_id = agent_id_from_address.address();
    ctx.require(
        address_from_agent_id == address_from_bytes(&address_to_bytes(&address_from_agent_id)),
        "eth raw agentID bytes conversion failed",
    );
    ctx.require(
        address_from_agent_id == address_from_string(&address_to_string(&address_from_agent_id)),
        "eth raw agentID string conversion failed",
    );

    let eth_address_from_agent_id = agent_id_from_address.eth_address();
    ctx.require(
        eth_address_from_agent_id == address_from_bytes(&address_to_bytes(&eth_address_from_agent_id)),
        "eth raw agentID bytes conversion failed",
    );
    ctx.require(
        eth_address_from_agent_id == address_from_string(&address_to_string(&eth_address_from_agent_id)),
        "eth raw agentID string conversion failed",
    );
}

pub fn view_check_eth_invalid_empty_address_from_string(
    _ctx: &ScViewContext,
    _f: &CheckEthInvalidEmptyAddressFromStringContext,
) {
    address_from_string("0x00");
}

pub fn func_activate(ctx: &ScFuncContext, f: &ActivateContext) {
    f.state.active().set_value(true);
    let deposit = ctx.allowance().base_tokens();
    let transfer = wasmlib::ScTransfer::base_tokens(deposit);
    ctx.transfer_allowed(&ctx.account_id(), &transfer);
    let delay = f.params.seconds().value();
    testwasmlib::ScFuncs::deactivate(ctx).func.delay(delay).post();
}

pub fn func_deactivate(ctx: &ScFuncContext, f: &DeactivateContext) {
    f.state.active().set_value(false);
}

pub fn view_get_active(ctx: &ScViewContext, f: &GetActiveContext) {
    f.results.active().set_value(f.state.active().value());
}
