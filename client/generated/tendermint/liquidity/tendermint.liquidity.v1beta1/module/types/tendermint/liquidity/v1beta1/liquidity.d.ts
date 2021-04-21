import { Writer, Reader } from "protobufjs/minimal";
import { Coin } from "../../../cosmos_proto/coin";
import { MsgDepositWithinBatch, MsgWithdrawWithinBatch, MsgSwapWithinBatch } from "../../../tendermint/liquidity/v1beta1/tx";
export declare const protobufPackage = "tendermint.liquidity.v1beta1";
/** Structure for the pool type to distinguish the characteristics of the reserve pools */
export interface PoolType {
    /**
     * The id of the pool_type to use as pool_type_id for pool creation.
     * Only pool-type-id 1 is supported
     * {"id":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}
     */
    id: number;
    /** name of the pool type */
    name: string;
    /** min number of reserveCoins for LiquidityPoolType only 2 is allowed on this spec */
    minReserveCoinNum: number;
    /** max number of reserveCoins for LiquidityPoolType only 2 is allowed on this spec */
    maxReserveCoinNum: number;
    /** description of the pool type */
    description: string;
}
/** Params defines the parameters for the liquidity module. */
export interface Params {
    /** list of available pool types */
    poolTypes: PoolType[];
    /** Minimum number of coins to be deposited to the liquidity pool upon pool creation */
    minInitDepositAmount: string;
    /** Initial mint amount of pool coin upon pool creation */
    initPoolCoinMintAmount: string;
    /** Limit the size of each liquidity pool in the beginning phase of Liquidity Module adoption to minimize risk, 0 means no limit */
    maxReserveCoinAmount: string;
    /** Fee to create a Liquidity Pool. */
    poolCreationFee: Coin[];
    /** Swap fee rate for every executed swap */
    swapFeeRate: string;
    /** Reserve coin withdrawal with less proportion by withdrawFeeRate */
    withdrawFeeRate: string;
    /** Maximum ratio of reserve coins that can be ordered at a swap order */
    maxOrderAmountRatio: string;
    /** The smallest unit batch height for every liquidity pool */
    unitBatchHeight: number;
}
/** The liquidity pool information */
export interface Pool {
    /** id of the pool */
    id: number;
    /** id of the pool_type */
    typeId: number;
    /** denoms of reserve coin pair of the pool */
    reserveCoinDenoms: string[];
    /** reserve account address of the pool */
    reserveAccountAddress: string;
    /** denom of pool coin of the pool */
    poolCoinDenom: string;
}
/** Metadata for the state of each pool for invariant checking after genesis export or import */
export interface PoolMetadata {
    /** id of the pool */
    poolId: number;
    /** pool coin issued at the pool */
    poolCoinTotalSupply: Coin | undefined;
    /** reserve coins deposited in the pool */
    reserveCoins: Coin[];
}
/** The batch or batches of a given liquidity pool that contains indexes of the deposit, withdraw, and swap messages. The index param increments by 1 if the pool id exists. */
export interface PoolBatch {
    /** id of the pool */
    poolId: number;
    /** index of this batch */
    index: number;
    /** height where this batch is begun */
    beginHeight: number;
    /** last index of DepositMsgStates */
    depositMsgIndex: number;
    /** last index of WithdrawMsgStates */
    withdrawMsgIndex: number;
    /** last index of SwapMsgStates */
    swapMsgIndex: number;
    /** true if executed, false if not executed yet */
    executed: boolean;
}
/** DepositMsgState defines the state of the deposit message that contains the state information as it is processed in the next batch or batches */
export interface DepositMsgState {
    /** height where this message is appended to the batch */
    msgHeight: number;
    /** index of this deposit message in this liquidity pool */
    msgIndex: number;
    /** true if executed on this batch, false if not executed yet */
    executed: boolean;
    /** true if executed successfully on this batch, false if failed */
    succeeded: boolean;
    /** true if ready to be deleted on kvstore, false if not ready to be deleted */
    toBeDeleted: boolean;
    /** MsgDepositWithinBatch */
    msg: MsgDepositWithinBatch | undefined;
}
/** WithdrawMsgState defines the state of the withdraw message that contains state information as it is processed in the next batch or batches */
export interface WithdrawMsgState {
    /** height where this message is appended to the batch */
    msgHeight: number;
    /** index of this withdraw message in this liquidity pool */
    msgIndex: number;
    /** true if executed on this batch, false if not executed yet */
    executed: boolean;
    /** true if executed successfully on this batch, false if failed */
    succeeded: boolean;
    /** true if ready to be deleted on kvstore, false if not ready to be deleted */
    toBeDeleted: boolean;
    /** MsgWithdrawWithinBatch */
    msg: MsgWithdrawWithinBatch | undefined;
}
/** SwapMsgState defines the state of swap message that contains state information as it is processed in the next batch or batches */
export interface SwapMsgState {
    /** height where this message is appended to the batch */
    msgHeight: number;
    /** index of this swap message in this liquidity pool */
    msgIndex: number;
    /** true if executed on this batch, false if not executed yet */
    executed: boolean;
    /** true if executed successfully on this batch, false if failed */
    succeeded: boolean;
    /** true if ready to be deleted on kvstore, false if not ready to be deleted */
    toBeDeleted: boolean;
    /** swap orders are cancelled when current height is equal or higher than ExpiryHeight */
    orderExpiryHeight: number;
    /** offer coin exchanged until now */
    exchangedOfferCoin: Coin | undefined;
    /** offer coin currently remaining to be exchanged */
    remainingOfferCoin: Coin | undefined;
    /** reserve fee for pays fee in half offer coin */
    reservedOfferCoinFee: Coin | undefined;
    /** MsgSwapWithinBatch */
    msg: MsgSwapWithinBatch | undefined;
}
export declare const PoolType: {
    encode(message: PoolType, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): PoolType;
    fromJSON(object: any): PoolType;
    toJSON(message: PoolType): unknown;
    fromPartial(object: DeepPartial<PoolType>): PoolType;
};
export declare const Params: {
    encode(message: Params, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Params;
    fromJSON(object: any): Params;
    toJSON(message: Params): unknown;
    fromPartial(object: DeepPartial<Params>): Params;
};
export declare const Pool: {
    encode(message: Pool, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): Pool;
    fromJSON(object: any): Pool;
    toJSON(message: Pool): unknown;
    fromPartial(object: DeepPartial<Pool>): Pool;
};
export declare const PoolMetadata: {
    encode(message: PoolMetadata, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): PoolMetadata;
    fromJSON(object: any): PoolMetadata;
    toJSON(message: PoolMetadata): unknown;
    fromPartial(object: DeepPartial<PoolMetadata>): PoolMetadata;
};
export declare const PoolBatch: {
    encode(message: PoolBatch, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): PoolBatch;
    fromJSON(object: any): PoolBatch;
    toJSON(message: PoolBatch): unknown;
    fromPartial(object: DeepPartial<PoolBatch>): PoolBatch;
};
export declare const DepositMsgState: {
    encode(message: DepositMsgState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): DepositMsgState;
    fromJSON(object: any): DepositMsgState;
    toJSON(message: DepositMsgState): unknown;
    fromPartial(object: DeepPartial<DepositMsgState>): DepositMsgState;
};
export declare const WithdrawMsgState: {
    encode(message: WithdrawMsgState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): WithdrawMsgState;
    fromJSON(object: any): WithdrawMsgState;
    toJSON(message: WithdrawMsgState): unknown;
    fromPartial(object: DeepPartial<WithdrawMsgState>): WithdrawMsgState;
};
export declare const SwapMsgState: {
    encode(message: SwapMsgState, writer?: Writer): Writer;
    decode(input: Reader | Uint8Array, length?: number): SwapMsgState;
    fromJSON(object: any): SwapMsgState;
    toJSON(message: SwapMsgState): unknown;
    fromPartial(object: DeepPartial<SwapMsgState>): SwapMsgState;
};
declare type Builtin = Date | Function | Uint8Array | string | number | undefined;
export declare type DeepPartial<T> = T extends Builtin ? T : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>> : T extends {} ? {
    [K in keyof T]?: DeepPartial<T[K]>;
} : Partial<T>;
export {};
