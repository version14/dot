import { describe, expect, it } from 'vitest';
import { AppError } from '../app.error';
import { ConflictError } from '../conflict.error';
import { ForbiddenError } from '../forbidden.error';
import { NotFoundError } from '../not-found.error';
import { UnauthorizedError } from '../unauthorized.error';
import { ValidationError } from '../validation.error';

describe('AppError', () => {
  it('sets message, statusCode, and code', () => {
    const err = new AppError('test', 422, 'CUSTOM');
    expect(err.message).toBe('test');
    expect(err.statusCode).toBe(422);
    expect(err.code).toBe('CUSTOM');
    expect(err).toBeInstanceOf(Error);
  });

  it('defaults statusCode to 500', () => {
    const err = new AppError('oops');
    expect(err.statusCode).toBe(500);
  });
});

describe('NotFoundError', () => {
  it('has statusCode 404 and code NOT_FOUND', () => {
    const err = new NotFoundError();
    expect(err.statusCode).toBe(404);
    expect(err.code).toBe('NOT_FOUND');
  });
});

describe('UnauthorizedError', () => {
  it('has statusCode 401 and code UNAUTHORIZED', () => {
    const err = new UnauthorizedError();
    expect(err.statusCode).toBe(401);
    expect(err.code).toBe('UNAUTHORIZED');
  });
});

describe('ForbiddenError', () => {
  it('has statusCode 403 and code FORBIDDEN', () => {
    const err = new ForbiddenError();
    expect(err.statusCode).toBe(403);
    expect(err.code).toBe('FORBIDDEN');
  });
});

describe('ConflictError', () => {
  it('has statusCode 409 and code CONFLICT', () => {
    const err = new ConflictError();
    expect(err.statusCode).toBe(409);
    expect(err.code).toBe('CONFLICT');
  });
});

describe('ValidationError', () => {
  it('has statusCode 400 and code VALIDATION_ERROR', () => {
    const err = new ValidationError();
    expect(err.statusCode).toBe(400);
    expect(err.code).toBe('VALIDATION_ERROR');
  });
});
