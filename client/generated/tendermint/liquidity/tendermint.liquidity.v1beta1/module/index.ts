// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgDepositWithinBatch } from "./types/tendermint/liquidity/v1beta1/tx";
import { MsgWithdrawWithinBatch } from "./types/tendermint/liquidity/v1beta1/tx";
import { MsgSwapWithinBatch } from "./types/tendermint/liquidity/v1beta1/tx";
import { MsgCreatePool } from "./types/tendermint/liquidity/v1beta1/tx";


const types = [
  ["/tendermint.liquidity.v1beta1.MsgDepositWithinBatch", MsgDepositWithinBatch],
  ["/tendermint.liquidity.v1beta1.MsgWithdrawWithinBatch", MsgWithdrawWithinBatch],
  ["/tendermint.liquidity.v1beta1.MsgSwapWithinBatch", MsgSwapWithinBatch],
  ["/tendermint.liquidity.v1beta1.MsgCreatePool", MsgCreatePool],
  
];
export const MissingWalletError = new Error("wallet is required");

const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw MissingWalletError;

  const client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }: SignAndBroadcastOptions = {fee: defaultFee, memo: ""}) => client.signAndBroadcast(address, msgs, fee,memo),
    msgDepositWithinBatch: (data: MsgDepositWithinBatch): EncodeObject => ({ typeUrl: "/tendermint.liquidity.v1beta1.MsgDepositWithinBatch", value: data }),
    msgWithdrawWithinBatch: (data: MsgWithdrawWithinBatch): EncodeObject => ({ typeUrl: "/tendermint.liquidity.v1beta1.MsgWithdrawWithinBatch", value: data }),
    msgSwapWithinBatch: (data: MsgSwapWithinBatch): EncodeObject => ({ typeUrl: "/tendermint.liquidity.v1beta1.MsgSwapWithinBatch", value: data }),
    msgCreatePool: (data: MsgCreatePool): EncodeObject => ({ typeUrl: "/tendermint.liquidity.v1beta1.MsgCreatePool", value: data }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
