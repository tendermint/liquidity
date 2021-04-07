/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { Pool, PoolBatch, Params, SwapMsgState, DepositMsgState, WithdrawMsgState, } from "../../../tendermint/liquidity/v1beta1/liquidity";
import { PageRequest, PageResponse } from "../../../cosmos_proto/pagination";
export const protobufPackage = "tendermint.liquidity.v1beta1";
const baseQueryLiquidityPoolRequest = { poolId: 0 };
export const QueryLiquidityPoolRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryLiquidityPoolRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryLiquidityPoolRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryLiquidityPoolRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        return message;
    },
};
const baseQueryLiquidityPoolResponse = {};
export const QueryLiquidityPoolResponse = {
    encode(message, writer = Writer.create()) {
        if (message.pool !== undefined) {
            Pool.encode(message.pool, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryLiquidityPoolResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pool = Pool.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryLiquidityPoolResponse,
        };
        if (object.pool !== undefined && object.pool !== null) {
            message.pool = Pool.fromJSON(object.pool);
        }
        else {
            message.pool = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.pool !== undefined &&
            (obj.pool = message.pool ? Pool.toJSON(message.pool) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryLiquidityPoolResponse,
        };
        if (object.pool !== undefined && object.pool !== null) {
            message.pool = Pool.fromPartial(object.pool);
        }
        else {
            message.pool = undefined;
        }
        return message;
    },
};
const baseQueryLiquidityPoolBatchRequest = { poolId: 0 };
export const QueryLiquidityPoolBatchRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryLiquidityPoolBatchRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryLiquidityPoolBatchRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryLiquidityPoolBatchRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        return message;
    },
};
const baseQueryLiquidityPoolBatchResponse = {};
export const QueryLiquidityPoolBatchResponse = {
    encode(message, writer = Writer.create()) {
        if (message.batch !== undefined) {
            PoolBatch.encode(message.batch, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryLiquidityPoolBatchResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.batch = PoolBatch.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryLiquidityPoolBatchResponse,
        };
        if (object.batch !== undefined && object.batch !== null) {
            message.batch = PoolBatch.fromJSON(object.batch);
        }
        else {
            message.batch = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.batch !== undefined &&
            (obj.batch = message.batch ? PoolBatch.toJSON(message.batch) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryLiquidityPoolBatchResponse,
        };
        if (object.batch !== undefined && object.batch !== null) {
            message.batch = PoolBatch.fromPartial(object.batch);
        }
        else {
            message.batch = undefined;
        }
        return message;
    },
};
const baseQueryLiquidityPoolsRequest = {};
export const QueryLiquidityPoolsRequest = {
    encode(message, writer = Writer.create()) {
        if (message.pagination !== undefined) {
            PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryLiquidityPoolsRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pagination = PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryLiquidityPoolsRequest,
        };
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryLiquidityPoolsRequest,
        };
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryLiquidityPoolsResponse = {};
export const QueryLiquidityPoolsResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.pools) {
            Pool.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryLiquidityPoolsResponse,
        };
        message.pools = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.pools.push(Pool.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryLiquidityPoolsResponse,
        };
        message.pools = [];
        if (object.pools !== undefined && object.pools !== null) {
            for (const e of object.pools) {
                message.pools.push(Pool.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.pools) {
            obj.pools = message.pools.map((e) => (e ? Pool.toJSON(e) : undefined));
        }
        else {
            obj.pools = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryLiquidityPoolsResponse,
        };
        message.pools = [];
        if (object.pools !== undefined && object.pools !== null) {
            for (const e of object.pools) {
                message.pools.push(Pool.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryParamsRequest = {};
export const QueryParamsRequest = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryParamsRequest };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(_) {
        const message = { ...baseQueryParamsRequest };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseQueryParamsRequest };
        return message;
    },
};
const baseQueryParamsResponse = {};
export const QueryParamsResponse = {
    encode(message, writer = Writer.create()) {
        if (message.params !== undefined) {
            Params.encode(message.params, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseQueryParamsResponse };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.params = Params.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseQueryParamsResponse };
        if (object.params !== undefined && object.params !== null) {
            message.params = Params.fromJSON(object.params);
        }
        else {
            message.params = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.params !== undefined &&
            (obj.params = message.params ? Params.toJSON(message.params) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseQueryParamsResponse };
        if (object.params !== undefined && object.params !== null) {
            message.params = Params.fromPartial(object.params);
        }
        else {
            message.params = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchSwapMsgsRequest = { poolId: 0 };
export const QueryPoolBatchSwapMsgsRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        if (message.pagination !== undefined) {
            PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchSwapMsgsRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.pagination = PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgsRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgsRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchSwapMsgRequest = { poolId: 0, msgIndex: 0 };
export const QueryPoolBatchSwapMsgRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        if (message.msgIndex !== 0) {
            writer.uint32(16).uint64(message.msgIndex);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchSwapMsgRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.msgIndex = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.msgIndex !== undefined && object.msgIndex !== null) {
            message.msgIndex = Number(object.msgIndex);
        }
        else {
            message.msgIndex = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.msgIndex !== undefined && (obj.msgIndex = message.msgIndex);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.msgIndex !== undefined && object.msgIndex !== null) {
            message.msgIndex = object.msgIndex;
        }
        else {
            message.msgIndex = 0;
        }
        return message;
    },
};
const baseQueryPoolBatchSwapMsgsResponse = {};
export const QueryPoolBatchSwapMsgsResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.swaps) {
            SwapMsgState.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchSwapMsgsResponse,
        };
        message.swaps = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.swaps.push(SwapMsgState.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgsResponse,
        };
        message.swaps = [];
        if (object.swaps !== undefined && object.swaps !== null) {
            for (const e of object.swaps) {
                message.swaps.push(SwapMsgState.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.swaps) {
            obj.swaps = message.swaps.map((e) => e ? SwapMsgState.toJSON(e) : undefined);
        }
        else {
            obj.swaps = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgsResponse,
        };
        message.swaps = [];
        if (object.swaps !== undefined && object.swaps !== null) {
            for (const e of object.swaps) {
                message.swaps.push(SwapMsgState.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchSwapMsgResponse = {};
export const QueryPoolBatchSwapMsgResponse = {
    encode(message, writer = Writer.create()) {
        if (message.swap !== undefined) {
            SwapMsgState.encode(message.swap, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchSwapMsgResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.swap = SwapMsgState.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgResponse,
        };
        if (object.swap !== undefined && object.swap !== null) {
            message.swap = SwapMsgState.fromJSON(object.swap);
        }
        else {
            message.swap = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.swap !== undefined &&
            (obj.swap = message.swap ? SwapMsgState.toJSON(message.swap) : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchSwapMsgResponse,
        };
        if (object.swap !== undefined && object.swap !== null) {
            message.swap = SwapMsgState.fromPartial(object.swap);
        }
        else {
            message.swap = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchDepositMsgsRequest = { poolId: 0 };
export const QueryPoolBatchDepositMsgsRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        if (message.pagination !== undefined) {
            PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchDepositMsgsRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.pagination = PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgsRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgsRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchDepositMsgRequest = { poolId: 0, msgIndex: 0 };
export const QueryPoolBatchDepositMsgRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        if (message.msgIndex !== 0) {
            writer.uint32(16).uint64(message.msgIndex);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchDepositMsgRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.msgIndex = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.msgIndex !== undefined && object.msgIndex !== null) {
            message.msgIndex = Number(object.msgIndex);
        }
        else {
            message.msgIndex = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.msgIndex !== undefined && (obj.msgIndex = message.msgIndex);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.msgIndex !== undefined && object.msgIndex !== null) {
            message.msgIndex = object.msgIndex;
        }
        else {
            message.msgIndex = 0;
        }
        return message;
    },
};
const baseQueryPoolBatchDepositMsgsResponse = {};
export const QueryPoolBatchDepositMsgsResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.deposits) {
            DepositMsgState.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchDepositMsgsResponse,
        };
        message.deposits = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.deposits.push(DepositMsgState.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgsResponse,
        };
        message.deposits = [];
        if (object.deposits !== undefined && object.deposits !== null) {
            for (const e of object.deposits) {
                message.deposits.push(DepositMsgState.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.deposits) {
            obj.deposits = message.deposits.map((e) => e ? DepositMsgState.toJSON(e) : undefined);
        }
        else {
            obj.deposits = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgsResponse,
        };
        message.deposits = [];
        if (object.deposits !== undefined && object.deposits !== null) {
            for (const e of object.deposits) {
                message.deposits.push(DepositMsgState.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchDepositMsgResponse = {};
export const QueryPoolBatchDepositMsgResponse = {
    encode(message, writer = Writer.create()) {
        if (message.deposit !== undefined) {
            DepositMsgState.encode(message.deposit, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchDepositMsgResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.deposit = DepositMsgState.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgResponse,
        };
        if (object.deposit !== undefined && object.deposit !== null) {
            message.deposit = DepositMsgState.fromJSON(object.deposit);
        }
        else {
            message.deposit = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.deposit !== undefined &&
            (obj.deposit = message.deposit
                ? DepositMsgState.toJSON(message.deposit)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchDepositMsgResponse,
        };
        if (object.deposit !== undefined && object.deposit !== null) {
            message.deposit = DepositMsgState.fromPartial(object.deposit);
        }
        else {
            message.deposit = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchWithdrawMsgsRequest = { poolId: 0 };
export const QueryPoolBatchWithdrawMsgsRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        if (message.pagination !== undefined) {
            PageRequest.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchWithdrawMsgsRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.pagination = PageRequest.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgsRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageRequest.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgsRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageRequest.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchWithdrawMsgRequest = { poolId: 0, msgIndex: 0 };
export const QueryPoolBatchWithdrawMsgRequest = {
    encode(message, writer = Writer.create()) {
        if (message.poolId !== 0) {
            writer.uint32(8).uint64(message.poolId);
        }
        if (message.msgIndex !== 0) {
            writer.uint32(16).uint64(message.msgIndex);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchWithdrawMsgRequest,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 2:
                    message.msgIndex = longToNumber(reader.uint64());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.msgIndex !== undefined && object.msgIndex !== null) {
            message.msgIndex = Number(object.msgIndex);
        }
        else {
            message.msgIndex = 0;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.msgIndex !== undefined && (obj.msgIndex = message.msgIndex);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgRequest,
        };
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.msgIndex !== undefined && object.msgIndex !== null) {
            message.msgIndex = object.msgIndex;
        }
        else {
            message.msgIndex = 0;
        }
        return message;
    },
};
const baseQueryPoolBatchWithdrawMsgsResponse = {};
export const QueryPoolBatchWithdrawMsgsResponse = {
    encode(message, writer = Writer.create()) {
        for (const v of message.withdraws) {
            WithdrawMsgState.encode(v, writer.uint32(10).fork()).ldelim();
        }
        if (message.pagination !== undefined) {
            PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchWithdrawMsgsResponse,
        };
        message.withdraws = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.withdraws.push(WithdrawMsgState.decode(reader, reader.uint32()));
                    break;
                case 2:
                    message.pagination = PageResponse.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgsResponse,
        };
        message.withdraws = [];
        if (object.withdraws !== undefined && object.withdraws !== null) {
            for (const e of object.withdraws) {
                message.withdraws.push(WithdrawMsgState.fromJSON(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromJSON(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        if (message.withdraws) {
            obj.withdraws = message.withdraws.map((e) => e ? WithdrawMsgState.toJSON(e) : undefined);
        }
        else {
            obj.withdraws = [];
        }
        message.pagination !== undefined &&
            (obj.pagination = message.pagination
                ? PageResponse.toJSON(message.pagination)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgsResponse,
        };
        message.withdraws = [];
        if (object.withdraws !== undefined && object.withdraws !== null) {
            for (const e of object.withdraws) {
                message.withdraws.push(WithdrawMsgState.fromPartial(e));
            }
        }
        if (object.pagination !== undefined && object.pagination !== null) {
            message.pagination = PageResponse.fromPartial(object.pagination);
        }
        else {
            message.pagination = undefined;
        }
        return message;
    },
};
const baseQueryPoolBatchWithdrawMsgResponse = {};
export const QueryPoolBatchWithdrawMsgResponse = {
    encode(message, writer = Writer.create()) {
        if (message.withdraw !== undefined) {
            WithdrawMsgState.encode(message.withdraw, writer.uint32(10).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseQueryPoolBatchWithdrawMsgResponse,
        };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.withdraw = WithdrawMsgState.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgResponse,
        };
        if (object.withdraw !== undefined && object.withdraw !== null) {
            message.withdraw = WithdrawMsgState.fromJSON(object.withdraw);
        }
        else {
            message.withdraw = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.withdraw !== undefined &&
            (obj.withdraw = message.withdraw
                ? WithdrawMsgState.toJSON(message.withdraw)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = {
            ...baseQueryPoolBatchWithdrawMsgResponse,
        };
        if (object.withdraw !== undefined && object.withdraw !== null) {
            message.withdraw = WithdrawMsgState.fromPartial(object.withdraw);
        }
        else {
            message.withdraw = undefined;
        }
        return message;
    },
};
export class QueryClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    LiquidityPools(request) {
        const data = QueryLiquidityPoolsRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "LiquidityPools", data);
        return promise.then((data) => QueryLiquidityPoolsResponse.decode(new Reader(data)));
    }
    LiquidityPool(request) {
        const data = QueryLiquidityPoolRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "LiquidityPool", data);
        return promise.then((data) => QueryLiquidityPoolResponse.decode(new Reader(data)));
    }
    LiquidityPoolBatch(request) {
        const data = QueryLiquidityPoolBatchRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "LiquidityPoolBatch", data);
        return promise.then((data) => QueryLiquidityPoolBatchResponse.decode(new Reader(data)));
    }
    PoolBatchSwapMsgs(request) {
        const data = QueryPoolBatchSwapMsgsRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "PoolBatchSwapMsgs", data);
        return promise.then((data) => QueryPoolBatchSwapMsgsResponse.decode(new Reader(data)));
    }
    PoolBatchSwapMsg(request) {
        const data = QueryPoolBatchSwapMsgRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "PoolBatchSwapMsg", data);
        return promise.then((data) => QueryPoolBatchSwapMsgResponse.decode(new Reader(data)));
    }
    PoolBatchDepositMsgs(request) {
        const data = QueryPoolBatchDepositMsgsRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "PoolBatchDepositMsgs", data);
        return promise.then((data) => QueryPoolBatchDepositMsgsResponse.decode(new Reader(data)));
    }
    PoolBatchDepositMsg(request) {
        const data = QueryPoolBatchDepositMsgRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "PoolBatchDepositMsg", data);
        return promise.then((data) => QueryPoolBatchDepositMsgResponse.decode(new Reader(data)));
    }
    PoolBatchWithdrawMsgs(request) {
        const data = QueryPoolBatchWithdrawMsgsRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "PoolBatchWithdrawMsgs", data);
        return promise.then((data) => QueryPoolBatchWithdrawMsgsResponse.decode(new Reader(data)));
    }
    PoolBatchWithdrawMsg(request) {
        const data = QueryPoolBatchWithdrawMsgRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "PoolBatchWithdrawMsg", data);
        return promise.then((data) => QueryPoolBatchWithdrawMsgResponse.decode(new Reader(data)));
    }
    Params(request) {
        const data = QueryParamsRequest.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Query", "Params", data);
        return promise.then((data) => QueryParamsResponse.decode(new Reader(data)));
    }
}
var globalThis = (() => {
    if (typeof globalThis !== "undefined")
        return globalThis;
    if (typeof self !== "undefined")
        return self;
    if (typeof window !== "undefined")
        return window;
    if (typeof global !== "undefined")
        return global;
    throw "Unable to locate global object";
})();
function longToNumber(long) {
    if (long.gt(Number.MAX_SAFE_INTEGER)) {
        throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
    }
    return long.toNumber();
}
if (util.Long !== Long) {
    util.Long = Long;
    configure();
}
