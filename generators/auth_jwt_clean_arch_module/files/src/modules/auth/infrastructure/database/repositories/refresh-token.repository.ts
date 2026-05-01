import { eq } from 'drizzle-orm';
import { db } from '../../../../../db';
import { refreshTokens } from '../../../../../db/schema';
import type { IRefreshTokenRepository } from '../../../domain/interfaces/refresh-token.repository.interface';

export class RefreshTokenRepository implements IRefreshTokenRepository {
  async create(data: { token: string; userId: string; expiresAt: Date }): Promise<void> {
    await db.insert(refreshTokens).values(data);
  }

  async findByToken(token: string): Promise<{ id: string; userId: string } | null> {
    const result = await db
      .select({ id: refreshTokens.id, userId: refreshTokens.userId })
      .from(refreshTokens)
      .where(eq(refreshTokens.token, token))
      .limit(1);
    return result[0] ?? null;
  }

  async deleteById(id: string): Promise<void> {
    await db.delete(refreshTokens).where(eq(refreshTokens.id, id));
  }
}
