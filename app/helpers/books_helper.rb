module BooksHelper
  COVER_STYLE_OPTIONS = {
    "glass"     => { label: "Glass", hint: "Soft layered gradients" },
    "rings"     => { label: "Rings", hint: "Orbital geometry" },
    "shapes"    => { label: "Shapes", hint: "Modern abstract forms" },
    "identicon" => { label: "Identicon", hint: "Sharp symmetric marks" }
  }.freeze

  DICEBEAR_THEME_PALETTES = {
    "black"   => %w[111827 1f2937 374151],
    "blue"    => %w[1d4ed8 2563eb 60a5fa],
    "green"   => %w[166534 16a34a 4ade80],
    "magenta" => %w[9d174d db2777 f472b6],
    "orange"  => %w[c2410c f97316 fdba74],
    "violet"  => %w[5b21b6 7c3aed a78bfa],
    "white"   => %w[e5e7eb f3f4f6 ffffff]
  }.freeze

  THEME_COLORS = {
    "blue"    => { l: 69.14, c: 0.174, h: 245.01 },
    "orange"  => { l: 70.22, c: 0.2,   h: 45.1 },
    "magenta" => { l: 61.8,  c: 0.19,  h: 354.64 },
    "green"   => { l: 71.53, c: 0.175, h: 155.352 },
    "violet"  => { l: 56.13, c: 0.244, h: 297.99 },
    "white"   => { l: 92.0,  c: 0.0,   h: 0 },
    "black"   => { l: 20.0,  c: 0.0,   h: 0 }
  }.freeze

  COVER_COLS = 11
  COVER_ROWS = 16
  COVER_HALF = (COVER_COLS / 2.0).ceil # 6 independent columns, mirrored

  def book_cover_style_options
    COVER_STYLE_OPTIONS
  end

  def dicebear_base_url
    configured_url = Rails.application.config.x.dicebear_api_url
    base_url = configured_url if configured_url.is_a?(String) && configured_url.present?
    base_url ||= ENV["DICEBEAR_API_URL"].presence
    base_url ||= "https://api.dicebear.com/9.x"
    base_url.sub(%r{/*$}, "")
  end

  def dicebear_theme_palettes
    DICEBEAR_THEME_PALETTES
  end

  def book_cover_art(seed:, theme:, cover_style:, **options)
    if normalize_cover_style(cover_style) == "blocks"
      book_cover_svg(seed, **options)
    else
      image_tag(dicebear_cover_url(seed, theme:, cover_style:), { alt: "Book cover art" }.merge(options))
    end
  end

  def dicebear_cover_url(seed, theme:, cover_style:)
    style = normalize_cover_style(cover_style)
    raise ArgumentError, "blocks cover style is rendered locally" if style == "blocks"

    "#{dicebear_base_url}/#{style}/svg?#{dicebear_cover_params(seed, theme:).to_query}"
  end

  def book_cover_svg(seed, **options)
    digest = Digest::SHA256.hexdigest(seed.to_s) # SHA256 for more bits (256 vs 128)
    cells = compute_identicon_cells(digest)
    css_class = [ "book__cover", options[:class] ].compact.join(" ")
    svg_options = options.except(:class)

    tag.svg(
      {
        viewBox: "0 0 #{COVER_COLS} #{COVER_ROWS}",
        width: 1600,
        height: (1600.0 * COVER_ROWS / COVER_COLS).round,
        xmlns: "http://www.w3.org/2000/svg",
        class: css_class,
        role: "img",
        "aria-label": "Book cover pattern"
      }.merge(svg_options)
    ) do
      rects = tag.rect(x: 0, y: 0, width: COVER_COLS, height: COVER_ROWS, fill: "var(--cover-shade-bg)")
      cells.each do |cell|
        rects += tag.rect(x: cell[:x], y: cell[:y], width: 1, height: 1, fill: "var(--cover-shade-#{cell[:shade]})")
      end
      rects
    end
  end

  def cover_shade_css(theme)
    base = THEME_COLORS.fetch(theme, THEME_COLORS["blue"])
    if base[:c] == 0.0
      # Achromatic (black/white) — vary lightness only
      center = base[:l]
      shades = {
        bg: oklch(center - 15, 0, 0),
        1  => oklch(center - 8, 0, 0),
        2  => oklch(center, 0, 0),
        3  => oklch(center + 8, 0, 0),
        4  => oklch(center + 15, 0, 0)
      }
    else
      shades = {
        bg: oklch(base[:l] - 20, base[:c], base[:h]),
        1  => oklch(base[:l] - 10, base[:c] * 0.9, base[:h]),
        2  => oklch(base[:l], base[:c], base[:h]),
        3  => oklch(base[:l] + 10, base[:c] * 0.8, base[:h]),
        4  => oklch(base[:l] + 20, base[:c] * 0.6, base[:h])
      }
    end
    shades
  end

  def book_toc_tag(book, &)
    tag.ol class: "toc", tabindex: 0,
      data: {
        controller: "arrangement",
        action: arrangement_actions,
        arrangement_cursor_class: "arrangement-cursor",
        arrangement_selected_class: "arrangement-selected",
        arrangement_placeholder_class: "arrangement-placeholder",
        arrangement_move_mode_class: "arrangement-move-mode",
        arrangement_url_value: book_leaves_moves_url(book)
      }, &
  end

  def book_part_create_button(book, kind, **, &)
    url = url_for [ book, kind.new ]

    button_to url, class: "btn btn--plain txt-medium fill-transparent disable-when-arranging disable-when-deleting", draggable: true,
      data: {
        action: "dragstart->arrangement#dragStartCreate dragend->arrangement#dragEndCreate",
        arrangement_url_param: url
      }, **, &
  end

  def link_to_first_leafable(leaves)
    if first_leaf = leaves.first
      link_to leafable_slug_path(first_leaf), data: hotkey_data_attributes("right"), class: "disable-when-arranging", hidden: true do
        tag.span(class: "btn") do
          image_tag("arrow-right.svg", aria: { hidden: true }, size: 24) + tag.span("Start reading", class: "for-screen-reader")
        end + tag.span(first_leaf.title, class: "overflow-ellipsis")
      end
    end
  end

  def link_to_previous_leafable(leaf, hotkey: true, for_edit: false)
    if previous_leaf = leaf.previous
      path = for_edit ? edit_leafable_path(previous_leaf) : leafable_slug_path(previous_leaf)
      link_to path, data: hotkey_data_attributes("left", enabled: hotkey), class: "btn" do
        image_tag("arrow-left.svg", aria: { hidden: true }, size: 24) + tag.span("Previous: #{ previous_leaf.title }", class: "for-screen-reader")
      end
    else
      link_to book_slug_path(leaf.book), data: hotkey_data_attributes("left", enabled: hotkey), class: "btn" do
        image_tag("arrow-left.svg", aria: { hidden: true }, size: 24) + tag.span("Table of contents: #{ leaf.book.title }", class: "for-screen-reader")
      end
    end
  end

  def link_to_next_leafable(leaf, hotkey: true, for_edit: false)
    if next_leaf = leaf.next
      path = for_edit ? edit_leafable_path(next_leaf) : leafable_slug_path(next_leaf)
      link_to path, data: hotkey_data_attributes("right", enabled: hotkey), class: "btn txt-medium min-width" do
        tag.span("Next: #{next_leaf.title }", class: "overflow-ellipsis") + image_tag("arrow-right.svg", aria: { hidden: true }, size: 24)
      end
    else
      link_to book_slug_path(leaf.book), data: hotkey_data_attributes("right", enabled: hotkey), class: "btn txt-medium" do
        tag.span("Table of contents: #{leaf.book.title }", class: "overflow-ellipsis") + image_tag("arrow-reverse.svg", aria: { hidden: true }, size: 24)
      end
    end
  end

  private
    def hotkey_data_attributes(key, enabled: true)
      if enabled
        { controller: "hotkey", action: "keydown.#{key}@document->hotkey#click touch:swipe-#{key}@window->hotkey#click" }
      end
    end

    def compute_identicon_cells(digest)
      bytes = [ digest ].pack("H*").bytes # 32 bytes from SHA256
      cells = []
      bit_pool = bytes.inject(0) { |acc, b| (acc << 8) | b }
      bit_index = 0
      shade_index = 128 # use upper half for shade selection

      COVER_ROWS.times do |row|
        COVER_HALF.times do |col|
          filled = (bit_pool >> bit_index) & 1 == 1
          bit_index += 1

          next unless filled

          shade = ((bit_pool >> shade_index) & 3) + 1
          shade_index += 2

          cells << { x: col, y: row, shade: shade }
          mirror_col = COVER_COLS - 1 - col
          cells << { x: mirror_col, y: row, shade: shade } if mirror_col != col
        end
      end

      cells
    end

    def oklch(l, c, h)
      l = l.clamp(0.0, 100.0)
      c = c.clamp(0.0, 0.4)
      "oklch(#{l.round(1)}% #{c.round(3)} #{h.round(2)})"
    end

    def dicebear_cover_params(seed, theme:)
      {
        seed: seed.presence || SecureRandom.hex(8),
        backgroundColor: DICEBEAR_THEME_PALETTES.fetch(theme.to_s, DICEBEAR_THEME_PALETTES["blue"])
      }
    end

    def normalize_cover_style(cover_style)
      style = cover_style.to_s
      COVER_STYLE_OPTIONS.key?(style) ? style : "glass"
    end
end
