# Books

## List books

```bash
writebook books list
writebook books list --json
```

Output columns: ID, TITLE, AUTHOR, PUBLISHED, THEME, COVER STYLE

## Create a book

```bash
writebook books create --title "Engineering Handbook" --author "Platform Team"
writebook books create --title "Design Guide" --theme magenta --cover-style identicon
writebook books create --title "API Docs" --cover-seed "a1b2c3d4e5f6g7h8" --json
```

| Flag | Required | Notes |
|------|----------|-------|
| `--title` | yes | |
| `--subtitle` | no | |
| `--author` | no | |
| `--theme` | no | `black`, `blue`, `green`, `magenta`, `orange`, `violet`, `white`. Default: `blue` |
| `--cover-style` | no | `glass`, `rings`, `shapes`, `identicon`. Default: `glass` |
| `--cover-seed` | no | 16-char hex string. Auto-generated if omitted |

## Show a book

```bash
writebook books show 1
writebook books show 1 --json
```

Shows book metadata plus a table of all its leaves (pages, sections, pictures) with their IDs, types, and titles. Use this to discover leaf IDs before operating on them.

## Update a book

```bash
writebook books update 1 --title "New Title"
writebook books update 1 --published
writebook books update 1 --theme magenta --cover-style rings
writebook books update 1 --published=false
```

| Flag | Notes |
|------|-------|
| `--title` | |
| `--subtitle` | |
| `--author` | |
| `--theme` | |
| `--cover-style` | |
| `--cover-seed` | |
| `--published` | Boolean flag. Toggles visibility to readers |

Only flags you pass are changed — everything else stays the same.

## Delete a book

```bash
writebook books delete 1
writebook books delete 1 --force   # skip confirmation prompt
```

This is a **hard delete** — the book and all its leaves are permanently destroyed.

## Themes and covers

**Themes** set the color scheme. The same theme values work on both books and sections.

**Cover styles** determine the generative art pattern:
- `glass` — soft layered gradients (default)
- `rings` — orbital geometry
- `shapes` — modern abstract forms
- `identicon` — sharp symmetric block pattern (rendered as SVG)

The `cover_seed` is deterministic: same seed + style + theme = identical cover every time. Omit it and the server generates a random one.
