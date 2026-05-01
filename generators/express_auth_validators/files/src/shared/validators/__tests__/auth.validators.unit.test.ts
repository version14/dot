import { describe, expect, it } from 'vitest';
import { loginSchema, refreshSchema, registerSchema } from '../auth.validators';

describe('registerSchema', () => {
  it('accepts valid email and password', () => {
    const result = registerSchema.safeParse({ email: 'user@example.com', password: 'password123' });
    expect(result.success).toBe(true);
  });

  it('rejects invalid email', () => {
    const result = registerSchema.safeParse({ email: 'not-an-email', password: 'password123' });
    expect(result.success).toBe(false);
  });

  it('rejects password shorter than 8 characters', () => {
    const result = registerSchema.safeParse({ email: 'user@example.com', password: 'short' });
    expect(result.success).toBe(false);
  });
});

describe('loginSchema', () => {
  it('accepts valid credentials', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: 'any' });
    expect(result.success).toBe(true);
  });

  it('rejects empty password', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: '' });
    expect(result.success).toBe(false);
  });
});

describe('refreshSchema', () => {
  it('accepts a non-empty token', () => {
    const result = refreshSchema.safeParse({ refreshToken: 'some.jwt.token' });
    expect(result.success).toBe(true);
  });

  it('rejects empty refreshToken', () => {
    const result = refreshSchema.safeParse({ refreshToken: '' });
    expect(result.success).toBe(false);
  });
});
