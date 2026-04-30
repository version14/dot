package authjwtmvcroute

import (
	"slices"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

const authRoute = `import { Router } from 'express';
import { login, logout, me, refresh, register } from '../controllers/auth.controller';
import { authMiddleware } from '../middleware/auth.middleware';

const router = Router();

router.post('/register', register);
router.post('/login', login);
router.get('/me', authMiddleware, me);
router.post('/refresh', refresh);
router.post('/logout', logout);

export default router;
`

const authControllerWithDB = `import type { Request, Response } from 'express';
import type { AuthRequest } from '../middleware/auth.middleware';
import bcrypt from 'bcryptjs';
import { eq } from 'drizzle-orm';
import { db } from '../db';
import { refreshTokens, users } from '../db/schema';
import { signRefreshToken, signToken, verifyToken } from '../lib/jwt';

async function issueTokens(
  userId: string,
  email: string,
): Promise<{ accessToken: string; refreshToken: string }> {
  const accessToken = signToken({ id: userId, email });
  const newRefreshToken = signRefreshToken({ id: userId, email });
  const decoded = verifyToken<{ exp: number }>(newRefreshToken);
  const expiresAt = new Date(decoded.exp * 1000);
  await db.insert(refreshTokens).values({ token: newRefreshToken, userId, expiresAt });
  return { accessToken, refreshToken: newRefreshToken };
}

export async function register(req: Request, res: Response): Promise<void> {
  const { email, password } = req.body as { email: string; password: string };
  const existing = await db.select().from(users).where(eq(users.email, email)).limit(1);
  if (existing.length > 0) {
    res.status(409).json({ error: 'Email already in use' });
    return;
  }
  const passwordHash = await bcrypt.hash(password, 10);
  const result = await db.insert(users).values({ email, passwordHash }).returning();
  const newUser = result[0];
  if (!newUser) {
    res.status(500).json({ error: 'Failed to create user' });
    return;
  }
  const tokens = await issueTokens(newUser.id, newUser.email);
  res.status(201).json(tokens);
}

export async function login(req: Request, res: Response): Promise<void> {
  const { email, password } = req.body as { email: string; password: string };
  const result = await db.select().from(users).where(eq(users.email, email)).limit(1);
  const user = result[0];
  if (!user) {
    res.status(401).json({ error: 'Invalid credentials' });
    return;
  }
  const valid = await bcrypt.compare(password, user.passwordHash);
  if (!valid) {
    res.status(401).json({ error: 'Invalid credentials' });
    return;
  }
  const tokens = await issueTokens(user.id, user.email);
  res.json(tokens);
}

export async function refresh(req: Request, res: Response): Promise<void> {
  const { refreshToken } = req.body as { refreshToken?: string };
  if (!refreshToken) {
    res.status(400).json({ error: 'Refresh token required' });
    return;
  }
  const getPayload = (): { id: string; email: string; exp: number } | null => {
    try {
      return verifyToken<{ id: string; email: string; exp: number }>(refreshToken);
    } catch {
      return null;
    }
  };
  const payload = getPayload();
  if (!payload) {
    res.status(401).json({ error: 'Invalid refresh token' });
    return;
  }
  const stored = await db
    .select()
    .from(refreshTokens)
    .where(eq(refreshTokens.token, refreshToken))
    .limit(1);
  const storedToken = stored[0];
  if (!storedToken) {
    res.status(401).json({ error: 'Refresh token not found' });
    return;
  }
  await db.delete(refreshTokens).where(eq(refreshTokens.id, storedToken.id));
  const tokens = await issueTokens(payload.id, payload.email);
  res.json(tokens);
}

export async function logout(req: Request, res: Response): Promise<void> {
  const { refreshToken } = req.body as { refreshToken?: string };
  if (refreshToken) {
    await db.delete(refreshTokens).where(eq(refreshTokens.token, refreshToken));
  }
  res.status(204).send();
}

export async function me(req: AuthRequest, res: Response): Promise<void> {
  res.json({ user: req.user });
}
`

const authControllerStub = `import type { Request, Response } from 'express';
import type { AuthRequest } from '../middleware/auth.middleware';

export async function register(_req: Request, res: Response): Promise<void> {
  res.status(501).json({ error: 'not implemented' });
}

export async function login(_req: Request, res: Response): Promise<void> {
  res.status(501).json({ error: 'not implemented' });
}

export async function refresh(_req: Request, res: Response): Promise<void> {
  res.status(501).json({ error: 'not implemented' });
}

export async function logout(_req: Request, res: Response): Promise<void> {
  res.status(501).json({ error: 'not implemented' });
}

export async function me(req: AuthRequest, res: Response): Promise<void> {
  res.json({ user: req.user });
}
`

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("src/routes/auth.route.ts", []byte(authRoute), state.ContentRaw)

	hasDB := slices.Contains(ctx.PreviousGens, "drizzle_postgres_adapter")
	if hasDB {
		ctx.State.WriteFile("src/controllers/auth.controller.ts", []byte(authControllerWithDB), state.ContentRaw)
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
	ctx.State.WriteFile("src/controllers/auth.controller.ts", []byte(authControllerStub), state.ContentRaw)
	return nil
}
