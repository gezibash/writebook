# Search

Full-text search across book content using FTS5 with porter stemming and BM25 ranking.

## Global search

```bash
writebook search "deployment process"
writebook search "deployment process" --json
```

Searches all accessible books and groups results by book.

## Book-scoped search

```bash
writebook search "deployment" --book 3
writebook search "deployment" --book 3 --json
```

Searches within a single book. Returns a flat list of results.

## Flags

| Flag | Required | Notes |
|------|----------|-------|
| `<query>` | yes | Positional argument — the search terms |
| `--book` | no | Limit search to a specific book ID |
| `--json` | no | Output raw JSON from the API |

## Output format

**Global (plain text)** — grouped by book:
```
Handbook (2 results)
  42   Page     Deploy Guide
                ...deployment process starts with...
  55   Section  Operations
                ...deployment checklist and...

Manual (1 result)
  12   Page     Getting Started
                ...initial deployment steps...
```

**Book-scoped (plain text)** — tabwriter table:
```
ID  TYPE     TITLE           SNIPPET
42  Page     Deploy Guide    ...deployment process starts with...
```

**JSON (global)** — array of book groups with nested results:
```json
[
  {
    "book_id": 1,
    "book_title": "Handbook",
    "results": [
      { "id": 42, "title": "Deploy Guide", "title_snippet": "Deploy Guide",
        "content_snippet": "...deployment process starts with...",
        "type": "Page", "book_id": 1, "book_title": "Handbook" }
    ]
  }
]
```

**JSON (book-scoped)** — flat array of results.

## Search tips

- Multiple words are ANDed together: `deploy staging` finds leaves containing both words
- Phrase search with quotes: `"deployment process"` matches the exact phrase
- Results are ranked by BM25 with title matches favored over content matches
- Results are limited to 50 per book
- Only active (non-trashed) leaves are searched

## API endpoints

```
GET /api/v1/search?q=<query>               # global
GET /api/v1/books/<id>/search?q=<query>     # book-scoped
```

Both return 422 if `q` is missing or blank.
