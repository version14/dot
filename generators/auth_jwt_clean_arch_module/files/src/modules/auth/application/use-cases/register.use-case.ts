import bcrypt from 'bcryptjs';
import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';
import type { IUserRepository } from '../../domain/interfaces/user.repository.interface';
import { ConflictError } from '../../../../shared/errors';
import { signRefreshToken, signToken, verifyToken } from '../../../../shared/services/jwt';

export class RegisterUseCase {
  constructor(
    private readonly userRepository: IUserRepository,
    private readonly refreshTokenRepository: IRefreshTokenRepository,
  ) {}

  async execute(email: string, password: string): Promise<{ accessToken: string; refreshToken: string }> {
    const existing = await this.userRepository.findByEmail(email);
    if (existing) {
      throw new ConflictError('Email already in use');
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
