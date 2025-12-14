import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, Between, MoreThanOrEqual, LessThanOrEqual } from 'typeorm';
import { EventLog } from './entities/event-log.entity';
import { QueryEventsDto } from './dto/query-events.dto';
import { CommandEvent } from '@phd/shared';
import { BlockchainService } from '../blockchain/blockchain.service';

@Injectable()
export class EventsService implements OnModuleInit {
  private readonly logger = new Logger(EventsService.name);

  constructor(
    @InjectRepository(EventLog)
    private eventLogRepository: Repository<EventLog>,
    private blockchainService: BlockchainService
  ) {}

  async onModuleInit() {
    // Start listening to blockchain events
    await this.startEventListener();
  }

  private async startEventListener() {
    try {
      this.logger.log('Starting blockchain event listener...');

      await this.blockchainService.startCommandEventPolling(
        async (event: CommandEvent) => {
          await this.logEvent(event);
        }
      );

      this.logger.log('✅ Blockchain event listener started successfully');
    } catch (error) {
      this.logger.error('Failed to start event listener', error);
    }
  }

  async logEvent(event: CommandEvent): Promise<EventLog> {
    this.logger.log(
      `Logging blockchain event: CommandID=${event.commandId}, Type=${event.commandType}`
    );

    const eventLog = new EventLog();
    eventLog.commandId = event.commandId;
    eventLog.timestamp = event.timestamp;
    eventLog.commandType = event.commandType;
    eventLog.blockNumber = event.blockNumber;
    eventLog.transactionHash = event.transactionHash;
    eventLog.eventType = 'command_triggered';
    eventLog.triggeredBy = undefined;
    eventLog.data = undefined;
    eventLog.metadata = event;

    return this.eventLogRepository.save(eventLog);
  }

  async findAll(query: QueryEventsDto) {
    const { page = 1, limit = 20, commandType, startDate, endDate } = query;
    const skip = (page - 1) * limit;

    const where: any = {};

    if (commandType !== undefined) {
      where.commandType = commandType;
    }

    if (startDate && endDate) {
      where.createdAt = Between(new Date(startDate), new Date(endDate));
    } else if (startDate) {
      where.createdAt = MoreThanOrEqual(new Date(startDate));
    } else if (endDate) {
      where.createdAt = LessThanOrEqual(new Date(endDate));
    }

    const [events, total] = await this.eventLogRepository.findAndCount({
      where,
      skip,
      take: limit,
      order: { createdAt: 'DESC' },
    });

    return {
      data: events,
      meta: {
        total,
        page,
        limit,
        totalPages: Math.ceil(total / limit),
      },
    };
  }

  async findOne(id: string): Promise<EventLog | null> {
    return this.eventLogRepository.findOne({ where: { id } });
  }

  async getStats() {
    const totalEvents = await this.eventLogRepository.count();

    const eventsByType = await this.eventLogRepository
      .createQueryBuilder('event')
      .select('event.command_type', 'commandType')
      .addSelect('COUNT(*)', 'count')
      .groupBy('event.command_type')
      .getRawMany();

    return {
      totalEvents,
      eventsByType,
    };
  }
}
