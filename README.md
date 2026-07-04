# DevSync

*Sync your understanding before you code.*

DevSync guides developers through codebase changes using Socratic questioning — ensuring deep understanding before any code is written or reviewed. It asks questions instead of dumping information, confirms understanding at each step, and only proceeds when the developer demonstrates comprehension.

## Single source of truth

Three canonical sources are verified by `make check`:

| Source | Purpose | Provider copies |
|--------|---------|----------------|
| `SKILL.md` | Socratic questioning skill (always-on) | 8 agent rule files |
| `commands/devsync.md` | `/devsync <path>` command (on-demand) | 8 provider command files |
| `commands/devsync-commit.md` | `/devsync-commit <ref>` command (on-demand) | 8 provider command files |

Edit either source, then run:

```bash
make check
```

## Agent support

| Agent | File | Mechanism |
|-------|------|-----------|
| **Cursor** | `.cursor/rules/devsync.mdc` | Auto-loaded rule |
| **Windsurf** | `.windsurf/rules/devsync.md` | Auto-loaded rule |
| **Cline** | `.clinerules/devsync.md` | Auto-loaded rule |
| **GitHub Copilot** | `.github/copilot-instructions.md` (+ `.github/prompts/` for commands) | Auto-loaded instructions + commands |
| **Claude Code** | `.agents/rules/devsync.md` | Auto-loaded rule |
| **OpenCode** | `.opencode/skills/devsync/SKILL.md` | Registered skill + frontmatter |
| **Kiro** | `.kiro/steering/devsync.md` | Auto-loaded steering |
| **Zed** | `.agents/skills/devsync/SKILL.md` (also reads `.github/copilot-instructions.md`) | Skill (on-demand via `/devsync`) + auto-loaded rules |
| **Gemini CLI** | `gemini-extension.json` (points at `SKILL.md`) | Extension install |

### Cursor

Copy `.cursor/rules/` into your project root.

### Windsurf

Copy `.windsurf/rules/` into your project root.

### Cline

Copy `.clinerules/` into your project root.

### GitHub Copilot

Copy these three files into your project — **do not** copy the entire `.github/` directory:

- `.github/copilot-instructions.md` — always-on instructions
- `.github/prompts/devsync.prompt.md` — `/devsync` command
- `.github/prompts/devsync-commit.prompt.md` — `/devsync-commit` command

`.github/` also contains `.github/workflows/check.yml` (CI workflow for this repo), which should not be copied into other projects.

### Claude Code

Copy `.agents/rules/` into your project root.

### OpenCode

Copy `.opencode/` into your project root. OpenCode auto-discovers skills under `.opencode/skills/` and commands under `.opencode/commands/`.

### Kiro

Copy `.kiro/steering/` into your project root, or to `~/.kiro/steering/` for global use.

### Zed

Zed auto-loads `.github/copilot-instructions.md` as project rules (already covered above). For the on-demand skill:

Copy `.agents/skills/` into your project root, or symlink from `~/.agents/skills/` for global use.

Zed does not support custom slash commands — use `/devsync` in the agent panel to invoke the skill.

### Gemini CLI

```bash
gemini extensions install https://github.com/yd/devsync
```

## Commands

Most providers support two commands for on-demand Socratic walkthroughs:

| Provider | Path walkthrough | Commit walkthrough |
|----------|-----------------|-------------------|
| **Cursor** | `/devsync` | `/devsync-commit` |
| **Windsurf** | `/devsync` | `/devsync-commit` |
| **Cline** | `/devsync` | `/devsync-commit` |
| **Claude Code** | `/devsync` | `/devsync-commit` |
| **GitHub Copilot Editor** | `/devsync` | `/devsync-commit` |
| **OpenCode** | `/devsync` | `/devsync-commit` |
| **Kiro** | `/devsync` | `/devsync-commit` |
| **Gemini CLI** | `/devsync` | `/devsync-commit` |

Zed and GitHub Copilot CLI do not support custom slash commands.

## How it works

1. **Change discovery** — reads the relevant files, diffs, or planned changes and groups them logically.
2. **Socratic prompting** — starts each group with an open-ended question instead of an explanation.
3. **Listen, then build** — confirms correct understanding, gently redirects misconceptions, fills gaps only on request.
4. **No info-dumping** — every explanation is in direct response to the developer's answer or explicit request.
5. **Confirmation gate** — only proceeds when the developer demonstrates understanding.
6. **Session wrap** — asks if anything is still unclear before declaring done.

## License

MIT
