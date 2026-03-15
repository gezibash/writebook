import { Controller } from "@hotwired/stimulus"

const COVER_COLS = 11
const COVER_ROWS = 16
const COVER_HALF = Math.ceil(COVER_COLS / 2)

export default class extends Controller {
  static targets = [
    "title", "titleInput",
    "subtitle", "subtitleInput",
    "author", "authorInput",
    "themeInput", "styleInput",
    "styleName", "styleHint",
    "seedInput", "seedValue",
    "blocksArt", "dicebearArt"
  ]

  static values = { dicebearEndpoint: String, themeColors: Object, styleOptions: Object }

  connect() {
    this.blocksRenderRequest = 0
    this.lastBlocksSignature = null
    this.lastDicebearUrl = null
    this.update()
  }

  update() {
    this.ensureSeed()
    this.sync("title")
    this.sync("subtitle")
    this.sync("author")
    this.updateSeedDisplay()
    this.updateThemeClass()
    this.updateStyleDisplay()
    this.updateArt()
  }

  refreshSeed() {
    if (!this.hasSeedInputTarget) return

    this.seedInputTarget.value = this.randomSeed()
    this.lastBlocksSignature = null
    this.lastDicebearUrl = null
    this.update()
  }

  cycleStyle(event) {
    if (this.styleInputTargets.length === 0) return

    const direction = Number.parseInt(event.currentTarget.dataset.direction || "1", 10)
    const currentIndex = this.styleInputTargets.findIndex((input) => input.checked)
    const startIndex = currentIndex >= 0 ? currentIndex : 0
    const nextIndex = (startIndex + direction + this.styleInputTargets.length) % this.styleInputTargets.length

    this.styleInputTargets[nextIndex].checked = true
    this.update()
  }

  sync(name) {
    const sourceTarget = `${name}Input`

    if (!this[`has${this.classify(name)}Target`] || !this[`has${this.classify(sourceTarget)}Target`]) return

    const input = this[`${sourceTarget}Target`]
    const value = input.value.trim()
    const text = value.length > 0 ? input.value : input.placeholder

    this[`${name}Targets`].forEach((element) => {
      element.textContent = text
      element.dataset.empty = value.length === 0
    })
  }

  updateArt() {
    const style = this.currentStyle()
    const isBlocks = style === "blocks"

    if (this.hasBlocksArtTarget) this.blocksArtTarget.hidden = !isBlocks

    if (isBlocks) {
      this.renderBlocksArt()
      if (this.hasDicebearArtTarget) this.dicebearArtTarget.hidden = true
      return
    }

    if (!this.hasDicebearArtTarget) return

    const url = this.dicebearUrl(style)
    this.dicebearArtTarget.hidden = false

    if (url !== this.lastDicebearUrl) {
      this.dicebearArtTarget.src = url
      this.lastDicebearUrl = url
    }
  }

  updateSeedDisplay() {
    if (this.hasSeedValueTarget) this.seedValueTarget.textContent = this.seed()
  }

  updateThemeClass() {
    Array.from(this.element.classList)
      .filter((className) => className.startsWith("theme--"))
      .forEach((className) => this.element.classList.remove(className))

    this.element.classList.add(`theme--${this.currentTheme()}`)
  }

  updateStyleDisplay() {
    const details = this.styleOptionsValue[this.currentStyle()] || {}

    if (this.hasStyleNameTarget) this.styleNameTarget.textContent = details.label || this.titleize(this.currentStyle())
    if (this.hasStyleHintTarget) this.styleHintTarget.textContent = details.hint || ""
  }

  async renderBlocksArt() {
    if (!this.hasBlocksArtTarget) return

    const signature = `${this.seed()}:${this.currentTheme()}`
    if (signature === this.lastBlocksSignature) return

    const request = ++this.blocksRenderRequest
    const digestBytes = await this.digestBytes(this.seed())

    if (request !== this.blocksRenderRequest) return

    const cells = this.computeIdenticonCells(digestBytes)
    const rects = [ '<rect x="0" y="0" width="11" height="16" fill="var(--cover-shade-bg)"></rect>' ]

    cells.forEach((cell) => {
      rects.push(`<rect x="${cell.x}" y="${cell.y}" width="1" height="1" fill="var(--cover-shade-${cell.shade})"></rect>`)
    })

    this.blocksArtTarget.innerHTML = rects.join("")
    this.lastBlocksSignature = signature
  }

  async digestBytes(value) {
    if (globalThis.crypto?.subtle) {
      const encoded = new TextEncoder().encode(value)
      const digest = await globalThis.crypto.subtle.digest("SHA-256", encoded)
      return Array.from(new Uint8Array(digest))
    }

    return this.fallbackDigestBytes(value)
  }

  fallbackDigestBytes(value) {
    const bytes = new Array(32).fill(0)

    Array.from(value).forEach((character, index) => {
      const code = character.codePointAt(0) || 0
      const slot = index % bytes.length
      bytes[slot] = (bytes[slot] ^ code ^ ((slot + 1) * 17)) & 0xff
    })

    return bytes
  }

  computeIdenticonCells(bytes) {
    let bitPool = 0n

    bytes.forEach((byte) => {
      bitPool = (bitPool << 8n) | BigInt(byte)
    })

    const cells = []
    let bitIndex = 0
    let shadeIndex = 128

    for (let row = 0; row < COVER_ROWS; row += 1) {
      for (let col = 0; col < COVER_HALF; col += 1) {
        const filled = ((bitPool >> BigInt(bitIndex)) & 1n) === 1n
        bitIndex += 1

        if (!filled) continue

        const shade = Number(((bitPool >> BigInt(shadeIndex)) & 3n) + 1n)
        shadeIndex += 2

        cells.push({ x: col, y: row, shade })

        const mirrorCol = COVER_COLS - 1 - col
        if (mirrorCol !== col) cells.push({ x: mirrorCol, y: row, shade })
      }
    }

    return cells
  }

  dicebearUrl(style) {
    const base = this.dicebearEndpointValue.replace(/\/+$/, "")
    const url = new URL(`${base}/${style}/svg`)

    url.searchParams.set("seed", this.seed())
    this.palette().forEach((color) => url.searchParams.append("backgroundColor[]", color))

    return url.toString()
  }

  seed() {
    return this.hasSeedInputTarget ? this.seedInputTarget.value.trim() : ""
  }

  ensureSeed() {
    if (!this.hasSeedInputTarget) return

    if (this.seedInputTarget.value.trim().length === 0) {
      this.seedInputTarget.value = this.randomSeed()
    }
  }

  randomSeed() {
    const bytes = new Uint8Array(8)

    if (globalThis.crypto?.getRandomValues) {
      globalThis.crypto.getRandomValues(bytes)
    } else {
      bytes.forEach((_, index) => {
        bytes[index] = Math.floor(Math.random() * 256)
      })
    }

    return Array.from(bytes, (byte) => byte.toString(16).padStart(2, "0")).join("")
  }

  palette() {
    return this.themeColorsValue[this.currentTheme()] || this.themeColorsValue.blue || []
  }

  currentTheme() {
    return this.checkedValue(this.themeInputTargets, "blue")
  }

  currentStyle() {
    return this.checkedValue(this.styleInputTargets, "blocks")
  }

  checkedValue(inputs, fallback) {
    return inputs.find((input) => input.checked)?.value || fallback
  }

  classify(name) {
    return name.charAt(0).toUpperCase() + name.slice(1)
  }

  titleize(value) {
    return value.charAt(0).toUpperCase() + value.slice(1)
  }
}
