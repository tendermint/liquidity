import { txClient, queryClient, MissingWalletError } from './module'
// @ts-ignore
import { SpVuexError } from '@starport/vuex'

import { PoolRecord } from "./module/types/tendermint/liquidity/v1beta1/genesis"
import { PoolType } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { Params } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { Pool } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { PoolMetadata } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { PoolBatch } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { DepositMsgState } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { WithdrawMsgState } from "./module/types/tendermint/liquidity/v1beta1/liquidity"
import { SwapMsgState } from "./module/types/tendermint/liquidity/v1beta1/liquidity"


export { PoolRecord, PoolType, Params, Pool, PoolMetadata, PoolBatch, DepositMsgState, WithdrawMsgState, SwapMsgState };

async function initTxClient(vuexGetters) {
	return await txClient(vuexGetters['common/wallet/signer'], {
		addr: vuexGetters['common/env/apiTendermint']
	})
}

async function initQueryClient(vuexGetters) {
	return await queryClient({
		addr: vuexGetters['common/env/apiCosmos']
	})
}

function mergeResults(value, next_values) {
	for (let prop of Object.keys(next_values)) {
		if (Array.isArray(next_values[prop])) {
			value[prop]=[...value[prop], ...next_values[prop]]
		}else{
			value[prop]=next_values[prop]
		}
	}
	return value
}

function getStructure(template) {
	let structure = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field: any = {}
		field.name = key
		field.type = typeof value
		structure.fields.push(field)
	}
	return structure
}

const getDefaultState = () => {
	return {
				LiquidityPools: {},
				LiquidityPool: {},
				LiquidityPoolBatch: {},
				PoolBatchSwapMsgs: {},
				PoolBatchSwapMsg: {},
				PoolBatchDepositMsgs: {},
				PoolBatchDepositMsg: {},
				PoolBatchWithdrawMsgs: {},
				PoolBatchWithdrawMsg: {},
				Params: {},
				
				_Structure: {
						PoolRecord: getStructure(PoolRecord.fromPartial({})),
						PoolType: getStructure(PoolType.fromPartial({})),
						Params: getStructure(Params.fromPartial({})),
						Pool: getStructure(Pool.fromPartial({})),
						PoolMetadata: getStructure(PoolMetadata.fromPartial({})),
						PoolBatch: getStructure(PoolBatch.fromPartial({})),
						DepositMsgState: getStructure(DepositMsgState.fromPartial({})),
						WithdrawMsgState: getStructure(WithdrawMsgState.fromPartial({})),
						SwapMsgState: getStructure(SwapMsgState.fromPartial({})),
						
		},
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(subscription)
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(subscription)
		}
	},
	getters: {
				getLiquidityPools: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.LiquidityPools[JSON.stringify(params)] ?? {}
		},
				getLiquidityPool: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.LiquidityPool[JSON.stringify(params)] ?? {}
		},
				getLiquidityPoolBatch: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.LiquidityPoolBatch[JSON.stringify(params)] ?? {}
		},
				getPoolBatchSwapMsgs: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.PoolBatchSwapMsgs[JSON.stringify(params)] ?? {}
		},
				getPoolBatchSwapMsg: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.PoolBatchSwapMsg[JSON.stringify(params)] ?? {}
		},
				getPoolBatchDepositMsgs: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.PoolBatchDepositMsgs[JSON.stringify(params)] ?? {}
		},
				getPoolBatchDepositMsg: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.PoolBatchDepositMsg[JSON.stringify(params)] ?? {}
		},
				getPoolBatchWithdrawMsgs: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.PoolBatchWithdrawMsgs[JSON.stringify(params)] ?? {}
		},
				getPoolBatchWithdrawMsg: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.PoolBatchWithdrawMsg[JSON.stringify(params)] ?? {}
		},
				getParams: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Params[JSON.stringify(params)] ?? {}
		},
				
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('Vuex module: tendermint.liquidity.v1beta1 initialized!')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach(async (subscription) => {
				try {
					await dispatch(subscription.action, subscription.payload)
				}catch(e) {
					throw new SpVuexError('Subscriptions: ' + e.message)
				}
			})
		},
		
		
		
		 		
		
		
		async QueryLiquidityPools({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryLiquidityPools(query)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await queryClient.queryLiquidityPools({...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'LiquidityPools', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryLiquidityPools', payload: { options: { all }, params: {...key},query }})
				return getters['getLiquidityPools']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryLiquidityPools', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryLiquidityPool({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryLiquidityPool( key.pool_id)).data
				
					
				commit('QUERY', { query: 'LiquidityPool', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryLiquidityPool', payload: { options: { all }, params: {...key},query }})
				return getters['getLiquidityPool']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryLiquidityPool', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryLiquidityPoolBatch({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryLiquidityPoolBatch( key.pool_id)).data
				
					
				commit('QUERY', { query: 'LiquidityPoolBatch', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryLiquidityPoolBatch', payload: { options: { all }, params: {...key},query }})
				return getters['getLiquidityPoolBatch']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryLiquidityPoolBatch', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryPoolBatchSwapMsgs({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryPoolBatchSwapMsgs( key.pool_id, query)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await queryClient.queryPoolBatchSwapMsgs( key.pool_id, {...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'PoolBatchSwapMsgs', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPoolBatchSwapMsgs', payload: { options: { all }, params: {...key},query }})
				return getters['getPoolBatchSwapMsgs']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryPoolBatchSwapMsgs', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryPoolBatchSwapMsg({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryPoolBatchSwapMsg( key.pool_id,  key.msg_index)).data
				
					
				commit('QUERY', { query: 'PoolBatchSwapMsg', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPoolBatchSwapMsg', payload: { options: { all }, params: {...key},query }})
				return getters['getPoolBatchSwapMsg']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryPoolBatchSwapMsg', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryPoolBatchDepositMsgs({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryPoolBatchDepositMsgs( key.pool_id, query)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await queryClient.queryPoolBatchDepositMsgs( key.pool_id, {...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'PoolBatchDepositMsgs', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPoolBatchDepositMsgs', payload: { options: { all }, params: {...key},query }})
				return getters['getPoolBatchDepositMsgs']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryPoolBatchDepositMsgs', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryPoolBatchDepositMsg({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryPoolBatchDepositMsg( key.pool_id,  key.msg_index)).data
				
					
				commit('QUERY', { query: 'PoolBatchDepositMsg', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPoolBatchDepositMsg', payload: { options: { all }, params: {...key},query }})
				return getters['getPoolBatchDepositMsg']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryPoolBatchDepositMsg', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryPoolBatchWithdrawMsgs({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryPoolBatchWithdrawMsgs( key.pool_id, query)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.nextKey!=null) {
					let next_values=(await queryClient.queryPoolBatchWithdrawMsgs( key.pool_id, {...query, 'pagination.key':(<any> value).pagination.nextKey})).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'PoolBatchWithdrawMsgs', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPoolBatchWithdrawMsgs', payload: { options: { all }, params: {...key},query }})
				return getters['getPoolBatchWithdrawMsgs']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryPoolBatchWithdrawMsgs', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryPoolBatchWithdrawMsg({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryPoolBatchWithdrawMsg( key.pool_id,  key.msg_index)).data
				
					
				commit('QUERY', { query: 'PoolBatchWithdrawMsg', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryPoolBatchWithdrawMsg', payload: { options: { all }, params: {...key},query }})
				return getters['getPoolBatchWithdrawMsg']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryPoolBatchWithdrawMsg', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryParams({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params: {...key}, query=null }) {
			try {
				const queryClient=await initQueryClient(rootGetters)
				let value= (await queryClient.queryParams()).data
				
					
				commit('QUERY', { query: 'Params', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryParams', payload: { options: { all }, params: {...key},query }})
				return getters['getParams']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new SpVuexError('QueryClient:QueryParams', 'API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		async sendMsgDepositWithinBatch({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgDepositWithinBatch(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgDepositWithinBatch:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgDepositWithinBatch:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		async sendMsgWithdrawWithinBatch({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgWithdrawWithinBatch(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgWithdrawWithinBatch:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgWithdrawWithinBatch:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		async sendMsgSwapWithinBatch({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgSwapWithinBatch(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgSwapWithinBatch:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSwapWithinBatch:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		async sendMsgCreatePool({ rootGetters }, { value, fee = [], memo = '' }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgCreatePool(value)
				const result = await txClient.signAndBroadcast([msg], {fee: { amount: fee, 
	gas: "200000" }, memo})
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgCreatePool:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreatePool:Send', 'Could not broadcast Tx: '+ e.message)
				}
			}
		},
		
		async MsgDepositWithinBatch({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgDepositWithinBatch(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgDepositWithinBatch:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgDepositWithinBatch:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		async MsgWithdrawWithinBatch({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgWithdrawWithinBatch(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgWithdrawWithinBatch:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgWithdrawWithinBatch:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		async MsgSwapWithinBatch({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgSwapWithinBatch(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgSwapWithinBatch:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgSwapWithinBatch:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		async MsgCreatePool({ rootGetters }, { value }) {
			try {
				const txClient=await initTxClient(rootGetters)
				const msg = await txClient.msgCreatePool(value)
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new SpVuexError('TxClient:MsgCreatePool:Init', 'Could not initialize signing client. Wallet is required.')
				}else{
					throw new SpVuexError('TxClient:MsgCreatePool:Create', 'Could not create message: ' + e.message)
					
				}
			}
		},
		
	}
}
