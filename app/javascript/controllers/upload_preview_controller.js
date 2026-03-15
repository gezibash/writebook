import { Controller } from "@hotwired/stimulus"

export default class extends Controller {
  static values = { defaultImage: String }
  static targets = [ "image", "input", "button", "svg" ]

  previewImage() {
    const file = this.inputTarget.files[0]

    if (file) {
      this.imageTarget.src = URL.createObjectURL(this.inputTarget.files[0]);
      this.imageTarget.onload = () => { URL.revokeObjectURL(this.imageTarget.src) }
      this.imageTarget.hidden = false
      if (this.hasButtonTarget) this.buttonTarget.hidden = false
      if (this.hasSvgTarget) this.svgTarget.hidden = true
    }
  }

  clear() {
    this.imageTarget.src = this.defaultImageValue
    this.imageTarget.hidden = true
    if (this.hasButtonTarget) this.buttonTarget.hidden = true
    this.inputTarget.value = ""
    if (this.hasSvgTarget) this.svgTarget.hidden = false
  }
}
