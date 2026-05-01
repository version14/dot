import { afterAll, beforeEach, describe, expect, it } from 'vitest';
import { db } from '../../../db';
import { refreshTokens, users } from '../users.table';

afterAll(async () => {
  await db.delete(refreshTokens);
  await db.delete(users);
});

beforeEach(async () => {
  await db.delete(refreshTokens);
  await db.delete(users);
});

describe('users schema (DB)', () => {
  it('inserts and retrieves a user', async () => {
    const result = await db
      .insert(users)
      .values({ email: 'schema-test@example.com', passwordHash: 'hashed' })
      .returning();
    const user = result[0];
    expect(user).toBeDefined();
    expect(user?.email).toBe('schema-test@example.com');
    expect(user?.id).toBeDefined();
  });

  it('enforces unique email constraint', async () => {
    await db.insert(users).values({ email: 'unique@example.com', passwordHash: 'h' });
    await expect(
      db.insert(users).values({ email: 'unique@example.com', passwordHash: 'h' }),
    ).rejects.toThrow();
  });

  it('inserts a refresh token linked to a user', async () => {
    const [user] = await db
      .insert(users)
      .values({ email: 'token-test@example.com', passwordHash: 'h' })
      .returning();
    if (!user) throw new Error('user insert failed');
    const [token] = await db
      .insert(refreshTokens)
      .values({
        token: 'test-token-123',
        userId: user.id,
        expiresAt: new Date(Date.now() + 86400000),
      })
      .returning();
    expect(token?.userId).toBe(user.id);
  });

  it('cascade deletes refresh tokens when user is deleted', async () => {
    const [user] = await db
      .insert(users)
      .values({ email: 'cascade-test@example.com', passwordHash: 'h' })
      .returning();
    if (!user) throw new Error('user insert failed');
    await db.insert(refreshTokens).values({
      token: 'cascade-token', userId: user.id, expiresAt: new Date(Date.now() + 86400000)
    });
    await db.delete(users);
    const remaining = await db.select().from(refreshTokens);
    expect(remaining).toHaveLength(0);
  });
});
