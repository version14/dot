import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/*.{unit,feature,e2e,db}.test.ts'],
    passWithNoTests: true,
    fileParallelism: false,
    env: {
      JWT_SECRET: 'test-secret-vitest',
      JWT_EXPIRES_IN: '1h',
      JWT_REFRESH_EXPIRES_IN: '7d',
    },
  },
});
