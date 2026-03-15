# Sections

Sections are dividers/groupings within a book. They have an optional body and their own theme color.

## List sections

```bash
writebook sections list 1
writebook sections list 1 --json
```

Output columns: ID, TITLE, THEME, STATUS, POSITION

## Create a section

```bash
writebook sections create 1 --title "Part 1: Getting Started"
writebook sections create 1 --title "Part 2" --body "Advanced topics" --theme green
writebook sections create 1 --title "Appendix" --position 10
```

| Flag | Required | Notes |
|------|----------|-------|
| `--title` | yes | |
| `--body` | no | Optional section body text |
| `--theme` | no | `black`, `blue`, `green`, `magenta`, `orange`, `violet`, `white` |
| `--position` | no | Insert at specific position |

Sections are visual dividers in the book's table of contents. Use them to group related pages under a heading.

## Show a section

```bash
writebook sections show 1 11
writebook sections show 1 11 --json
```

## Update a section

```bash
writebook sections update 1 11 --title "Part 1: Fundamentals"
writebook sections update 1 11 --theme magenta --body "Updated intro"
```

## Delete a section

```bash
writebook sections delete 1 11
```

Soft delete — sets status to "trashed".

## Typical book structure

A well-organized book alternates sections and pages:

```bash
BOOK_ID=3

writebook sections create $BOOK_ID --title "Part 1: Setup" --theme blue
writebook pages create $BOOK_ID --title "Installation" --body-file ./install.md
writebook pages create $BOOK_ID --title "Configuration" --body-file ./config.md

writebook sections create $BOOK_ID --title "Part 2: Usage" --theme green
writebook pages create $BOOK_ID --title "Basic Commands" --body-file ./basics.md
writebook pages create $BOOK_ID --title "Advanced" --body-file ./advanced.md
```
