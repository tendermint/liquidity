/* eslint-disable */
import {
  Pool,
  PoolMetadata,
  PoolBatch,
  DepositMsgState,
  WithdrawMsgState,
  SwapMsgState,
  Params,
} from "../../../tendermint/liquidity/v1beta1/liquidity";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "tendermint.liquidity.v1beta1";

/** record of the state of each pool after genesis export or import, used to check variables */
export interface PoolRecord {
  pool: Pool | undefined;
  poolMetadata: PoolMetadata | undefined;
  poolBatch: PoolBatch | undefined;
  depositMsgStates: DepositMsgState[];
  withdrawMsgStates: WithdrawMsgState[];
  swapMsgStates: SwapMsgState[];
}

/** GenesisState defines the liquidity module's genesis state. */
export interface GenesisState {
  /** params defines all the parameters of related to liquidity. */
  params: Params | undefined;
  poolRecords: PoolRecord[];
}

const basePoolRecord: object = {};

export const PoolRecord = {
  encode(message: PoolRecord, writer: Writer = Writer.create()): Writer {
    if (message.pool !== undefined) {
      Pool.encode(message.pool, writer.uint32(10).fork()).ldelim();
    }
    if (message.poolMetadata !== undefined) {
      PoolMetadata.encode(
        message.poolMetadata,
        writer.uint32(18).fork()
      ).ldelim();
    }
    if (message.poolBatch !== undefined) {
      PoolBatch.encode(message.poolBatch, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.depositMsgStates) {
      DepositMsgState.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    for (const v of message.withdrawMsgStates) {
      WithdrawMsgState.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.swapMsgStates) {
      SwapMsgState.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PoolRecord {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePoolRecord } as PoolRecord;
    message.depositMsgStates = [];
    message.withdrawMsgStates = [];
    message.swapMsgStates = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pool = Pool.decode(reader, reader.uint32());
          break;
        case 2:
          message.poolMetadata = PoolMetadata.decode(reader, reader.uint32());
          break;
        case 3:
          message.poolBatch = PoolBatch.decode(reader, reader.uint32());
          break;
        case 4:
          message.depositMsgStates.push(
            DepositMsgState.decode(reader, reader.uint32())
          );
          break;
        case 5:
          message.withdrawMsgStates.push(
            WithdrawMsgState.decode(reader, reader.uint32())
          );
          break;
        case 6:
          message.swapMsgStates.push(
            SwapMsgState.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PoolRecord {
    const message = { ...basePoolRecord } as PoolRecord;
    message.depositMsgStates = [];
    message.withdrawMsgStates = [];
    message.swapMsgStates = [];
    if (object.pool !== undefined && object.pool !== null) {
      message.pool = Pool.fromJSON(object.pool);
    } else {
      message.pool = undefined;
    }
    if (object.poolMetadata !== undefined && object.poolMetadata !== null) {
      message.poolMetadata = PoolMetadata.fromJSON(object.poolMetadata);
    } else {
      message.poolMetadata = undefined;
    }
    if (object.poolBatch !== undefined && object.poolBatch !== null) {
      message.poolBatch = PoolBatch.fromJSON(object.poolBatch);
    } else {
      message.poolBatch = undefined;
    }
    if (
      object.depositMsgStates !== undefined &&
      object.depositMsgStates !== null
    ) {
      for (const e of object.depositMsgStates) {
        message.depositMsgStates.push(DepositMsgState.fromJSON(e));
      }
    }
    if (
      object.withdrawMsgStates !== undefined &&
      object.withdrawMsgStates !== null
    ) {
      for (const e of object.withdrawMsgStates) {
        message.withdrawMsgStates.push(WithdrawMsgState.fromJSON(e));
      }
    }
    if (object.swapMsgStates !== undefined && object.swapMsgStates !== null) {
      for (const e of object.swapMsgStates) {
        message.swapMsgStates.push(SwapMsgState.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: PoolRecord): unknown {
    const obj: any = {};
    message.pool !== undefined &&
      (obj.pool = message.pool ? Pool.toJSON(message.pool) : undefined);
    message.poolMetadata !== undefined &&
      (obj.poolMetadata = message.poolMetadata
        ? PoolMetadata.toJSON(message.poolMetadata)
        : undefined);
    message.poolBatch !== undefined &&
      (obj.poolBatch = message.poolBatch
        ? PoolBatch.toJSON(message.poolBatch)
        : undefined);
    if (message.depositMsgStates) {
      obj.depositMsgStates = message.depositMsgStates.map((e) =>
        e ? DepositMsgState.toJSON(e) : undefined
      );
    } else {
      obj.depositMsgStates = [];
    }
    if (message.withdrawMsgStates) {
      obj.withdrawMsgStates = message.withdrawMsgStates.map((e) =>
        e ? WithdrawMsgState.toJSON(e) : undefined
      );
    } else {
      obj.withdrawMsgStates = [];
    }
    if (message.swapMsgStates) {
      obj.swapMsgStates = message.swapMsgStates.map((e) =>
        e ? SwapMsgState.toJSON(e) : undefined
      );
    } else {
      obj.swapMsgStates = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<PoolRecord>): PoolRecord {
    const message = { ...basePoolRecord } as PoolRecord;
    message.depositMsgStates = [];
    message.withdrawMsgStates = [];
    message.swapMsgStates = [];
    if (object.pool !== undefined && object.pool !== null) {
      message.pool = Pool.fromPartial(object.pool);
    } else {
      message.pool = undefined;
    }
    if (object.poolMetadata !== undefined && object.poolMetadata !== null) {
      message.poolMetadata = PoolMetadata.fromPartial(object.poolMetadata);
    } else {
      message.poolMetadata = undefined;
    }
    if (object.poolBatch !== undefined && object.poolBatch !== null) {
      message.poolBatch = PoolBatch.fromPartial(object.poolBatch);
    } else {
      message.poolBatch = undefined;
    }
    if (
      object.depositMsgStates !== undefined &&
      object.depositMsgStates !== null
    ) {
      for (const e of object.depositMsgStates) {
        message.depositMsgStates.push(DepositMsgState.fromPartial(e));
      }
    }
    if (
      object.withdrawMsgStates !== undefined &&
      object.withdrawMsgStates !== null
    ) {
      for (const e of object.withdrawMsgStates) {
        message.withdrawMsgStates.push(WithdrawMsgState.fromPartial(e));
      }
    }
    if (object.swapMsgStates !== undefined && object.swapMsgStates !== null) {
      for (const e of object.swapMsgStates) {
        message.swapMsgStates.push(SwapMsgState.fromPartial(e));
      }
    }
    return message;
  },
};

const baseGenesisState: object = {};

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.poolRecords) {
      PoolRecord.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.poolRecords = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.poolRecords.push(PoolRecord.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.poolRecords = [];
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    if (object.poolRecords !== undefined && object.poolRecords !== null) {
      for (const e of object.poolRecords) {
        message.poolRecords.push(PoolRecord.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.poolRecords) {
      obj.poolRecords = message.poolRecords.map((e) =>
        e ? PoolRecord.toJSON(e) : undefined
      );
    } else {
      obj.poolRecords = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.poolRecords = [];
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    if (object.poolRecords !== undefined && object.poolRecords !== null) {
      for (const e of object.poolRecords) {
        message.poolRecords.push(PoolRecord.fromPartial(e));
      }
    }
    return message;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;
