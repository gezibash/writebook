---
name: writebook
description: Manage Writebook via the CLI — create and publish books, add pages with markdown, organize with sections, upload images as pictures. Use this skill whenever the user mentions Writebook, wants to create or edit a book, add chapters or pages, upload images, publish content, manage sections, or work with book content. Also triggers for "the book app", "publishing", book authoring workflows, or importing markdown files into books.
---

# Writebook CLI

Writebook is a self-hosted book publishing app. The `writebook` CLI manages books, pages, sections, and pictures from the terminal.

## Reference files

Read these when you need exact flags, examples, or details for a specific resource:

- `references/auth.md` — login (interactive and non-interactive), env vars, agentic usage
- `references/books.md` — creating books, themes, cover styles, publishing
- `references/pages.md` — adding markdown content, bulk import from files
- `references/sections.md` — organizing books with dividers
- `references/pictures.md` — uploading images, bulk image import
- `references/search.md` — full-text search across books, global and book-scoped

## Quick orientation

A **Book** contains **Leaves**. Each leaf is one of:
- **Page** — markdown content (chapters, articles)
- **Section** — a divider/grouping with optional body and theme
- **Picture** — an uploaded image with optional caption

Leaves are ordered by **position** within a book.

## Setup

Check if the CLI is already configured:
```bash
cat ~/.config/writebook/config.toml 2>/dev/null
```

If it has `url` and `token`, you're ready. Otherwise:

**Join via invite link** (creates account + saves token in one step):
```bash
writebook join http://localhost:3007/join/gFoO-0Lkb-UcFa \
  --name "Agent" --email agent@example.com --password secret123
```

**Login with existing account:**
```bash
writebook login --url https://books.example.com --email user@example.com --password secret
```

See `references/auth.md` for env vars, agentic patterns, and token details.

All commands accept `--json` for machine-readable output.

## End-to-end example: create a book from scratch

```bash
# Create the book
writebook books create --title "Engineering Handbook" --author "Platform Team" --theme blue

# Add structure and content (using book ID from output)
writebook sections create 3 --title "Part 1: Getting Started" --theme green
writebook pages create 3 --title "Setup Guide" --body-file ./setup.md
writebook pages create 3 --title "Architecture" --body-file ./architecture.md
writebook pictures create 3 --title "System Diagram" --image ./diagram.png --caption "High-level architecture"

writebook sections create 3 --title "Part 2: Operations" --theme orange
writebook pages create 3 --title "Deployment" --body-file ./deploy.md
writebook pages create 3 --title "Monitoring" --body-file ./monitoring.md

# Publish
writebook books update 3 --published
```

## End-to-end example: import a docs folder

```bash
BOOK_ID=$(writebook books create --title "Documentation" --json | jq .id)

for f in docs/*.md; do
  TITLE=$(head -1 "$f" | sed 's/^#\+ *//')
  writebook pages create "$BOOK_ID" --title "$TITLE" --body-file "$f"
done

writebook books update "$BOOK_ID" --published
```

## Important patterns

**Discover leaf IDs** — use `books show` to see all leaves in a book before operating on them:
```bash
writebook books show 3
writebook books show 3 --json | jq '.leaves[] | {id, type, title}'
```

**Deletes are soft** for leaves (pages, sections, pictures set status to "trashed") but **hard** for books (permanently destroyed).

**Position** controls ordering. Pass `--position N` on create to place a leaf at a specific spot. Without it, new leaves go to the end.

**`--body-file`** for pages reads content from a file. Use it for anything longer than a sentence instead of `--body`.

**Pictures use file upload** — always pass `--image ./path.jpg`. Accepted formats: PNG, JPEG, WebP.

**Search** finds content across all books or within one: `writebook search "deploy"` or `writebook search "deploy" --book 3`. See `references/search.md`.
