import express from 'express';
import request from 'supertest';
import { describe, expect, it } from 'vitest';
import { limiter } from '../rate-limit.middleware';

const app = express();
app.use(limiter);
app.get('/test', (_req, res) => {
  res.json({ ok: true });
});

describe('rate limit middleware (feature)', () => {
  it('allows requests within the limit and adds rate limit headers', async () => {
    const res = await request(app).get('/test');
    expect(res.status).toBe(200);
    expect(res.body.ok).toBe(true);
    expect(res.headers.ratelimit).toBeDefined();
  });
});
