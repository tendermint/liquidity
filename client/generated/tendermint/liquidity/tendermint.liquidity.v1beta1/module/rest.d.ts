export interface ProtobufAny {
    typeUrl?: string;
    /** @format byte */
    value?: string;
}
export interface RpcStatus {
    /** @format int32 */
    code?: number;
    message?: string;
    details?: ProtobufAny[];
}
/**
* Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.
*/
export interface V1Beta1Coin {
    denom?: string;
    amount?: string;
}
export interface V1Beta1DepositMsgState {
    /**
     * @format int64
     * @example 1000
     */
    msgHeight?: string;
    /**
     * @format uint64
     * @example 1
     */
    msgIndex?: string;
    /** @example true */
    executed?: boolean;
    /** @example true */
    succeeded?: boolean;
    /** @example true */
    toBeDeleted?: boolean;
    /**
     * `MsgDepositWithinBatch defines` an `sdk.Msg` type that supports submitting a deposit requests to the liquidity pool batch
     * The deposit is submitted with the specified `pool_id` and reserve `deposit_coins`
     * The deposit requests are stacked in the liquidity pool batch and are not immediately processed
     * Batch deposit requests are processed in the `endblock` at the same time as other requests.
     *
     * See: https://github.com/tendermint/liquidity/blob/develop/x/liquidity/spec/04_messages.md
     */
    msg?: V1Beta1MsgDepositWithinBatch;
}
/**
 * MsgCreatePoolResponse defines the Msg/CreatePool response type.
 */
export declare type V1Beta1MsgCreatePoolResponse = object;
/**
* `MsgDepositWithinBatch defines` an `sdk.Msg` type that supports submitting a deposit requests to the liquidity pool batch
The deposit is submitted with the specified `pool_id` and reserve `deposit_coins`
The deposit requests are stacked in the liquidity pool batch and are not immediately processed
Batch deposit requests are processed in the `endblock` at the same time as other requests.

See: https://github.com/tendermint/liquidity/blob/develop/x/liquidity/spec/04_messages.md
*/
export interface V1Beta1MsgDepositWithinBatch {
    /**
     * account address of the origin of this message
     * @format sdk.AccAddress
     * @example cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun
     */
    depositorAddress?: string;
    /**
     * @format uint64
     * @example 1
     */
    poolId?: string;
    /**
     * @format sdk.Coins
     * @example [{"denom":"denomX","amount":"1000000"},{"denom":"denomY","amount":"2000000"}]
     */
    depositCoins?: V1Beta1Coin[];
}
/**
 * MsgDepositWithinBatchResponse defines the Msg/DepositWithinBatch response type.
 */
export declare type V1Beta1MsgDepositWithinBatchResponse = object;
/**
* `MsgSwapWithinBatch` defines an sdk.Msg type that submits a swap offer request to the liquidity pool batch
Submit swap offer to the liquidity pool batch with the specified the `pool_id`, `swap_type_id`,
`demand_coin_denom` with the coin and the price you're offering
The `offer_coin_fee` must be half of the offer coin amount * current `params.swap_fee_rate` for reservation to pay fees
This request is added to the pool and executed at the end of the batch (`endblock`)
You must submit the request using the same fields as the pool
Only the default `swap_type_id`1 is supported
The detailed swap algorithm is shown here.

See: https://github.com/tendermint/liquidity/tree/develop/doc
https://github.com/tendermint/liquidity/blob/develop/x/liquidity/spec/04_messages.md
*/
export interface V1Beta1MsgSwapWithinBatch {
    /**
     * account address of the origin of this message
     * @format sdk.AccAddress
     * @example cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun
     */
    swapRequesterAddress?: string;
    /**
     * @format uint64
     * @example 1
     */
    poolId?: string;
    /**
     * id of swap type. Must match the value in the pool.
     * @format uint32
     * @example 1
     */
    swapTypeId?: number;
    /**
     * offer sdk.coin for the swap request, must match the denom in the pool.
     * @format sdk.Coin
     * @example {"denom":"denomX","amount":"1000000"}
     */
    offerCoin?: V1Beta1Coin;
    /**
     * denom of demand coin to be exchanged on the swap request, must match the denom in the pool.
     * @example denomB
     */
    demandCoinDenom?: string;
    /**
     * Coin defines a token with a denomination and an amount.
     *
     * NOTE: The amount field is an Int which implements the custom method
     * signatures required by gogoproto.
     * @format sdk.Coin
     * @example {"denom":"denomX","amount":"5000"}
     */
    offerCoinFee?: V1Beta1Coin;
    /**
     * @format sdk.Dec
     * @example 1.1
     */
    orderPrice?: string;
}
/**
 * MsgSwapWithinBatchResponse defines the Msg/Swap response type.
 */
export declare type V1Beta1MsgSwapWithinBatchResponse = object;
/**
* `MsgWithdrawWithinBatch` defines an `sdk.Msg` type that submits a withdraw request to the liquidity pool batch
Withdraw submit to the batch from the Liquidity pool with the specified `pool_id`, `pool_coin` of the pool
this requests are stacked in the batch of the liquidity pool, not immediately processed and
processed in the `endblock` at once with other requests.

See: https://github.com/tendermint/liquidity/blob/develop/x/liquidity/spec/04_messages.md
*/
export interface V1Beta1MsgWithdrawWithinBatch {
    /**
     * account address of the origin of this message
     * @format sdk.AccAddress
     * @example cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun
     */
    withdrawerAddress?: string;
    /**
     * @format uint64
     * @example 1
     */
    poolId?: string;
    /**
     * Coin defines a token with a denomination and an amount.
     *
     * NOTE: The amount field is an Int which implements the custom method
     * signatures required by gogoproto.
     * @format sdk.Coin
     * @example {"denom":"poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4","amount":"1000"}
     */
    poolCoin?: V1Beta1Coin;
}
/**
 * MsgWithdrawWithinBatchResponse defines the Msg/WithdrawWithinBatch response type.
 */
export declare type V1Beta1MsgWithdrawWithinBatchResponse = object;
/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
    /**
     * key is a value returned in PageResponse.next_key to begin
     * querying the next page most efficiently. Only one of offset or key
     * should be set.
     * @format byte
     */
    key?: string;
    /**
     * offset is a numeric offset that can be used when key is unavailable.
     * It is less efficient than using key. Only one of offset or key should
     * be set.
     * @format uint64
     */
    offset?: string;
    /**
     * limit is the total number of results to be returned in the result page.
     * If left empty it will default to a value to be set by each app.
     * @format uint64
     */
    limit?: string;
    /**
     * count_total is set to true  to indicate that the result set should include
     * a count of the total number of items available for pagination in UIs.
     * count_total is only respected when offset is used. It is ignored when key
     * is set.
     */
    countTotal?: boolean;
}
/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
    /** @format byte */
    nextKey?: string;
    /** @format uint64 */
    total?: string;
}
/**
 * Params defines the parameters for the liquidity module.
 */
export interface V1Beta1Params {
    poolTypes?: V1Beta1PoolType[];
    /**
     * @format sdk.Int
     * @example 1000000
     */
    minInitDepositAmount?: string;
    /**
     * @format sdk.Int
     * @example 1000000
     */
    initPoolCoinMintAmount?: string;
    /**
     * @format sdk.Int
     * @example 1000000000000
     */
    maxReserveCoinAmount?: string;
    /**
     * Fee to create a Liquidity Pool.
     * @format sdk.Coins
     * @example [{"denom":"uatom","amount":"100000000"}]
     */
    poolCreationFee?: V1Beta1Coin[];
    /**
     * @format sdk.Dec
     * @example 0.003
     */
    swapFeeRate?: string;
    /**
     * @format sdk.Dec
     * @example 0.003
     */
    withdrawFeeRate?: string;
    /**
     * @format sdk.Dec
     * @example 0.003
     */
    maxOrderAmountRatio?: string;
    /**
     * @format uint32
     * @example 1
     */
    unitBatchHeight?: number;
}
export interface V1Beta1Pool {
    /**
     * @format uint64
     * @example 1
     */
    id?: string;
    /**
     * @format uint32
     * @example 1
     */
    typeId?: number;
    /** @example ["denomX","denomY"] */
    reserveCoinDenoms?: string[];
    /**
     * @format sdk.AccAddress
     * @example cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3
     */
    reserveAccountAddress?: string;
    /** @example poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4 */
    poolCoinDenom?: string;
}
/**
 * The batch or batches of a given liquidity pool that contains indexes of the deposit, withdraw, and swap messages. The index param increments by 1 if the pool id exists.
 */
export interface V1Beta1PoolBatch {
    /**
     * @format uint64
     * @example 1
     */
    poolId?: string;
    /**
     * @format uint64
     * @example 1
     */
    index?: string;
    /**
     * @format int64
     * @example 1000
     */
    beginHeight?: string;
    /**
     * @format uint64
     * @example 1
     */
    depositMsgIndex?: string;
    /**
     * @format uint64
     * @example 1
     */
    withdrawMsgIndex?: string;
    /**
     * @format uint64
     * @example 1
     */
    swapMsgIndex?: string;
    /** @example true */
    executed?: boolean;
}
export interface V1Beta1PoolType {
    /**
     * @format uint32
     * @example 1
     */
    id?: number;
    /** @example ConstantProductLiquidityPool */
    name?: string;
    /**
     * @format uint32
     * @example 2
     */
    minReserveCoinNum?: number;
    /**
     * @format uint32
     * @example 2
     */
    maxReserveCoinNum?: number;
    description?: string;
}
/**
 * the response type for the QueryLiquidityPoolBatchResponse RPC method. It returns the liquidity pool batch corresponding to the requested pool_id.
 */
export interface V1Beta1QueryLiquidityPoolBatchResponse {
    /** The batch or batches of a given liquidity pool that contains indexes of the deposit, withdraw, and swap messages. The index param increments by 1 if the pool id exists. */
    batch?: V1Beta1PoolBatch;
}
/**
 * the response type for the QueryLiquidityPoolResponse RPC method. It returns the liquidity pool corresponding to the requested pool_id.
 */
export interface V1Beta1QueryLiquidityPoolResponse {
    pool?: V1Beta1Pool;
}
/**
 * the response type for the QueryLiquidityPoolsResponse RPC method. This includes list of all liquidity pools currently existed and paging results containing next_key and total count.
 */
export interface V1Beta1QueryLiquidityPoolsResponse {
    pools?: V1Beta1Pool[];
    /** pagination defines the pagination in the response. not working on this version. */
    pagination?: V1Beta1PageResponse;
}
/**
 * the response type for the QueryParamsResponse RPC method. This includes current parameter of the liquidity module.
 */
export interface V1Beta1QueryParamsResponse {
    /** params holds all the parameters of this module. */
    params?: V1Beta1Params;
}
export interface V1Beta1QueryPoolBatchDepositMsgResponse {
    deposit?: V1Beta1DepositMsgState;
}
/**
 * the response type for the QueryPoolBatchDeposit RPC method. This includes a list of all currently existing deposit messages of the batch and paging results containing next_key and total count.
 */
export interface V1Beta1QueryPoolBatchDepositMsgsResponse {
    deposits?: V1Beta1DepositMsgState[];
    /** pagination defines the pagination in the response. not working on this version. */
    pagination?: V1Beta1PageResponse;
}
export interface V1Beta1QueryPoolBatchSwapMsgResponse {
    swap?: V1Beta1SwapMsgState;
}
/**
 * the response type for the QueryPoolBatchSwapMsgs RPC method. This includes list of all currently existing swap messages of the batch and paging results containing next_key and total count.
 */
export interface V1Beta1QueryPoolBatchSwapMsgsResponse {
    swaps?: V1Beta1SwapMsgState[];
    /** pagination defines the pagination in the response. not working on this version. */
    pagination?: V1Beta1PageResponse;
}
export interface V1Beta1QueryPoolBatchWithdrawMsgResponse {
    withdraw?: V1Beta1WithdrawMsgState;
}
/**
 * the response type for the QueryPoolBatchWithdraw RPC method. This includes a list of all currently existing withdraw messages of the batch and paging results containing next_key and total count.
 */
export interface V1Beta1QueryPoolBatchWithdrawMsgsResponse {
    withdraws?: V1Beta1WithdrawMsgState[];
    /** pagination defines the pagination in the response. not working on this version. */
    pagination?: V1Beta1PageResponse;
}
export interface V1Beta1SwapMsgState {
    /**
     * @format int64
     * @example 1000
     */
    msgHeight?: string;
    /**
     * @format uint64
     * @example 1
     */
    msgIndex?: string;
    /** @example true */
    executed?: boolean;
    /** @example true */
    succeeded?: boolean;
    /** @example true */
    toBeDeleted?: boolean;
    /**
     * @format int64
     * @example 1000
     */
    orderExpiryHeight?: string;
    /**
     * Coin defines a token with a denomination and an amount.
     *
     * NOTE: The amount field is an Int which implements the custom method
     * signatures required by gogoproto.
     * @format sdk.Coin
     * @example {"denom":"denomX","amount":"600000"}
     */
    exchangedOfferCoin?: V1Beta1Coin;
    /**
     * Coin defines a token with a denomination and an amount.
     *
     * NOTE: The amount field is an Int which implements the custom method
     * signatures required by gogoproto.
     * @format sdk.Coin
     * @example {"denom":"denomX","amount":"400000"}
     */
    remainingOfferCoin?: V1Beta1Coin;
    /**
     * Coin defines a token with a denomination and an amount.
     *
     * NOTE: The amount field is an Int which implements the custom method
     * signatures required by gogoproto.
     * @format sdk.Coin
     * @example {"denom":"denomX","amount":"5000"}
     */
    reservedOfferCoinFee?: V1Beta1Coin;
    /**
     * `MsgSwapWithinBatch` defines an sdk.Msg type that submits a swap offer request to the liquidity pool batch
     * Submit swap offer to the liquidity pool batch with the specified the `pool_id`, `swap_type_id`,
     * `demand_coin_denom` with the coin and the price you're offering
     * The `offer_coin_fee` must be half of the offer coin amount * current `params.swap_fee_rate` for reservation to pay fees
     * This request is added to the pool and executed at the end of the batch (`endblock`)
     * You must submit the request using the same fields as the pool
     * Only the default `swap_type_id`1 is supported
     * The detailed swap algorithm is shown here.
     *
     * See: https://github.com/tendermint/liquidity/tree/develop/doc
     * https://github.com/tendermint/liquidity/blob/develop/x/liquidity/spec/04_messages.md
     */
    msg?: V1Beta1MsgSwapWithinBatch;
}
export interface V1Beta1WithdrawMsgState {
    /**
     * @format int64
     * @example 1000
     */
    msgHeight?: string;
    /**
     * @format uint64
     * @example 1
     */
    msgIndex?: string;
    /** @example true */
    executed?: boolean;
    /** @example true */
    succeeded?: boolean;
    /** @example true */
    toBeDeleted?: boolean;
    /**
     * `MsgWithdrawWithinBatch` defines an `sdk.Msg` type that submits a withdraw request to the liquidity pool batch
     * Withdraw submit to the batch from the Liquidity pool with the specified `pool_id`, `pool_coin` of the pool
     * this requests are stacked in the batch of the liquidity pool, not immediately processed and
     * processed in the `endblock` at once with other requests.
     *
     * See: https://github.com/tendermint/liquidity/blob/develop/x/liquidity/spec/04_messages.md
     */
    msg?: V1Beta1MsgWithdrawWithinBatch;
}
export declare type QueryParamsType = Record<string | number, any>;
export declare type ResponseFormat = keyof Omit<Body, "body" | "bodyUsed">;
export interface FullRequestParams extends Omit<RequestInit, "body"> {
    /** set parameter to `true` for call `securityWorker` for this request */
    secure?: boolean;
    /** request path */
    path: string;
    /** content type of request body */
    type?: ContentType;
    /** query params */
    query?: QueryParamsType;
    /** format of response (i.e. response.json() -> format: "json") */
    format?: keyof Omit<Body, "body" | "bodyUsed">;
    /** request body */
    body?: unknown;
    /** base url */
    baseUrl?: string;
    /** request cancellation token */
    cancelToken?: CancelToken;
}
export declare type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;
export interface ApiConfig<SecurityDataType = unknown> {
    baseUrl?: string;
    baseApiParams?: Omit<RequestParams, "baseUrl" | "cancelToken" | "signal">;
    securityWorker?: (securityData: SecurityDataType) => RequestParams | void;
}
export interface HttpResponse<D extends unknown, E extends unknown = unknown> extends Response {
    data: D;
    error: E;
}
declare type CancelToken = Symbol | string | number;
export declare enum ContentType {
    Json = "application/json",
    FormData = "multipart/form-data",
    UrlEncoded = "application/x-www-form-urlencoded"
}
export declare class HttpClient<SecurityDataType = unknown> {
    baseUrl: string;
    private securityData;
    private securityWorker;
    private abortControllers;
    private baseApiParams;
    constructor(apiConfig?: ApiConfig<SecurityDataType>);
    setSecurityData: (data: SecurityDataType) => void;
    private addQueryParam;
    protected toQueryString(rawQuery?: QueryParamsType): string;
    protected addQueryParams(rawQuery?: QueryParamsType): string;
    private contentFormatters;
    private mergeRequestParams;
    private createAbortSignal;
    abortRequest: (cancelToken: CancelToken) => void;
    request: <T = any, E = any>({ body, secure, path, type, query, format, baseUrl, cancelToken, ...params }: FullRequestParams) => Promise<HttpResponse<T, E>>;
}
/**
 * @title tendermint/liquidity/v1beta1/genesis.proto
 * @version version not set
 */
export declare class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
    /**
     * @description It returns all parameters of the liquidity module.
     *
     * @tags Query
     * @name QueryParams
     * @summary Get all parameters of the liquidity module.
     * @request GET:/tendermint/liquidity/v1beta1/params
     */
    queryParams: (params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryParamsResponse, RpcStatus>>;
    /**
     * @description It returns list of all liquidity pools with pagination result.
     *
     * @tags Query
     * @name QueryLiquidityPools
     * @summary Get existing liquidity pools.
     * @request GET:/tendermint/liquidity/v1beta1/pools
     */
    queryLiquidityPools: (query?: {
        "pagination.key"?: string;
        "pagination.offset"?: string;
        "pagination.limit"?: string;
        "pagination.countTotal"?: boolean;
    }, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryLiquidityPoolsResponse, RpcStatus>>;
    /**
     * @description It returns the liquidity pool corresponding to the pool_id.
     *
     * @tags Query
     * @name QueryLiquidityPool
     * @summary Get specific liquidity pool.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}
     */
    queryLiquidityPool: (poolId: string, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryLiquidityPoolResponse, RpcStatus>>;
    /**
     * @description It returns the current batch of the pool corresponding to the pool_id.
     *
     * @tags Query
     * @name QueryLiquidityPoolBatch
     * @summary Get the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch
     */
    queryLiquidityPoolBatch: (poolId: string, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryLiquidityPoolBatchResponse, RpcStatus>>;
    /**
     * @description It returns list of all deposit messages in the current batch of the pool with pagination result.
     *
     * @tags Query
     * @name QueryPoolBatchDepositMsgs
     * @summary Get all deposit messages in the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/deposits
     */
    queryPoolBatchDepositMsgs: (poolId: string, query?: {
        "pagination.key"?: string;
        "pagination.offset"?: string;
        "pagination.limit"?: string;
        "pagination.countTotal"?: boolean;
    }, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryPoolBatchDepositMsgsResponse, RpcStatus>>;
    /**
     * @description It returns the deposit message corresponding to the msg_index in the pool's current batch.
     *
     * @tags Query
     * @name QueryPoolBatchDepositMsg
     * @summary Get specific deposit message in the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/deposits/{msgIndex}
     */
    queryPoolBatchDepositMsg: (poolId: string, msgIndex: string, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryPoolBatchDepositMsgResponse, RpcStatus>>;
    /**
     * @description It returns list of all swap messages in the current batch of the pool with pagination result.
     *
     * @tags Query
     * @name QueryPoolBatchSwapMsgs
     * @summary Get all swap messages in the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/swaps
     */
    queryPoolBatchSwapMsgs: (poolId: string, query?: {
        "pagination.key"?: string;
        "pagination.offset"?: string;
        "pagination.limit"?: string;
        "pagination.countTotal"?: boolean;
    }, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryPoolBatchSwapMsgsResponse, RpcStatus>>;
    /**
     * @description It returns the swap message corresponding to the msg_index in the pool's current batch
     *
     * @tags Query
     * @name QueryPoolBatchSwapMsg
     * @summary Get specific swap message in the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/swaps/{msgIndex}
     */
    queryPoolBatchSwapMsg: (poolId: string, msgIndex: string, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryPoolBatchSwapMsgResponse, RpcStatus>>;
    /**
     * @description It returns list of all withdraw messages in the current batch of the pool with pagination result.
     *
     * @tags Query
     * @name QueryPoolBatchWithdrawMsgs
     * @summary Get all withdraw messages in the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/withdraws
     */
    queryPoolBatchWithdrawMsgs: (poolId: string, query?: {
        "pagination.key"?: string;
        "pagination.offset"?: string;
        "pagination.limit"?: string;
        "pagination.countTotal"?: boolean;
    }, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryPoolBatchWithdrawMsgsResponse, RpcStatus>>;
    /**
     * @description It returns the withdraw message corresponding to the msg_index in the pool's current batch.
     *
     * @tags Query
     * @name QueryPoolBatchWithdrawMsg
     * @summary Get specific withdraw message in the pool's current batch.
     * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/withdraws/{msgIndex}
     */
    queryPoolBatchWithdrawMsg: (poolId: string, msgIndex: string, params?: RequestParams) => Promise<HttpResponse<V1Beta1QueryPoolBatchWithdrawMsgResponse, RpcStatus>>;
}
export {};
