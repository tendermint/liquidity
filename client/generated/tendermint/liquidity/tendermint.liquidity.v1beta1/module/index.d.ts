import { StdFee } from "@cosmjs/launchpad";
import { OfflineSigner, EncodeObject } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgDepositWithinBatch } from "./types/tendermint/liquidity/v1beta1/tx";
import { MsgWithdrawWithinBatch } from "./types/tendermint/liquidity/v1beta1/tx";
import { MsgSwapWithinBatch } from "./types/tendermint/liquidity/v1beta1/tx";
import { MsgCreatePool } from "./types/tendermint/liquidity/v1beta1/tx";
export declare const MissingWalletError: Error;
interface TxClientOptions {
    addr: string;
}
interface SignAndBroadcastOptions {
    fee: StdFee;
    memo?: string;
}
declare const txClient: (wallet: OfflineSigner, { addr: addr }?: TxClientOptions) => Promise<{
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }?: SignAndBroadcastOptions) => Promise<import("@cosmjs/stargate").BroadcastTxResponse>;
    msgDepositWithinBatch: (data: MsgDepositWithinBatch) => EncodeObject;
    msgWithdrawWithinBatch: (data: MsgWithdrawWithinBatch) => EncodeObject;
    msgSwapWithinBatch: (data: MsgSwapWithinBatch) => EncodeObject;
    msgCreatePool: (data: MsgCreatePool) => EncodeObject;
}>;
interface QueryClientOptions {
    addr: string;
}
declare const queryClient: ({ addr: addr }?: QueryClientOptions) => Promise<Api<unknown>>;
export { txClient, queryClient, };
