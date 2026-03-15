# Writebook

A fork of [ONCE Writebook](https://once.com/writebook) by 37signals — the simple, self-hosted book publishing tool — extended with a **REST API** and **Go CLI** for agentic and programmatic workflows.

## What's Writebook?

Writebook is part of the [ONCE](https://once.com) lineup from 37signals. You buy it once, own the code, and host it yourself. It lets you write and publish books on the web — no CMS, no static site generator, just a focused Rails app for putting words in front of readers.

This fork keeps all of that and adds the tooling needed to manage books without a browser.

## What This Fork Adds

### REST API (`/api/v1/`)

A full JSON API for managing books, pages, and sections programmatically:

| Endpoint | Methods |
|---|---|
| `/api/v1/tokens` | `POST` — exchange email/password for a bearer token |
| `/api/v1/books` | `GET` `POST` |
| `/api/v1/books/:id` | `GET` `PATCH` `DELETE` |
| `/api/v1/books/:book_id/pages` | `POST` |
| `/api/v1/books/:book_id/pages/:id` | `GET` `PATCH` `DELETE` |
| `/api/v1/books/:book_id/sections` | `POST` |
| `/api/v1/books/:book_id/sections/:id` | `GET` `PATCH` `DELETE` |

Token auth via `Authorization: Bearer <token>` header. All responses are JSON.

### Go CLI (`cli/`)

A command-line interface built with Cobra and Viper:

```bash
writebook login                              # Authenticate and save token
writebook books list                         # List all books
writebook books create --title "My Book"     # Create a book
writebook books show 1                       # Show book with table of contents
writebook pages create 1 --title "Ch 1" --body "# Hello"
writebook pages create 1 --title "Ch 2" --body-file chapter2.md
writebook sections create 1 --title "Part 2" --theme violet
writebook books delete 1                     # With confirmation prompt
```

Every command supports `--json` for machine-readable output — useful for piping into `jq` or driving from scripts and AI agents.

### Agentic Workflows

The API and CLI were designed with AI agents in mind. An agent can:

1. Authenticate via the API
2. Create and structure books programmatically
3. Write content page by page (accepting markdown)
4. Use `--json` output to parse results and make decisions

No browser needed. No manual clicking. Just API calls.

## Quick Start

### Docker (recommended)

```bash
docker run -d --name writebook -p 3007:80 \
  -e SECRET_KEY_BASE=$(openssl rand -hex 64) \
  -e DISABLE_SSL=1 \
  ghcr.io/gezibash/writebook:latest
```

Open `http://localhost:3007` to complete first-run setup.

### Install the CLI

```bash
curl -fsSL https://raw.githubusercontent.com/gezibash/writebook/main/install.sh | sh
```

Or download a binary from [Releases](https://github.com/gezibash/writebook/releases).

Then:

```bash
writebook login
writebook books list
```

### From Source

```bash
# Rails app
mise exec -- bin/setup
mise exec -- bin/dev

# CLI
cd cli
go build -o writebook .
```

## Cover Styles

Generated covers now support a curated set of styles: `blocks`, `glass`, `rings`, `shapes`, and `identicon`.

By default, DiceBear-backed styles use the public v9 API. To point Writebook at a self-hosted DiceBear instance instead, set:

```bash
export DICEBEAR_API_URL=https://dicebear.example.com/9.x
```

The value should be the API root that already includes the DiceBear version segment.

## Configuration

The CLI reads config from environment variables or `~/.writebook.yaml`:

| Variable | Description |
|---|---|
| `WRITEBOOK_URL` | Base URL of your Writebook instance |
| `WRITEBOOK_TOKEN` | Bearer token (obtained via `writebook login`) |

`writebook login` saves both to `~/.writebook.yaml` automatically.

## Release Process

Push a version tag to trigger the release pipeline:

```bash
git tag v1.1.0
git push origin v1.1.0
```

This builds:
- **CLI binaries** for linux/darwin/windows on amd64/arm64 (via GoReleaser)
- **Docker images** for linux/amd64 and linux/arm64 (pushed to GHCR)

## Stack

- **Rails 8** with SQLite — the Writebook app
- **Go** with Cobra/Viper — the CLI
- **GitHub Actions** with GoReleaser — CI/CD
- **Docker** with multi-arch builds — deployment

## License

Writebook is a [ONCE](https://once.com) product by 37signals. See your ONCE license for terms. The API and CLI additions in this fork are provided as-is.
