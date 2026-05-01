import express from 'express';
import request from 'supertest';
import { describe, expect, it } from 'vitest';
import { AppError } from '../../errors/app.error';
import { errorMiddleware } from '../error.middleware';

const app = express();
app.get('/app-error', (_req, _res, next) => {
  next(new AppError('Not found', 404, 'NOT_FOUND'));
});
app.get('/generic-error', (_req, _res, next) => {
  next(new Error('unexpected failure'));
});
app.use(errorMiddleware);

describe('errorMiddleware (feature)', () => {
  it('returns AppError status and body', async () => {
    const res = await request(app).get('/app-error');
    expect(res.status).toBe(404);
    expect(res.body).toEqual({ error: 'Not found', code: 'NOT_FOUND' });
  });

  it('returns 500 for generic Error', async () => {
    const res = await request(app).get('/generic-error');
    expect(res.status).toBe(500);
    expect(res.body.error).toBe('unexpected failure');
  });
});
