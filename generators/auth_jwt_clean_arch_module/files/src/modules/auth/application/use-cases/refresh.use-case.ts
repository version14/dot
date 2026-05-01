import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';
import { UnauthorizedError, NotFoundError } from '../../../../shared/errors';
import { signRefreshToken, signToken, verifyToken } from '../../../../shared/services/jwt';

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
      throw new UnauthorizedError('Invalid refresh token');
    }
    const stored = await this.refreshTokenRepository.findByToken(token);
    if (!stored) {
      throw new NotFoundError('Refresh token not found');
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
