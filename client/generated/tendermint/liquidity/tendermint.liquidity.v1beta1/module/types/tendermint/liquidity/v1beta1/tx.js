/* eslint-disable */
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { Coin } from "../../../cosmos_proto/coin";
export const protobufPackage = "tendermint.liquidity.v1beta1";
const baseMsgCreatePool = { poolCreatorAddress: "", poolTypeId: 0 };
export const MsgCreatePool = {
    encode(message, writer = Writer.create()) {
        if (message.poolCreatorAddress !== "") {
            writer.uint32(10).string(message.poolCreatorAddress);
        }
        if (message.poolTypeId !== 0) {
            writer.uint32(16).uint32(message.poolTypeId);
        }
        for (const v of message.depositCoins) {
            Coin.encode(v, writer.uint32(34).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgCreatePool };
        message.depositCoins = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.poolCreatorAddress = reader.string();
                    break;
                case 2:
                    message.poolTypeId = reader.uint32();
                    break;
                case 4:
                    message.depositCoins.push(Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgCreatePool };
        message.depositCoins = [];
        if (object.poolCreatorAddress !== undefined &&
            object.poolCreatorAddress !== null) {
            message.poolCreatorAddress = String(object.poolCreatorAddress);
        }
        else {
            message.poolCreatorAddress = "";
        }
        if (object.poolTypeId !== undefined && object.poolTypeId !== null) {
            message.poolTypeId = Number(object.poolTypeId);
        }
        else {
            message.poolTypeId = 0;
        }
        if (object.depositCoins !== undefined && object.depositCoins !== null) {
            for (const e of object.depositCoins) {
                message.depositCoins.push(Coin.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.poolCreatorAddress !== undefined &&
            (obj.poolCreatorAddress = message.poolCreatorAddress);
        message.poolTypeId !== undefined && (obj.poolTypeId = message.poolTypeId);
        if (message.depositCoins) {
            obj.depositCoins = message.depositCoins.map((e) => e ? Coin.toJSON(e) : undefined);
        }
        else {
            obj.depositCoins = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgCreatePool };
        message.depositCoins = [];
        if (object.poolCreatorAddress !== undefined &&
            object.poolCreatorAddress !== null) {
            message.poolCreatorAddress = object.poolCreatorAddress;
        }
        else {
            message.poolCreatorAddress = "";
        }
        if (object.poolTypeId !== undefined && object.poolTypeId !== null) {
            message.poolTypeId = object.poolTypeId;
        }
        else {
            message.poolTypeId = 0;
        }
        if (object.depositCoins !== undefined && object.depositCoins !== null) {
            for (const e of object.depositCoins) {
                message.depositCoins.push(Coin.fromPartial(e));
            }
        }
        return message;
    },
};
const baseMsgCreatePoolResponse = {};
export const MsgCreatePoolResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgCreatePoolResponse };
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
        const message = { ...baseMsgCreatePoolResponse };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = { ...baseMsgCreatePoolResponse };
        return message;
    },
};
const baseMsgDepositWithinBatch = { depositorAddress: "", poolId: 0 };
export const MsgDepositWithinBatch = {
    encode(message, writer = Writer.create()) {
        if (message.depositorAddress !== "") {
            writer.uint32(10).string(message.depositorAddress);
        }
        if (message.poolId !== 0) {
            writer.uint32(16).uint64(message.poolId);
        }
        for (const v of message.depositCoins) {
            Coin.encode(v, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgDepositWithinBatch };
        message.depositCoins = [];
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.depositorAddress = reader.string();
                    break;
                case 2:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 3:
                    message.depositCoins.push(Coin.decode(reader, reader.uint32()));
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgDepositWithinBatch };
        message.depositCoins = [];
        if (object.depositorAddress !== undefined &&
            object.depositorAddress !== null) {
            message.depositorAddress = String(object.depositorAddress);
        }
        else {
            message.depositorAddress = "";
        }
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.depositCoins !== undefined && object.depositCoins !== null) {
            for (const e of object.depositCoins) {
                message.depositCoins.push(Coin.fromJSON(e));
            }
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.depositorAddress !== undefined &&
            (obj.depositorAddress = message.depositorAddress);
        message.poolId !== undefined && (obj.poolId = message.poolId);
        if (message.depositCoins) {
            obj.depositCoins = message.depositCoins.map((e) => e ? Coin.toJSON(e) : undefined);
        }
        else {
            obj.depositCoins = [];
        }
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgDepositWithinBatch };
        message.depositCoins = [];
        if (object.depositorAddress !== undefined &&
            object.depositorAddress !== null) {
            message.depositorAddress = object.depositorAddress;
        }
        else {
            message.depositorAddress = "";
        }
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.depositCoins !== undefined && object.depositCoins !== null) {
            for (const e of object.depositCoins) {
                message.depositCoins.push(Coin.fromPartial(e));
            }
        }
        return message;
    },
};
const baseMsgDepositWithinBatchResponse = {};
export const MsgDepositWithinBatchResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseMsgDepositWithinBatchResponse,
        };
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
        const message = {
            ...baseMsgDepositWithinBatchResponse,
        };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = {
            ...baseMsgDepositWithinBatchResponse,
        };
        return message;
    },
};
const baseMsgWithdrawWithinBatch = { withdrawerAddress: "", poolId: 0 };
export const MsgWithdrawWithinBatch = {
    encode(message, writer = Writer.create()) {
        if (message.withdrawerAddress !== "") {
            writer.uint32(10).string(message.withdrawerAddress);
        }
        if (message.poolId !== 0) {
            writer.uint32(16).uint64(message.poolId);
        }
        if (message.poolCoin !== undefined) {
            Coin.encode(message.poolCoin, writer.uint32(26).fork()).ldelim();
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgWithdrawWithinBatch };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.withdrawerAddress = reader.string();
                    break;
                case 2:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 3:
                    message.poolCoin = Coin.decode(reader, reader.uint32());
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgWithdrawWithinBatch };
        if (object.withdrawerAddress !== undefined &&
            object.withdrawerAddress !== null) {
            message.withdrawerAddress = String(object.withdrawerAddress);
        }
        else {
            message.withdrawerAddress = "";
        }
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.poolCoin !== undefined && object.poolCoin !== null) {
            message.poolCoin = Coin.fromJSON(object.poolCoin);
        }
        else {
            message.poolCoin = undefined;
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.withdrawerAddress !== undefined &&
            (obj.withdrawerAddress = message.withdrawerAddress);
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.poolCoin !== undefined &&
            (obj.poolCoin = message.poolCoin
                ? Coin.toJSON(message.poolCoin)
                : undefined);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgWithdrawWithinBatch };
        if (object.withdrawerAddress !== undefined &&
            object.withdrawerAddress !== null) {
            message.withdrawerAddress = object.withdrawerAddress;
        }
        else {
            message.withdrawerAddress = "";
        }
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.poolCoin !== undefined && object.poolCoin !== null) {
            message.poolCoin = Coin.fromPartial(object.poolCoin);
        }
        else {
            message.poolCoin = undefined;
        }
        return message;
    },
};
const baseMsgWithdrawWithinBatchResponse = {};
export const MsgWithdrawWithinBatchResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseMsgWithdrawWithinBatchResponse,
        };
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
        const message = {
            ...baseMsgWithdrawWithinBatchResponse,
        };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = {
            ...baseMsgWithdrawWithinBatchResponse,
        };
        return message;
    },
};
const baseMsgSwapWithinBatch = {
    swapRequesterAddress: "",
    poolId: 0,
    swapTypeId: 0,
    demandCoinDenom: "",
    orderPrice: "",
};
export const MsgSwapWithinBatch = {
    encode(message, writer = Writer.create()) {
        if (message.swapRequesterAddress !== "") {
            writer.uint32(10).string(message.swapRequesterAddress);
        }
        if (message.poolId !== 0) {
            writer.uint32(16).uint64(message.poolId);
        }
        if (message.swapTypeId !== 0) {
            writer.uint32(24).uint32(message.swapTypeId);
        }
        if (message.offerCoin !== undefined) {
            Coin.encode(message.offerCoin, writer.uint32(34).fork()).ldelim();
        }
        if (message.demandCoinDenom !== "") {
            writer.uint32(42).string(message.demandCoinDenom);
        }
        if (message.offerCoinFee !== undefined) {
            Coin.encode(message.offerCoinFee, writer.uint32(50).fork()).ldelim();
        }
        if (message.orderPrice !== "") {
            writer.uint32(58).string(message.orderPrice);
        }
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = { ...baseMsgSwapWithinBatch };
        while (reader.pos < end) {
            const tag = reader.uint32();
            switch (tag >>> 3) {
                case 1:
                    message.swapRequesterAddress = reader.string();
                    break;
                case 2:
                    message.poolId = longToNumber(reader.uint64());
                    break;
                case 3:
                    message.swapTypeId = reader.uint32();
                    break;
                case 4:
                    message.offerCoin = Coin.decode(reader, reader.uint32());
                    break;
                case 5:
                    message.demandCoinDenom = reader.string();
                    break;
                case 6:
                    message.offerCoinFee = Coin.decode(reader, reader.uint32());
                    break;
                case 7:
                    message.orderPrice = reader.string();
                    break;
                default:
                    reader.skipType(tag & 7);
                    break;
            }
        }
        return message;
    },
    fromJSON(object) {
        const message = { ...baseMsgSwapWithinBatch };
        if (object.swapRequesterAddress !== undefined &&
            object.swapRequesterAddress !== null) {
            message.swapRequesterAddress = String(object.swapRequesterAddress);
        }
        else {
            message.swapRequesterAddress = "";
        }
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = Number(object.poolId);
        }
        else {
            message.poolId = 0;
        }
        if (object.swapTypeId !== undefined && object.swapTypeId !== null) {
            message.swapTypeId = Number(object.swapTypeId);
        }
        else {
            message.swapTypeId = 0;
        }
        if (object.offerCoin !== undefined && object.offerCoin !== null) {
            message.offerCoin = Coin.fromJSON(object.offerCoin);
        }
        else {
            message.offerCoin = undefined;
        }
        if (object.demandCoinDenom !== undefined &&
            object.demandCoinDenom !== null) {
            message.demandCoinDenom = String(object.demandCoinDenom);
        }
        else {
            message.demandCoinDenom = "";
        }
        if (object.offerCoinFee !== undefined && object.offerCoinFee !== null) {
            message.offerCoinFee = Coin.fromJSON(object.offerCoinFee);
        }
        else {
            message.offerCoinFee = undefined;
        }
        if (object.orderPrice !== undefined && object.orderPrice !== null) {
            message.orderPrice = String(object.orderPrice);
        }
        else {
            message.orderPrice = "";
        }
        return message;
    },
    toJSON(message) {
        const obj = {};
        message.swapRequesterAddress !== undefined &&
            (obj.swapRequesterAddress = message.swapRequesterAddress);
        message.poolId !== undefined && (obj.poolId = message.poolId);
        message.swapTypeId !== undefined && (obj.swapTypeId = message.swapTypeId);
        message.offerCoin !== undefined &&
            (obj.offerCoin = message.offerCoin
                ? Coin.toJSON(message.offerCoin)
                : undefined);
        message.demandCoinDenom !== undefined &&
            (obj.demandCoinDenom = message.demandCoinDenom);
        message.offerCoinFee !== undefined &&
            (obj.offerCoinFee = message.offerCoinFee
                ? Coin.toJSON(message.offerCoinFee)
                : undefined);
        message.orderPrice !== undefined && (obj.orderPrice = message.orderPrice);
        return obj;
    },
    fromPartial(object) {
        const message = { ...baseMsgSwapWithinBatch };
        if (object.swapRequesterAddress !== undefined &&
            object.swapRequesterAddress !== null) {
            message.swapRequesterAddress = object.swapRequesterAddress;
        }
        else {
            message.swapRequesterAddress = "";
        }
        if (object.poolId !== undefined && object.poolId !== null) {
            message.poolId = object.poolId;
        }
        else {
            message.poolId = 0;
        }
        if (object.swapTypeId !== undefined && object.swapTypeId !== null) {
            message.swapTypeId = object.swapTypeId;
        }
        else {
            message.swapTypeId = 0;
        }
        if (object.offerCoin !== undefined && object.offerCoin !== null) {
            message.offerCoin = Coin.fromPartial(object.offerCoin);
        }
        else {
            message.offerCoin = undefined;
        }
        if (object.demandCoinDenom !== undefined &&
            object.demandCoinDenom !== null) {
            message.demandCoinDenom = object.demandCoinDenom;
        }
        else {
            message.demandCoinDenom = "";
        }
        if (object.offerCoinFee !== undefined && object.offerCoinFee !== null) {
            message.offerCoinFee = Coin.fromPartial(object.offerCoinFee);
        }
        else {
            message.offerCoinFee = undefined;
        }
        if (object.orderPrice !== undefined && object.orderPrice !== null) {
            message.orderPrice = object.orderPrice;
        }
        else {
            message.orderPrice = "";
        }
        return message;
    },
};
const baseMsgSwapWithinBatchResponse = {};
export const MsgSwapWithinBatchResponse = {
    encode(_, writer = Writer.create()) {
        return writer;
    },
    decode(input, length) {
        const reader = input instanceof Uint8Array ? new Reader(input) : input;
        let end = length === undefined ? reader.len : reader.pos + length;
        const message = {
            ...baseMsgSwapWithinBatchResponse,
        };
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
        const message = {
            ...baseMsgSwapWithinBatchResponse,
        };
        return message;
    },
    toJSON(_) {
        const obj = {};
        return obj;
    },
    fromPartial(_) {
        const message = {
            ...baseMsgSwapWithinBatchResponse,
        };
        return message;
    },
};
export class MsgClientImpl {
    constructor(rpc) {
        this.rpc = rpc;
    }
    CreatePool(request) {
        const data = MsgCreatePool.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Msg", "CreatePool", data);
        return promise.then((data) => MsgCreatePoolResponse.decode(new Reader(data)));
    }
    DepositWithinBatch(request) {
        const data = MsgDepositWithinBatch.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Msg", "DepositWithinBatch", data);
        return promise.then((data) => MsgDepositWithinBatchResponse.decode(new Reader(data)));
    }
    WithdrawWithinBatch(request) {
        const data = MsgWithdrawWithinBatch.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Msg", "WithdrawWithinBatch", data);
        return promise.then((data) => MsgWithdrawWithinBatchResponse.decode(new Reader(data)));
    }
    Swap(request) {
        const data = MsgSwapWithinBatch.encode(request).finish();
        const promise = this.rpc.request("tendermint.liquidity.v1beta1.Msg", "Swap", data);
        return promise.then((data) => MsgSwapWithinBatchResponse.decode(new Reader(data)));
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
