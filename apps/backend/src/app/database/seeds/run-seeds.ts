import 'reflect-metadata';
import { config } from 'dotenv';
import { resolve } from 'path';
import dataSource from '../config/typeorm.config';
import { seedAdmins } from './admin.seed';

// Load environment variables
config({ path: resolve(__dirname, '../../../../../.env') });

async function runSeeds() {
  try {
    console.log('🌱 Starting database seeding...');

    // Initialize data source
    await dataSource.initialize();
    console.log('✅ Database connection established');

    // Run seeds
    await seedAdmins(dataSource);

    console.log('🎉 Seeding completed successfully!');
  } catch (error) {
    console.error('❌ Error during seeding:', error);
    process.exit(1);
  } finally {
    await dataSource.destroy();
  }
}

runSeeds();
