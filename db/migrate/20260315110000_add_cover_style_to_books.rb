class AddCoverStyleToBooks < ActiveRecord::Migration[8.0]
  def change
    add_column :books, :cover_style, :string, null: false, default: "glass"
  end
end
