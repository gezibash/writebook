require "securerandom"

class AddCoverSeedToBooks < ActiveRecord::Migration[8.0]
  class BookRecord < ApplicationRecord
    self.table_name = "books"
  end

  def up
    add_column :books, :cover_seed, :string

    BookRecord.reset_column_information
    BookRecord.find_each do |book|
      book.update_columns(cover_seed: SecureRandom.hex(8))
    end

    change_column_null :books, :cover_seed, false
  end

  def down
    remove_column :books, :cover_seed
  end
end
