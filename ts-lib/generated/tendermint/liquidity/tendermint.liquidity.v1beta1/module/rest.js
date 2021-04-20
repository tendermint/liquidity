/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */
export var ContentType;
(function (ContentType) {
    ContentType["Json"] = "application/json";
    ContentType["FormData"] = "multipart/form-data";
    ContentType["UrlEncoded"] = "application/x-www-form-urlencoded";
})(ContentType || (ContentType = {}));
export class HttpClient {
    constructor(apiConfig = {}) {
        this.baseUrl = "";
        this.securityData = null;
        this.securityWorker = null;
        this.abortControllers = new Map();
        this.baseApiParams = {
            credentials: "same-origin",
            headers: {},
            redirect: "follow",
            referrerPolicy: "no-referrer",
        };
        this.setSecurityData = (data) => {
            this.securityData = data;
        };
        this.contentFormatters = {
            [ContentType.Json]: (input) => input !== null && (typeof input === "object" || typeof input === "string") ? JSON.stringify(input) : input,
            [ContentType.FormData]: (input) => Object.keys(input || {}).reduce((data, key) => {
                data.append(key, input[key]);
                return data;
            }, new FormData()),
            [ContentType.UrlEncoded]: (input) => this.toQueryString(input),
        };
        this.createAbortSignal = (cancelToken) => {
            if (this.abortControllers.has(cancelToken)) {
                const abortController = this.abortControllers.get(cancelToken);
                if (abortController) {
                    return abortController.signal;
                }
                return void 0;
            }
            const abortController = new AbortController();
            this.abortControllers.set(cancelToken, abortController);
            return abortController.signal;
        };
        this.abortRequest = (cancelToken) => {
            const abortController = this.abortControllers.get(cancelToken);
            if (abortController) {
                abortController.abort();
                this.abortControllers.delete(cancelToken);
            }
        };
        this.request = ({ body, secure, path, type, query, format = "json", baseUrl, cancelToken, ...params }) => {
            const secureParams = (secure && this.securityWorker && this.securityWorker(this.securityData)) || {};
            const requestParams = this.mergeRequestParams(params, secureParams);
            const queryString = query && this.toQueryString(query);
            const payloadFormatter = this.contentFormatters[type || ContentType.Json];
            return fetch(`${baseUrl || this.baseUrl || ""}${path}${queryString ? `?${queryString}` : ""}`, {
                ...requestParams,
                headers: {
                    ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
                    ...(requestParams.headers || {}),
                },
                signal: cancelToken ? this.createAbortSignal(cancelToken) : void 0,
                body: typeof body === "undefined" || body === null ? null : payloadFormatter(body),
            }).then(async (response) => {
                const r = response;
                r.data = null;
                r.error = null;
                const data = await response[format]()
                    .then((data) => {
                    if (r.ok) {
                        r.data = data;
                    }
                    else {
                        r.error = data;
                    }
                    return r;
                })
                    .catch((e) => {
                    r.error = e;
                    return r;
                });
                if (cancelToken) {
                    this.abortControllers.delete(cancelToken);
                }
                if (!response.ok)
                    throw data;
                return data;
            });
        };
        Object.assign(this, apiConfig);
    }
    addQueryParam(query, key) {
        const value = query[key];
        return (encodeURIComponent(key) +
            "=" +
            encodeURIComponent(Array.isArray(value) ? value.join(",") : typeof value === "number" ? value : `${value}`));
    }
    toQueryString(rawQuery) {
        const query = rawQuery || {};
        const keys = Object.keys(query).filter((key) => "undefined" !== typeof query[key]);
        return keys
            .map((key) => typeof query[key] === "object" && !Array.isArray(query[key])
            ? this.toQueryString(query[key])
            : this.addQueryParam(query, key))
            .join("&");
    }
    addQueryParams(rawQuery) {
        const queryString = this.toQueryString(rawQuery);
        return queryString ? `?${queryString}` : "";
    }
    mergeRequestParams(params1, params2) {
        return {
            ...this.baseApiParams,
            ...params1,
            ...(params2 || {}),
            headers: {
                ...(this.baseApiParams.headers || {}),
                ...(params1.headers || {}),
                ...((params2 && params2.headers) || {}),
            },
        };
    }
}
/**
 * @title tendermint/liquidity/v1beta1/genesis.proto
 * @version version not set
 */
export class Api extends HttpClient {
    constructor() {
        super(...arguments);
        /**
         * @description It returns all parameters of the liquidity module.
         *
         * @tags Query
         * @name QueryParams
         * @summary Get all parameters of the liquidity module.
         * @request GET:/tendermint/liquidity/v1beta1/params
         */
        this.queryParams = (params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/params`,
            method: "GET",
            format: "json",
            ...params,
        });
        /**
         * @description It returns list of all liquidity pools with pagination result.
         *
         * @tags Query
         * @name QueryLiquidityPools
         * @summary Get existing liquidity pools.
         * @request GET:/tendermint/liquidity/v1beta1/pools
         */
        this.queryLiquidityPools = (query, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * @description It returns the liquidity pool corresponding to the pool_id.
         *
         * @tags Query
         * @name QueryLiquidityPool
         * @summary Get specific liquidity pool.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}
         */
        this.queryLiquidityPool = (poolId, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}`,
            method: "GET",
            format: "json",
            ...params,
        });
        /**
         * @description It returns the current batch of the pool corresponding to the pool_id.
         *
         * @tags Query
         * @name QueryLiquidityPoolBatch
         * @summary Get the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch
         */
        this.queryLiquidityPoolBatch = (poolId, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch`,
            method: "GET",
            format: "json",
            ...params,
        });
        /**
         * @description It returns list of all deposit messages in the current batch of the pool with pagination result.
         *
         * @tags Query
         * @name QueryPoolBatchDepositMsgs
         * @summary Get all deposit messages in the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/deposits
         */
        this.queryPoolBatchDepositMsgs = (poolId, query, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch/deposits`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * @description It returns the deposit message corresponding to the msg_index in the pool's current batch.
         *
         * @tags Query
         * @name QueryPoolBatchDepositMsg
         * @summary Get specific deposit message in the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/deposits/{msgIndex}
         */
        this.queryPoolBatchDepositMsg = (poolId, msgIndex, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch/deposits/${msgIndex}`,
            method: "GET",
            format: "json",
            ...params,
        });
        /**
         * @description It returns list of all swap messages in the current batch of the pool with pagination result.
         *
         * @tags Query
         * @name QueryPoolBatchSwapMsgs
         * @summary Get all swap messages in the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/swaps
         */
        this.queryPoolBatchSwapMsgs = (poolId, query, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch/swaps`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * @description It returns the swap message corresponding to the msg_index in the pool's current batch
         *
         * @tags Query
         * @name QueryPoolBatchSwapMsg
         * @summary Get specific swap message in the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/swaps/{msgIndex}
         */
        this.queryPoolBatchSwapMsg = (poolId, msgIndex, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch/swaps/${msgIndex}`,
            method: "GET",
            format: "json",
            ...params,
        });
        /**
         * @description It returns list of all withdraw messages in the current batch of the pool with pagination result.
         *
         * @tags Query
         * @name QueryPoolBatchWithdrawMsgs
         * @summary Get all withdraw messages in the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/withdraws
         */
        this.queryPoolBatchWithdrawMsgs = (poolId, query, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch/withdraws`,
            method: "GET",
            query: query,
            format: "json",
            ...params,
        });
        /**
         * @description It returns the withdraw message corresponding to the msg_index in the pool's current batch.
         *
         * @tags Query
         * @name QueryPoolBatchWithdrawMsg
         * @summary Get specific withdraw message in the pool's current batch.
         * @request GET:/tendermint/liquidity/v1beta1/pools/{poolId}/batch/withdraws/{msgIndex}
         */
        this.queryPoolBatchWithdrawMsg = (poolId, msgIndex, params = {}) => this.request({
            path: `/tendermint/liquidity/v1beta1/pools/${poolId}/batch/withdraws/${msgIndex}`,
            method: "GET",
            format: "json",
            ...params,
        });
    }
}
