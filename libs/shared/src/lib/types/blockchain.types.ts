/**
 * Blockchain and smart contract types
 */

export interface ContractConfig {
  address: string;
  abi: any[];
  network: string;
}

export interface TransactionReceipt {
  transactionHash: string;
  blockNumber: number;
  blockHash: string;
  from: string;
  to: string;
  gasUsed: string;
  status: number;
}

export interface BlockchainEvent {
  event: string;
  args: Record<string, any>;
  blockNumber: number;
  transactionHash: string;
  address: string;
}

export interface WalletConfig {
  privateKey?: string;
  mnemonic?: string;
  address: string;
}

export interface NetworkConfig {
  name: string;
  rpcUrl: string;
  chainId: number;
  explorer?: string;
}

export const HEDERA_NETWORKS: Record<string, NetworkConfig> = {
  testnet: {
    name: 'Hedera Testnet',
    rpcUrl: 'https://testnet.hashio.io/api',
    chainId: 296,
    explorer: 'https://hashscan.io/testnet',
  },
  mainnet: {
    name: 'Hedera Mainnet',
    rpcUrl: 'https://mainnet.hashio.io/api',
    chainId: 295,
    explorer: 'https://hashscan.io/mainnet',
  },
  local: {
    name: 'Local Network',
    rpcUrl: 'http://localhost:8545',
    chainId: 31337,
  },
};
