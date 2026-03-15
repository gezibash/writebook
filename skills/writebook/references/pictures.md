# Pictures

Pictures are image leaves within a book — photos, diagrams, screenshots. They support PNG, JPEG, and WebP formats.

## List pictures

```bash
writebook pictures list 1
writebook pictures list 1 --json
```

Output columns: ID, TITLE, CAPTION, HAS IMAGE, POSITION

## Create a picture

```bash
writebook pictures create 1 --title "Architecture Diagram" --image ./diagram.png
writebook pictures create 1 --title "Photo" --image ./photo.jpg --caption "Team offsite 2026"
writebook pictures create 1 --title "Logo" --image ./logo.webp --position 1
```

| Flag | Required | Notes |
|------|----------|-------|
| `--title` | yes | |
| `--image` | yes | Path to image file (PNG, JPEG, WebP) |
| `--caption` | no | Displayed below the image |
| `--position` | no | Insert at specific position |

The server generates a `:large` variant resized to fit within 1500x1500px.

## Show a picture

```bash
writebook pictures show 1 12
writebook pictures show 1 12 --json
```

JSON output includes `has_image` (boolean) and `image_url` (server path to the stored image).

## Update a picture

```bash
writebook pictures update 1 12 --caption "Updated caption"
writebook pictures update 1 12 --image ./better-photo.jpg --title "Better Photo"
writebook pictures update 1 12 --title "Renamed"
```

You can replace the image, update the caption, change the title, or any combination.

## Delete a picture

```bash
writebook pictures delete 1 12
```

Soft delete — sets status to "trashed".

## Bulk import images from a folder

```bash
BOOK_ID=3

for img in screenshots/*.{jpg,png,webp}; do
  [ -f "$img" ] || continue
  NAME=$(basename "$img" | sed 's/\.[^.]*$//' | tr '-_' '  ')
  writebook pictures create "$BOOK_ID" --title "$NAME" --image "$img"
done
```

## Add images alongside markdown content

```bash
BOOK_ID=3

writebook pages create $BOOK_ID --title "Setup Guide" --body-file ./setup.md
writebook pictures create $BOOK_ID --title "Setup Screenshot" --image ./setup-screenshot.png --caption "The setup wizard"
writebook pages create $BOOK_ID --title "Configuration" --body-file ./config.md
```

Pictures appear inline in the book between pages, so position them where they make sense in the reading flow.
