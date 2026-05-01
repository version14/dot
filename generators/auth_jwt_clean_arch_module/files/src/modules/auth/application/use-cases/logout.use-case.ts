import type { IRefreshTokenRepository } from '../../domain/interfaces/refresh-token.repository.interface';

export class LogoutUseCase {
  constructor(private readonly refreshTokenRepository: IRefreshTokenRepository) {}

  async execute(token: string): Promise<void> {
    const stored = await this.refreshTokenRepository.findByToken(token);
    if (stored) {
      await this.refreshTokenRepository.deleteById(stored.id);
    }
  }
}
