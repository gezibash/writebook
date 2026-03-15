require "securerandom"

class Book < ApplicationRecord
  include Accessable, Sluggable

  COVER_STYLES = %w[ glass rings shapes identicon ].freeze

  attribute :cover_seed, :string, default: -> { SecureRandom.hex(8) }

  has_many :leaves, dependent: :destroy
  has_one_attached :cover, dependent: :purge_later

  scope :ordered, -> { order(:title) }
  scope :published, -> { where(published: true) }

  enum :theme, %w[ black blue green magenta orange violet white ].index_by(&:itself), suffix: true, default: :blue
  enum :cover_style, COVER_STYLES.index_by(&:itself), suffix: true, default: :glass

  before_validation :ensure_cover_seed

  def press(leafable, leaf_params)
    leaves.create! leaf_params.merge(leafable: leafable)
  end

  private
    def ensure_cover_seed
      self.cover_seed = SecureRandom.hex(8) if cover_seed.blank?
    end
end
