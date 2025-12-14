import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  Index,
} from 'typeorm';
import { CommandType } from '@phd/shared';

@Entity('event_logs')
export class EventLog {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ type: 'bigint', name: 'command_id', transformer: {
    to: (value: number) => value,
    from: (value: string) => parseInt(value, 10)
  } })
  @Index()
  commandId: number;

  @Column({ type: 'bigint', transformer: {
    to: (value: number) => value,
    from: (value: string) => parseInt(value, 10)
  } })
  timestamp: number;

  @Column({
    type: 'enum',
    enum: CommandType,
    name: 'command_type',
  })
  commandType: CommandType;

  @Column({ type: 'int', nullable: true, name: 'block_number' })
  blockNumber?: number;

  @Column({ type: 'varchar', nullable: true, name: 'transaction_hash' })
  transactionHash?: string;

  @Column({ type: 'text', nullable: true })
  data?: string;

  @Column({ type: 'varchar', nullable: true, name: 'triggered_by' })
  triggeredBy?: string;

  @Column({ type: 'jsonb', nullable: true })
  metadata: Record<string, any>;

  @Column({ name: 'event_type' })
  @Index()
  eventType: string;

  @CreateDateColumn({ name: 'created_at' })
  createdAt: Date;
}
