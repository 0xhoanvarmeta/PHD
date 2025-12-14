import { DataSource } from 'typeorm';
import * as bcrypt from 'bcrypt';
import { Admin } from '../../admin/entities/admin.entity';

export async function seedAdmins(dataSource: DataSource): Promise<void> {
  const adminRepository = dataSource.getRepository(Admin);

  const defaultAdmins = [
    {
      username: 'admin',
      email: 'admin@phd.local',
      password: 'Admin@123', // Should be changed after first login
    },
  ];

  for (const adminData of defaultAdmins) {
    const existing = await adminRepository.findOne({
      where: { username: adminData.username },
    });

    if (!existing) {
      const admin = adminRepository.create({
        username: adminData.username,
        email: adminData.email,
        passwordHash: await bcrypt.hash(adminData.password, 12),
        isActive: true,
      });
      await adminRepository.save(admin);
      console.log(`✅ Seeded admin: ${adminData.username}`);
    } else {
      console.log(`ℹ️  Admin already exists: ${adminData.username}`);
    }
  }
}
