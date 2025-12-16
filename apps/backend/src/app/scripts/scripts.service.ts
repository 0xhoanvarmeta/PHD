import {
  Injectable,
  NotFoundException,
  ForbiddenException,
  Logger,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository, Like } from 'typeorm';
import { Script } from './entities/script.entity';
import { CreateScriptDto } from './dto/create-script.dto';
import { UpdateScriptDto } from './dto/update-script.dto';
import { QueryScriptsDto } from './dto/query-scripts.dto';
import { Admin } from '../admin/entities/admin.entity';
import { BlockchainService } from '../blockchain/blockchain.service';

@Injectable()
export class ScriptsService {
  private readonly logger = new Logger(ScriptsService.name);

  constructor(
    @InjectRepository(Script)
    private scriptRepository: Repository<Script>,
    private blockchainService: BlockchainService
  ) {}

  async create(
    createScriptDto: CreateScriptDto,
    admin: Admin
  ): Promise<Script> {
    const script = this.scriptRepository.create({
      ...createScriptDto,
      commandType: createScriptDto.jsonData.commandType,
      createdById: admin.id,
    });

    return this.scriptRepository.save(script);
  }

  async findAll(query: QueryScriptsDto) {
    const { page = 1, limit = 20, search, commandType } = query;
    const skip = (page - 1) * limit;

    const where: any = {};

    if (search) {
      where.name = Like(`%${search}%`);
    }

    if (commandType !== undefined) {
      where.commandType = commandType;
    }

    const [scripts, total] = await this.scriptRepository.findAndCount({
      where,
      relations: ['createdBy'],
      skip,
      take: limit,
      order: { createdAt: 'DESC' },
    });

    return {
      data: scripts,
      meta: {
        total,
        page,
        limit,
        totalPages: Math.ceil(total / limit),
      },
    };
  }

  async findOne(id: string): Promise<Script> {
    const script = await this.scriptRepository.findOne({
      where: { id },
      relations: ['createdBy'],
    });

    if (!script) {
      throw new NotFoundException(`Script with ID ${id} not found`);
    }

    return script;
  }

  async update(
    id: string,
    updateScriptDto: UpdateScriptDto,
    admin: Admin
  ): Promise<Script> {
    const script = await this.findOne(id);

    // Check if admin owns the script
    if (script.createdById !== admin.id) {
      throw new ForbiddenException('You can only update your own scripts');
    }

    // Update command type if jsonData is provided
    if (updateScriptDto.jsonData) {
      script.commandType = updateScriptDto.jsonData.commandType;
    }

    Object.assign(script, updateScriptDto);

    return this.scriptRepository.save(script);
  }

  async remove(id: string, admin: Admin): Promise<void> {
    const script = await this.findOne(id);

    // Check if admin owns the script
    if (script.createdById !== admin.id) {
      throw new ForbiddenException('You can only delete your own scripts');
    }

    await this.scriptRepository.remove(script);
  }

  async trigger(id: string, admin: Admin) {
    const script = await this.findOne(id);

    this.logger.log(
      `Triggering script ${script.id} (${script.name}) by admin ${admin.username}`
    );

    // Convert JSON data to string for blockchain
    const dataString = JSON.stringify(script.jsonData.data);

    // Trigger command on blockchain
    const transactionHash = await this.blockchainService.triggerCommand(
      script.commandType,
      dataString,
      script.id
    );

    this.logger.log(
      `Script ${script.id} triggered successfully. Transaction: ${transactionHash}`
    );

    return {
      success: true,
      transactionHash,
      script: {
        id: script.id,
        name: script.name,
        commandType: script.commandType,
      },
    };
  }
}
