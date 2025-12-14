import { Injectable, Logger } from '@nestjs/common';
import { BlockchainService } from '../blockchain/blockchain.service';
import {
  Command,
  CommandPayload,
  CommandResponse,
  CommandEvent,
} from '@phd/shared';

@Injectable()
export class CommandsService {
  private readonly logger = new Logger(CommandsService.name);
  private commandHistory: CommandEvent[] = [];

  constructor(private blockchainService: BlockchainService) {
    this.initializeEventListener();
  }

  private initializeEventListener() {
    this.blockchainService
      .startCommandEventPolling((event: CommandEvent) => {
        this.logger.log(`New command received: ${JSON.stringify(event)}`);
        this.commandHistory.push(event);

        // Keep only last 100 events
        if (this.commandHistory.length > 100) {
          this.commandHistory.shift();
        }
      })
      .catch((error) => {
        this.logger.error('Failed to initialize event listener', error);
      });
  }

  async triggerCommand(payload: CommandPayload): Promise<CommandResponse> {
    try {
      const txHash = await this.blockchainService.triggerCommand(
        payload.commandType,
        payload.data
      );

      const command = await this.blockchainService.getCurrentCommand();

      return {
        success: true,
        message: `Command triggered successfully. Transaction: ${txHash}`,
        command,
      };
    } catch (error) {
      this.logger.error('Failed to trigger command', error);
      return {
        success: false,
        message: error.message,
      };
    }
  }

  async getCurrentCommand(): Promise<Command> {
    return this.blockchainService.getCurrentCommand();
  }

  async getLatestCommandId(): Promise<number> {
    return this.blockchainService.getLatestCommandId();
  }

  getCommandHistory(): CommandEvent[] {
    return this.commandHistory;
  }

  async getBlockchainInfo() {
    const info = this.blockchainService.getProviderInfo();
    const adminAddress = await this.blockchainService.getAdminAddress();

    return {
      ...info,
      adminAddress,
      historyCount: this.commandHistory.length,
    };
  }
}
