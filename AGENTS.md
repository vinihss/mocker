# AGENTS.md - OpenCode Multi-Agent Template

Template repo to copy to `~/.config/opencode/` or project `.opencode/` directory.

## Quick Start

```bash
cp -r * ~/.config/opencode/
```

## Enabled MCPs

- `sequential-thinking` - Chain-of-thought reasoning
- `fetch` - HTTP requests
- `context7` - Codebase context
- `git` - Git operations
- `playwright` - Browser automation
- `mcpbrowser` - Real browser with session
- `atlassian` - Jira + Confluence (OAuth required)
- `spec-kit` - Spec-driven development
- `browsermcp` - Browser MCP

Disabled (require API keys): `github`, `memory`, `docker`, `brave-search`, `sentry`, `postgres`, `notion`, `google-workspace`, `linear`

## Primary Agent

Use orchestrator as entry point: `opencode --agent orchestrator`

## Customization

Edit `agents/*.md` files to customize agent behavior. Each agent has YAML frontmatter defining its tools and permissions.

## Configuration

- Model: `opencode/big-pickle`
- Main config: `opencode.json`
- MCPs configured in `opencode.json` `"mcp"` section

## Environment

Optional env vars (not set in this repo):
- `${GITHUB_TOKEN}` - Enable GitHub MCP
- `${NOTION_TOKEN}` - Enable Notion MCP
- `${DATABASE_URL}` - Enable Postgres MCP