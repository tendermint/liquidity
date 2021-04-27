import { Pool, PoolMetadata, PoolBatch, DepositMsgState, WithdrawMsgState, SwapMsgState, Params } from "../../../tendermint/liquidity/v1beta1/liquidity";
import { Writer, Reader } from "protobufjs/minimal";
export declare const protobufPackage = "tendermint.liquidity.v1beta1";
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
export declare const PoolRecord: {
    encode(message: PoolRecord, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): PoolRecord;
    fromJSON(object: any): PoolRecord;
    toJSON(message: PoolRecord): unknown;
    fromPartial(object: DeepPartial<PoolRecord>): PoolRecord;
};
export declare const GenesisState: {
    encode(message: GenesisState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): GenesisState;
    fromJSON(object: any): GenesisState;
    toJSON(message: GenesisState): unknown;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
