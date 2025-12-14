import { MigrationInterface, QueryRunner } from 'typeorm';

export class InitialSchema1702000000000 implements MigrationInterface {
  name = 'InitialSchema1702000000000';

  public async up(queryRunner: QueryRunner): Promise<void> {
    // Enable uuid extension
    await queryRunner.query(
      `CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`
    );

    // Create command_type enum
    await queryRunner.query(
      `CREATE TYPE "command_type_enum" AS ENUM ('0', '1')`
    );

    // Create admins table
    await queryRunner.query(
      `CREATE TABLE "admins" (
        "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
        "username" character varying NOT NULL,
        "email" character varying NOT NULL,
        "password_hash" character varying NOT NULL,
        "is_active" boolean NOT NULL DEFAULT true,
        "created_at" TIMESTAMP NOT NULL DEFAULT now(),
        "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
        CONSTRAINT "UQ_admins_username" UNIQUE ("username"),
        CONSTRAINT "UQ_admins_email" UNIQUE ("email"),
        CONSTRAINT "PK_admins" PRIMARY KEY ("id")
      )`
    );

    // Create scripts table
    await queryRunner.query(
      `CREATE TABLE "scripts" (
        "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
        "name" character varying NOT NULL,
        "description" text,
        "json_data" jsonb NOT NULL,
        "command_type" "command_type_enum" NOT NULL,
        "created_by" uuid NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT now(),
        "updated_at" TIMESTAMP NOT NULL DEFAULT now(),
        CONSTRAINT "PK_scripts" PRIMARY KEY ("id")
      )`
    );

    // Create event_logs table
    await queryRunner.query(
      `CREATE TABLE "event_logs" (
        "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
        "command_id" bigint NOT NULL,
        "timestamp" bigint NOT NULL,
        "command_type" "command_type_enum" NOT NULL,
        "block_number" integer,
        "transaction_hash" character varying,
        "data" text,
        "triggered_by" character varying,
        "metadata" jsonb,
        "event_type" character varying NOT NULL,
        "created_at" TIMESTAMP NOT NULL DEFAULT now(),
        CONSTRAINT "PK_event_logs" PRIMARY KEY ("id")
      )`
    );

    // Create indexes
    await queryRunner.query(
      `CREATE INDEX "IDX_event_logs_command_id" ON "event_logs" ("command_id")`
    );

    await queryRunner.query(
      `CREATE INDEX "IDX_event_logs_event_type" ON "event_logs" ("event_type")`
    );

    await queryRunner.query(
      `CREATE INDEX "IDX_scripts_created_by" ON "scripts" ("created_by")`
    );

    // Add foreign key constraint
    await queryRunner.query(
      `ALTER TABLE "scripts"
       ADD CONSTRAINT "FK_scripts_created_by"
       FOREIGN KEY ("created_by")
       REFERENCES "admins"("id")
       ON DELETE CASCADE
       ON UPDATE NO ACTION`
    );
  }

  public async down(queryRunner: QueryRunner): Promise<void> {
    // Drop foreign key
    await queryRunner.query(
      `ALTER TABLE "scripts" DROP CONSTRAINT "FK_scripts_created_by"`
    );

    // Drop indexes
    await queryRunner.query(`DROP INDEX "IDX_scripts_created_by"`);
    await queryRunner.query(`DROP INDEX "IDX_event_logs_event_type"`);
    await queryRunner.query(`DROP INDEX "IDX_event_logs_command_id"`);

    // Drop tables
    await queryRunner.query(`DROP TABLE "event_logs"`);
    await queryRunner.query(`DROP TABLE "scripts"`);
    await queryRunner.query(`DROP TABLE "admins"`);

    // Drop enum
    await queryRunner.query(`DROP TYPE "command_type_enum"`);
  }
}
