package authjwtcleanarchmodule

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

const userEntity = `export interface UserEntity {
  id: string;
  email: string;
  passwordHash: string;
  createdAt: Date;
  updatedAt: Date;
}
`

const userRepositoryInterface = `import type { UserEntity } from '../entities/user.entity';

export interface IUserRepository {
  findByEmail(email: string): Promise<UserEntity | null>;
  create(data: { email: string; passwordHash: string }): Promise<UserEntity>;
}
`

const refreshTokenRepositoryInterface = `export interface IRefreshTokenRepository {
  create(data: { token: string; userId: string; expiresAt: Date }): Promise<void>;
  findByToken(token: string): Promise<{ id: string; userId: string } | null>;
  deleteById(id: string): Promise<void>;
}
`

const loginUseCase = `import bcrypt from 'bcryptjs';
import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';
import type { IUserRepository } from '../../domain/interfaces/user.repository.interface';
import { signRefreshToken, signToken, verifyToken } from '../../../../lib/jwt';

export class LoginUseCase {
  constructor(
    private readonly userRepository: IUserRepository,
    private readonly refreshTokenRepository: IRefreshTokenRepository,
  ) {}

  async execute(email: string, password: string): Promise<{ accessToken: string; refreshToken: string }> {
    const user = await this.userRepository.findByEmail(email);
    if (!user) {
      throw new Error('Invalid credentials');
    }
    const valid = await bcrypt.compare(password, user.passwordHash);
    if (!valid) {
      throw new Error('Invalid credentials');
    }
    const accessToken = signToken({ id: user.id, email: user.email });
    const newRefreshToken = signRefreshToken({ id: user.id, email: user.email });
    const decoded = verifyToken<{ exp: number }>(newRefreshToken);
    const expiresAt = new Date(decoded.exp * 1000);
    await this.refreshTokenRepository.create({ token: newRefreshToken, userId: user.id, expiresAt });
    return { accessToken, refreshToken: newRefreshToken };
  }
}
`

const registerUseCase = `import bcrypt from 'bcryptjs';
import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';
import type { IUserRepository } from '../../domain/interfaces/user.repository.interface';
import { signRefreshToken, signToken, verifyToken } from '../../../../lib/jwt';

export class RegisterUseCase {
  constructor(
    private readonly userRepository: IUserRepository,
    private readonly refreshTokenRepository: IRefreshTokenRepository,
  ) {}

  async execute(email: string, password: string): Promise<{ accessToken: string; refreshToken: string }> {
    const existing = await this.userRepository.findByEmail(email);
    if (existing) {
      throw new Error('Email already in use');
    }
    const passwordHash = await bcrypt.hash(password, 10);
    const user = await this.userRepository.create({ email, passwordHash });
    const accessToken = signToken({ id: user.id, email: user.email });
    const newRefreshToken = signRefreshToken({ id: user.id, email: user.email });
    const decoded = verifyToken<{ exp: number }>(newRefreshToken);
    const expiresAt = new Date(decoded.exp * 1000);
    await this.refreshTokenRepository.create({ token: newRefreshToken, userId: user.id, expiresAt });
    return { accessToken, refreshToken: newRefreshToken };
  }
}
`

const refreshUseCase = `import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';
import { signRefreshToken, signToken, verifyToken } from '../../../../lib/jwt';

export class RefreshUseCase {
  constructor(private readonly refreshTokenRepository: IRefreshTokenRepository) {}

  async execute(token: string): Promise<{ accessToken: string; refreshToken: string }> {
    const getPayload = (): { id: string; email: string; exp: number } | null => {
      try {
        return verifyToken<{ id: string; email: string; exp: number }>(token);
      } catch {
        return null;
      }
    };
    const payload = getPayload();
    if (!payload) {
      throw new Error('Invalid refresh token');
    }
    const stored = await this.refreshTokenRepository.findByToken(token);
    if (!stored) {
      throw new Error('Refresh token not found');
    }
    await this.refreshTokenRepository.deleteById(stored.id);
    const newRefreshToken = signRefreshToken({ id: payload.id, email: payload.email });
    const decoded = verifyToken<{ exp: number }>(newRefreshToken);
    const expiresAt = new Date(decoded.exp * 1000);
    await this.refreshTokenRepository.create({
      token: newRefreshToken,
      userId: stored.userId,
      expiresAt,
    });
    const accessToken = signToken({ id: payload.id, email: payload.email });
    return { accessToken, refreshToken: newRefreshToken };
  }
}
`

const logoutUseCase = `import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';

export class LogoutUseCase {
  constructor(private readonly refreshTokenRepository: IRefreshTokenRepository) {}

  async execute(token: string): Promise<void> {
    const stored = await this.refreshTokenRepository.findByToken(token);
    if (stored) {
      await this.refreshTokenRepository.deleteById(stored.id);
    }
  }
}
`

const authController = `import type { NextFunction, Request, Response } from 'express';
import type { AuthRequest } from '../../../../middleware/auth.middleware';
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

export async function register(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const { email, password } = req.body as { email: string; password: string };
    const tokens = await registerUseCase.execute(email, password);
    res.status(201).json(tokens);
  } catch (error) {
    next(error);
  }
}

export async function login(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const { email, password } = req.body as { email: string; password: string };
    const tokens = await loginUseCase.execute(email, password);
    res.json(tokens);
  } catch (error) {
    next(error);
  }
}

export async function refresh(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const { refreshToken } = req.body as { refreshToken?: string };
    if (!refreshToken) {
      res.status(400).json({ error: 'Refresh token required' });
      return;
    }
    const tokens = await refreshUseCase.execute(refreshToken);
    res.json(tokens);
  } catch (error) {
    next(error);
  }
}

export async function logout(req: Request, res: Response, next: NextFunction): Promise<void> {
  try {
    const { refreshToken } = req.body as { refreshToken?: string };
    if (refreshToken) {
      await logoutUseCase.execute(refreshToken);
    }
    res.status(204).send();
  } catch (error) {
    next(error);
  }
}

export async function me(req: AuthRequest, res: Response): Promise<void> {
  res.json({ user: req.user });
}
`

const authRouteCleanArch = `import { Router } from 'express';
import { login, logout, me, refresh, register } from '../modules/auth/application/controllers/auth.controller';
import { authMiddleware } from '../middleware/auth.middleware';

const router = Router();

router.post('/register', register);
router.post('/login', login);
router.get('/me', authMiddleware, me);
router.post('/refresh', refresh);
router.post('/logout', logout);

export default router;
`

const userRepositoryImpl = `import { eq } from 'drizzle-orm';
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
`

const refreshTokenRepositoryImpl = `import { eq } from 'drizzle-orm';
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
`

func (g *Generator) Generate(ctx *dotapi.Context) error {
	files := map[string]string{
		"src/modules/auth/domain/entities/user.entity.ts":                                   userEntity,
		"src/modules/auth/domain/interfaces/user.repository.interface.ts":                   userRepositoryInterface,
		"src/modules/auth/domain/interfaces/refresh-token.repository.interface.ts":          refreshTokenRepositoryInterface,
		"src/modules/auth/application/use-cases/login.use-case.ts":                          loginUseCase,
		"src/modules/auth/application/use-cases/register.use-case.ts":                       registerUseCase,
		"src/modules/auth/application/use-cases/refresh.use-case.ts":                        refreshUseCase,
		"src/modules/auth/application/use-cases/logout.use-case.ts":                         logoutUseCase,
		"src/modules/auth/application/controllers/auth.controller.ts":                       authController,
		"src/routes/auth.route.ts":                                                          authRouteCleanArch,
		"src/modules/auth/infrastructure/database/repositories/user.repository.ts":          userRepositoryImpl,
		"src/modules/auth/infrastructure/database/repositories/refresh-token.repository.ts": refreshTokenRepositoryImpl,
	}

	for path, content := range files {
		ctx.State.WriteFile(path, []byte(content), state.ContentRaw)
	}

	return ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"bcryptjs": "^2.4.3",
			},
			"devDependencies": map[string]interface{}{
				"@types/bcryptjs": "^2.4.6",
			},
		})
		return nil
	})
}
