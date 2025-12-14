import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { ethers } from 'ethers';
import {
  Command,
  CommandType,
  CommandEvent,
  HEDERA_NETWORKS,
  DeviceControlABI,
  DeviceControlABIJson,
} from '@phd/shared';

@Injectable()
export class BlockchainService implements OnModuleInit {
  private readonly logger = new Logger(BlockchainService.name);
  private provider: ethers.JsonRpcProvider;
  private contract: ethers.Contract;
  private wallet: ethers.Wallet;
  private pollingInterval?: NodeJS.Timeout;
  private lastProcessedBlock = 0;
  private iface = new ethers.Interface(DeviceControlABIJson);

  constructor(private configService: ConfigService) {}

  async onModuleInit() {
    await this.initializeBlockchain();
  }

  private async initializeBlockchain() {
    try {
      const network = this.configService.get<string>(
        'BLOCKCHAIN_NETWORK',
        'testnet'
      );
      const contractAddress =
        this.configService.get<string>('CONTRACT_ADDRESS');
      const privateKey = this.configService.get<string>('PRIVATE_KEY');

      if (!contractAddress) {
        this.logger.warn('CONTRACT_ADDRESS not configured');
        return;
      }

      const networkConfig = HEDERA_NETWORKS[network];
      if (!networkConfig) {
        throw new Error(`Invalid network: ${network}`);
      }

      this.provider = new ethers.JsonRpcProvider(networkConfig.rpcUrl);
      this.logger.log(`Connected to ${networkConfig.name}`);

      if (privateKey) {
        this.wallet = new ethers.Wallet(privateKey, this.provider);
        this.logger.log(`Wallet address: ${this.wallet.address}`);
      }

      this.contract = new ethers.Contract(
        contractAddress,
        DeviceControlABIJson,
        this.wallet || this.provider
      );

      this.logger.log(`Contract initialized at ${contractAddress}`);
    } catch (error) {
      this.logger.error('Failed to initialize blockchain', error);
      throw error;
    }
  }

  async triggerCommand(
    commandType: CommandType,
    data: string
  ): Promise<string> {
    try {
      if (!this.wallet) {
        throw new Error('Wallet not configured. Cannot send transactions.');
      }

      this.logger.log(
        `Triggering command: ${CommandType[commandType]} - ${data}`
      );

      const tx = await this.contract.Trigger(commandType, data);
      this.logger.log(`Transaction sent: ${tx.hash}`);

      const receipt = await tx.wait();
      this.logger.log(`Transaction confirmed in block ${receipt.blockNumber}`);

      return tx.hash;
    } catch (error) {
      this.logger.error('Failed to trigger command', error);
      throw error;
    }
  }

  async getCurrentCommand(): Promise<Command> {
    try {
      const [id, commandType, data, timestamp] =
        await this.contract.GetFunction();

      return {
        id: Number(id),
        commandType: Number(commandType),
        data,
        timestamp: Number(timestamp),
      };
    } catch (error) {
      this.logger.error('Failed to get current command', error);
      throw error;
    }
  }

  async getLatestCommandId(): Promise<number> {
    try {
      const commandId = await this.contract.getLatestCommandId();
      return Number(commandId);
    } catch (error) {
      this.logger.error('Failed to get latest command ID', error);
      throw error;
    }
  }

  // async listenToCommandEvents(
  //   callback: (event: CommandEvent) => void
  // ): Promise<void> {
  //   try {
  //     if (!this.contract) {
  //       this.logger.warn(
  //         'Contract not initialized. Skipping event listener setup.'
  //       );
  //       return;
  //     }

  //     this.logger.log('Starting to listen for CommandTriggered events...');

  //     this.contract.on(
  //       'CommandTriggered',
  //       (commandId, timestamp, commandType, event) => {
  //         const commandEvent: CommandEvent = {
  //           commandId: Number(commandId),
  //           timestamp: Number(timestamp),
  //           commandType: Number(commandType),
  //           blockNumber: event.log.blockNumber,
  //           transactionHash: event.log.transactionHash,
  //         };

  //         this.logger.log(
  //           `CommandTriggered event: ID=${commandEvent.commandId}, Type=${
  //             CommandType[commandEvent.commandType]
  //           }`
  //         );

  //         callback(commandEvent);
  //       }
  //     );
  //   } catch (error) {
  //     this.logger.error('Failed to listen to events', error);
  //     throw error;
  //   }
  // }

  async startCommandEventPolling(
    callback: (event: CommandEvent) => void,
    intervalMs = 5000
  ) {
    const contractAddress = this.contract.target as string;

    const eventFragment = this.iface.getEvent('CommandTriggered');
    if (!eventFragment) {
      throw new Error('Failed to get event fragment');
    }
    const topic = ethers.id(eventFragment.format());

    if (this.lastProcessedBlock === 0) {
      this.lastProcessedBlock = await this.provider.getBlockNumber();
    }

    this.pollingInterval = setInterval(async () => {
      try {
        const latestBlock = await this.provider.getBlockNumber();
        if (latestBlock <= this.lastProcessedBlock) return;

        const logs = await this.provider.getLogs({
          address: contractAddress,
          topics: [topic],
          fromBlock: this.lastProcessedBlock + 1,
          toBlock: latestBlock,
        });

        for (const log of logs) {
          const parsed = this.iface.parseLog(log);
          if (!parsed) continue;
          const commandEvent: CommandEvent = {
            commandId: Number(parsed.args.commandId),
            timestamp: Number(parsed.args.timestamp),
            commandType: Number(parsed.args.commandType),
            blockNumber: log.blockNumber,
            transactionHash: log.transactionHash,
          };

          callback(commandEvent);
        }

        this.lastProcessedBlock = latestBlock;
      } catch (err) {
        this.logger.error('Polling failed', err);
      }
    }, intervalMs);
  }

  async getAdminAddress(): Promise<string> {
    try {
      return await this.contract.admin();
    } catch (error) {
      this.logger.error('Failed to get admin address', error);
      throw error;
    }
  }

  getProviderInfo() {
    return {
      network: this.provider ? 'connected' : 'disconnected',
      contractAddress: this.contract?.target,
      walletAddress: this.wallet?.address,
    };
  }
}
