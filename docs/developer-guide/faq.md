# FAQ

Frequently asked questions about developing in this project.

---

## General

**Q: Where do I start?**
See [Getting Started](../getting-started/README.md) for the full setup guide.

**Q: I found a bug. How do I report it?**
Open a [Bug Report issue](../../../issues/new/choose) using the provided template.

**Q: I want to add a feature. Where do I begin?**
Open a [Feature Request issue](../../../issues/new/choose) first to discuss the idea before writing code.

---

## Development

**Q: Tests are failing locally but passing in CI (or vice versa).**
- Make sure your `.env` matches the values expected by the test suite
- Run `pnpm install --frozen-lockfile` to ensure your dependencies match the lockfile

**Q: How do I run a single test file?**
```bash
pnpm vitest run path/to/your.test.ts
```

**Q: The build fails with type errors I didn't introduce.**
Pull the latest `main` and run `pnpm install` — a dependency may have been updated.

---

## Contributing

**Q: How large should a PR be?**
Aim for PRs that can be reviewed in under 30 minutes. Split larger changes into multiple PRs if possible.

**Q: Do I need to write tests for every change?**
Yes for new features and bug fixes. Documentation-only PRs are exempt.

**Q: Who merges PRs?**
Maintainers merge PRs once they have one approving review and all CI checks are green.

---

Still stuck? Open a [Discussion](../../../discussions).
