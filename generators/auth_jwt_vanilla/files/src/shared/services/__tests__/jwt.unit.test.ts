import { describe, expect, it } from 'vitest';
import { signRefreshToken, signToken, verifyToken } from '../jwt';

describe('jwt', () => {
  it('signToken and verifyToken round-trip', () => {
    const payload = { id: 'user-1', email: 'user@example.com' };
    const token = signToken(payload);
    const decoded = verifyToken<{ id: string; email: string }>(token);
    expect(decoded.id).toBe(payload.id);
    expect(decoded.email).toBe(payload.email);
  });

  it('signRefreshToken and verifyToken round-trip', () => {
    const payload = { id: 'user-2' };
    const token = signRefreshToken(payload);
    const decoded = verifyToken<{ id: string }>(token);
    expect(decoded.id).toBe(payload.id);
  });

  it('verifyToken throws on invalid token', () => {
    expect(() => verifyToken('invalid.token.here')).toThrow();
  });
});
