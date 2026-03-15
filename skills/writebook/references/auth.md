# Authentication

The CLI needs a URL and bearer token stored in `~/.config/writebook/config.toml`. There are four ways to set this up.

## Join via invite link

If you have a join URL (e.g. from an admin), this creates your account and saves the token in one step:

```bash
writebook join http://localhost:3007/join/gFoO-0Lkb-UcFa
```

Interactive prompts for name, email, and password. Or pass flags for non-interactive use:

```bash
writebook join http://localhost:3007/join/gFoO-0Lkb-UcFa \
  --name "Build Agent" \
  --email agent@example.com \
  --password secret123
```

The command parses the URL to extract the server address and join code, creates the user account, obtains a bearer token, and saves everything to `~/.config/writebook/config.toml`.

## Login (existing account)

```bash
writebook login
```

Prompts for URL, email, and password interactively. Or pass flags:

```bash
writebook login --url https://books.example.com --email user@example.com --password secret
```

All three flags are optional — any flag you omit falls back to an interactive prompt.

## Environment variables

If you already have a token, skip login entirely:

```bash
export WRITEBOOK_URL=https://books.example.com
export WRITEBOOK_TOKEN=a1b2c3d4e5f6...
```

The CLI checks env vars (with `WRITEBOOK_` prefix) before the config file. This is useful for CI/CD pipelines or ephemeral environments where you don't want to write `~/.config/writebook/config.toml`.

## Check current auth status

```bash
cat ~/.config/writebook/config.toml 2>/dev/null
```

If it shows `url` and `token`, you're authenticated.

## Agentic usage

The recommended flow for agents is:

1. **First time** — join via invite link with flags:
   ```bash
   writebook join "$JOIN_URL" --name "Agent" --email agent@example.com --password secret123
   ```

2. **Subsequent runs** — the token is saved in `~/.config/writebook/config.toml`, so all commands just work:
   ```bash
   writebook books list
   ```

3. **Stateless environments** (containers, CI) — set env vars directly:
   ```bash
   export WRITEBOOK_URL=https://books.example.com
   export WRITEBOOK_TOKEN=existing-token-here
   writebook books list
   ```

## Token details

- Tokens are 64-character hex strings generated server-side
- They don't expire (persist until the user record changes)
- One token per user — logging in again returns the same token
- Stored in `~/.config/writebook/config.toml` as plain TOML:
  ```toml
  url = "https://books.example.com"
  token = "a1b2c3d4e5f6..."
  ```
