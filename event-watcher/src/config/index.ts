import { ChainName, EVMChainName } from '@deltaswapio/deltaswap-sdk/lib/cjs/utils/consts';

export const env = {
  NODE_ENV: process.env.NODE_ENV || 'development',

  P2P_NETWORK: String(process.env.P2P_NETWORK).toLowerCase() || 'mainnet',
  LOG_DIR: process.env.LOG_DIR,
  LOG_LEVEL: process.env.LOG_LEVEL || 'info',

  DB_SOURCE: process.env.DB_SOURCE || 'local',

  JSON_WH_TXS_FILE: process.env.JSON_WH_TXS_FILE || './deltaswapTxs.json',
  JSON_GLOBAL_TXS_FILE: process.env.JSON_GLOBAL_TXS_FILE || './globalTxs.json',
  JSON_LAST_BLOCKS_FILE: process.env.JSON_LAST_BLOCKS_FILE || './lastBlocksByChain.json',

  PORT: process.env.PORT,

  MONGODB_URI: process.env.MONGODB_URI,
  MONGODB_DATABASE: process.env.MONGODB_DATABASE,

  SNS_SOURCE: process.env.SNS_SOURCE,
  AWS_SNS_REGION: process.env.AWS_SNS_REGION,
  AWS_SNS_TOPIC_ARN: process.env.AWS_SNS_TOPIC_ARN,
  AWS_SNS_SUBJECT: process.env.AWS_SNS_SUBJECT,
  AWS_ACCESS_KEY_ID: process.env.AWS_ACCESS_KEY_ID,
  AWS_SECRET_ACCESS_KEY: process.env.AWS_SECRET_ACCESS_KEY,

  CHAINS: process.env.CHAINS,

  ACALA_RPC: process.env.ACALA_RPC,
  ALGORAND_RPC: process.env.ALGORAND_RPC,
  APTOS_RPC: process.env.APTOS_RPC,
  ARBITRUM_RPC: process.env.ARBITRUM_RPC,
  AVALANCHE_RPC: process.env.AVALANCHE_RPC,
  BASE_RPC: process.env.BASE_RPC,
  BSC_RPC: process.env.BSC_RPC,
  CELO_RPC: process.env.CELO_RPC,
  ETHEREUM_RPC: process.env.ETHEREUM_RPC,
  FANTOM_RPC: process.env.FANTOM_RPC,
  INJECTIVE_RPC: process.env.INJECTIVE_RPC,
  KARURA_RPC: process.env.KARURA_RPC,
  KLAYTN_RPC: process.env.KLAYTN_RPC,
  MOONBEAM_RPC: process.env.MOONBEAM_RPC,
  NEAR_RPC: process.env.NEAR_RPC,
  OASIS_RPC: process.env.OASIS_RPC,
  OPTIMISM_RPC: process.env.OPTIMISM_RPC,
  PLANQ_RPC: process.env.PLANQ_RPC,
  POLYGON_RPC: process.env.POLYGON_RPC,
  SEI_RPC: process.env.SEI_RPC,
  SOLANA_RPC: process.env.SOLANA_RPC,
  SUI_RPC: process.env.SUI_RPC,
  TERRA_RPC: process.env.TERRA_RPC,
  TERRA2_RPC: process.env.TERRA2_RPC,
  DELTACHAIN_RPC: process.env.DELTACHAIN_RPC,
  XPLA_RPC: process.env.XPLA_RPC,

  ACALA_RPS: process.env.ACALA_RPS,
  ALGORAND_RPS: process.env.ALGORAND_RPS,
  APTOS_RPS: process.env.APTOS_RPS,
  ARBITRUM_RPS: process.env.ARBITRUM_RPS,
  AVALANCHE_RPS: process.env.AVALANCHE_RPS,
  BASE_RPS: process.env.BASE_RPS,
  BSC_RPS: process.env.BSC_RPS,
  CELO_RPS: process.env.CELO_RPS,
  ETHEREUM_RPS: process.env.ETHEREUM_RPS,
  FANTOM_RPS: process.env.FANTOM_RPS,
  INJECTIVE_RPS: process.env.INJECTIVE_RPS,
  KARURA_RPS: process.env.KARURA_RPS,
  KLAYTN_RPS: process.env.KLAYTN_RPS,
  MOONBEAM_RPS: process.env.MOONBEAM_RPS,
  NEAR_RPS: process.env.NEAR_RPS,
  OASIS_RPS: process.env.OASIS_RPS,
  OPTIMISM_RPS: process.env.OPTIMISM_RPS,
  PLANQ_RPS: process.env.PLANQ_RPS,
  POLYGON_RPS: process.env.POLYGON_RPS,
  SEI_RPS: process.env.SEI_RPS,
  SOLANA_RPS: process.env.SOLANA_RPS,
  SUI_RPS: process.env.SUI_RPS,
  TERRA_RPS: process.env.TERRA_RPS,
  TERRA2_RPS: process.env.TERRA2_RPS,
  DELTACHAIN_RPS: process.env.DELTACHAIN_RPS,
  XPLA_RPS: process.env.XPLA_RPS,
} as const;

// EVM Chains not supported
// aurora, gnosis, neon, sepolia

export const evmChains: EVMChainName[] = [
 /* 'acala',
  'arbitrum',
  'avalanche',
  'base',
  'celo',
  'ethereum',
  'fantom',
  'karura',
  'klaytn',
  'moonbeam',
  'oasis',
  'optimism',
  'polygon',*/
  'planq',
  'bsc',
];

export const supportedChains: ChainName[] = [
  ...evmChains,
 /* 'algorand',
  'aptos',
  'injective',
  'near',
  'sei',
  'solana',
  'sui',
  'terra',
  'terra2',
  'xpla',*/
];
