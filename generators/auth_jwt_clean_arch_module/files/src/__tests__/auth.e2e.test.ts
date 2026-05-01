import request from 'supertest';
import { describe, expect, it } from 'vitest';
import app from '../app';

describe('auth routes (E2E)', () => {
  it('GET /health returns status ok', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });

  it('GET /auth/me without token returns 401', async () => {
    const res = await request(app).get('/auth/me');
    expect(res.status).toBe(401);
  });
});
