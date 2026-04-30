package authjwtusersschema

import (
	"strings"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

const usersTable = `import { pgTable, text, timestamp, uuid } from 'drizzle-orm/pg-core';

export const users = pgTable('users', {
  id: uuid('id').primaryKey().defaultRandom(),
  email: text('email').notNull().unique(),
  passwordHash: text('password_hash').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
  updatedAt: timestamp('updated_at').defaultNow().notNull(),
});

export const refreshTokens = pgTable('refresh_tokens', {
  id: uuid('id').primaryKey().defaultRandom(),
  token: text('token').notNull().unique(),
  userId: uuid('user_id')
    .notNull()
    .references(() => users.id, { onDelete: 'cascade' }),
  expiresAt: timestamp('expires_at').notNull(),
  createdAt: timestamp('created_at').defaultNow().notNull(),
});

export type User = typeof users.$inferSelect;
export type NewUser = typeof users.$inferInsert;
export type RefreshToken = typeof refreshTokens.$inferSelect;
`

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("src/db/schema/users.table.ts", []byte(usersTable), state.ContentRaw)

	existing := "export {};\n"
	if f, ok := ctx.State.GetFile("src/db/schema/index.ts"); ok {
		existing = string(f.Content)
	}
	if strings.TrimSpace(existing) == "export {};" {
		existing = ""
	}
	updated := existing + "export * from './users.table';\n"
	ctx.State.WriteFile("src/db/schema/index.ts", []byte(updated), state.ContentRaw)
	return nil
}
