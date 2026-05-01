import type { NextFunction, Request, Response } from 'express';
import type { AuthRequest } from '../../../../shared/middlewares/auth.middleware';
import { RefreshTokenRepository } from '../../infrastructure/database/repositories/refresh-token.repository';
import { UserRepository } from '../../infrastructure/database/repositories/user.repository';
import { LoginUseCase } from '../use-cases/login.use-case';
import { LogoutUseCase } from '../use-cases/logout.use-case';
import { RefreshUseCase } from '../use-cases/refresh.use-case';
import { RegisterUseCase } from '../use-cases/register.use-case';

const userRepository = new UserRepository();
const refreshTokenRepository = new RefreshTokenRepository();
const loginUseCase = new LoginUseCase(userRepository, refreshTokenRepository);
const registerUseCase = new RegisterUseCase(userRepository, refreshTokenRepository);
const refreshUseCase = new RefreshUseCase(refreshTokenRepository);
const logoutUseCase = new LogoutUseCase(refreshTokenRepository);

const COOKIE_OPTIONS = {
  httpOnly: true,
  secure: process.env.NODE_ENV === 'production',
  sameSite: 'lax' as const,
};

function setAuthCookies(res: Response, tokens: { accessToken: string; refreshToken: string }) {
  res.cookie('access_token', tokens.accessToken, COOKIE_OPTIONS);
  res.cookie('refresh_token', tokens.refreshToken, COOKIE_OPTIONS);
}

export async function register(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const { email, password } = req.body as { email: string; password: string };
    const tokens = await registerUseCase.execute(email, password);
    setAuthCookies(res, tokens);
    res.status(201).json({ message: 'Registered successfully' });
  } catch (error) {
    next(error);
  }
}

export async function login(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const { email, password } = req.body as { email: string; password: string };
    const tokens = await loginUseCase.execute(email, password);
    setAuthCookies(res, tokens);
    res.json({ message: 'Logged in successfully' });
  } catch (error) {
    next(error);
  }
}

export async function refresh(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const refreshToken = req.cookies?.refresh_token;
    if (!refreshToken) {
      res.status(400).json({ error: 'Refresh token required' });
      return;
    }
    const tokens = await refreshUseCase.execute(refreshToken);
    setAuthCookies(res, tokens);
    res.json({ message: 'Tokens refreshed' });
  } catch (error) {
    next(error);
  }
}

export async function logout(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const refreshToken = req.cookies?.refresh_token;
    if (refreshToken) {
      await logoutUseCase.execute(refreshToken);
    }
    res.clearCookie('access_token');
    res.clearCookie('refresh_token');
    res.status(204).send();
  } catch (error) {
    next(error);
  }
}

export async function me(req: AuthRequest, res: Response): Promise<void> {
  res.json({ user: req.user });
}
