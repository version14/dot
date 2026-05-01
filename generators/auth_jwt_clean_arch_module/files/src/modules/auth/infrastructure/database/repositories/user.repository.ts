import { eq } from 'drizzle-orm';
import { db } from '../../../../../db';
import { users } from '../../../../../db/schema';
import type { UserEntity } from '../../../domain/entities/user.entity';
import type { IUserRepository } from '../../../domain/interfaces/user.repository.interface';

export class UserRepository implements IUserRepository {
  async findByEmail(email: string): Promise<UserEntity | null> {
    const result = await db.select().from(users).where(eq(users.email, email)).limit(1);
    return result[0] ?? null;
  }

  async create(data: { email: string; passwordHash: string }): Promise<UserEntity> {
    const result = await db.insert(users).values(data).returning();
    const created = result[0];
    if (!created) {
      throw new Error('Failed to create user');
    }
    return created;
  }
}
