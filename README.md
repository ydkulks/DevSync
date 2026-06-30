# DevSync

*Sync your understanding before you code.*

DevSync guides developers through codebase changes using Socratic questioning — ensuring deep understanding before any code is written or reviewed. It asks questions instead of dumping information, confirms understanding at each step, and only proceeds when the developer demonstrates comprehension.

## Single source of truth

[`SKILL.md`](SKILL.md) is the canonical body. All agent-specific copies are verified against it. Edit `SKILL.md`, then run:

```bash
make check
```

## Agent support

| Agent | File | Mechanism |
|-------|------|-----------|
| **Cursor** | `.cursor/rules/devsync.mdc` | Auto-loaded rule |
| **Windsurf** | `.windsurf/rules/devsync.md` | Auto-loaded rule |
| **Cline** | `.clinerules/devsync.md` | Auto-loaded rule |
| **GitHub Copilot** | `.github/copilot-instructions.md` | Auto-loaded instructions |
| **Claude Code** | `.agents/rules/devsync.md` | Auto-loaded rule |
| **OpenCode** | `.opencode/skills/devsync/SKILL.md` | Registered skill + frontmatter |
| **Kiro** | `.kiro/steering/devsync.md` | Auto-loaded steering |
| **Zed** | `.zed/devsync.md` | Auto-loaded instructions |
| **Gemini CLI** | `gemini-extension.json` (points at `SKILL.md`) | Extension install |

### Cursor

Copy `.cursor/rules/` into your project root.

### Windsurf

Copy `.windsurf/rules/` into your project root.

### Cline

Copy `.clinerules/` into your project root.

### GitHub Copilot

Copy `.github/copilot-instructions.md` into your project root.

### Claude Code

Copy `.agents/rules/` into your project root.

### OpenCode

Copy `.opencode/` into your project root. OpenCode auto-discovers skills under `.opencode/skills/`.

### Kiro

Copy `.kiro/steering/` into your project root, or to `~/.kiro/steering/` for global use.

### Zed

Copy `.zed/devsync.md` into your project root.

### Gemini CLI

```bash
gemini extensions install https://github.com/yd/devsync
```

## How it works

1. **Change discovery** — reads the relevant files, diffs, or planned changes and groups them logically.
2. **Socratic prompting** — starts each group with an open-ended question instead of an explanation.
3. **Listen, then build** — confirms correct understanding, gently redirects misconceptions, fills gaps only on request.
4. **No info-dumping** — every explanation is in direct response to the developer's answer or explicit request.
5. **Confirmation gate** — only proceeds when the developer demonstrates understanding.
6. **Session wrap** — asks if anything is still unclear before declaring done.

## License

MIT
