# Pages

Pages are markdown content leaves within a book — chapters, articles, notes.

## List pages

```bash
writebook pages list 1              # list pages in book 1
writebook pages list 1 --json
```

Output columns: ID, TITLE, STATUS, POSITION

## Create a page

```bash
writebook pages create 1 --title "Chapter 1" --body "# Hello\n\nMarkdown here"
writebook pages create 1 --title "Chapter 2" --body-file ./chapter2.md
writebook pages create 1 --title "Intro" --body "text" --position 1
```

| Flag | Required | Notes |
|------|----------|-------|
| `--title` | yes | |
| `--body` | no | Inline markdown string |
| `--body-file` | no | Read body from a file. Mutually exclusive with `--body` |
| `--position` | no | Insert at specific position. Without it, goes to the end |

**`--body` vs `--body-file`**: Use `--body` for short inline content. Use `--body-file` for anything longer than a sentence — it reads the file as-is, preserving all formatting.

## Show a page

```bash
writebook pages show 1 10           # book_id=1, page_id=10
writebook pages show 1 10 --json
```

Human-readable output shows title, status, position, and the full body text.

## Update a page

```bash
writebook pages update 1 10 --title "New Title"
writebook pages update 1 10 --body "Updated content"
writebook pages update 1 10 --body-file ./revised-chapter.md
writebook pages update 1 10 --title "Final" --body-file ./final.md
```

Only flags you pass are changed.

## Delete a page

```bash
writebook pages delete 1 10
```

This is a **soft delete** — sets status to "trashed". The leaf still exists in the database but won't appear in listings.

## Bulk import from a directory

```bash
BOOK_ID=$(writebook books create --title "Documentation" --json | jq .id)

for f in docs/*.md; do
  TITLE=$(head -1 "$f" | sed 's/^#\+ *//')
  writebook pages create "$BOOK_ID" --title "$TITLE" --body-file "$f"
done
```

## Reorder pages

Use `--position` on create to control ordering. Position values are floats — existing leaves shift to accommodate. To insert a page between positions 2 and 3, use `--position 2`.
