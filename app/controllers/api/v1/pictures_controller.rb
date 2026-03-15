module Api
  module V1
    class PicturesController < BaseController
      before_action :set_book
      before_action -> { ensure_editable(@book) }, except: %i[ index show ]
      before_action :set_leaf, only: %i[ show update destroy ]

      def index
        pictures = @book.leaves.active.with_leafables.positioned.where(leafable_type: "Picture")
        render json: pictures.map { |l| picture_json(l) }
      end

      def show
        render json: picture_json(@leaf)
      end

      def create
        picture = Picture.new(picture_params)
        leaf = @book.press(picture, leaf_params)

        if position = params[:position]&.to_i
          leaf.move_to_position(position)
        end

        render json: picture_json(leaf), status: :created
      end

      def update
        @leaf.edit(leafable_params: picture_params, leaf_params: leaf_params)
        render json: picture_json(@leaf.reload)
      end

      def destroy
        @leaf.trashed!
        head :no_content
      end

      private
        def set_book
          @book = Book.accessable_or_published.find(params[:book_id])
        end

        def set_leaf
          @leaf = @book.leaves.active.find(params[:id])
        end

        def leaf_params
          params.fetch(:leaf, {}).permit(:title).reverse_merge(title: "Untitled")
        end

        def picture_params
          params.fetch(:picture, {}).permit(:image, :caption)
        end

        def picture_json(leaf)
          json = {
            id: leaf.id,
            title: leaf.title,
            type: "Picture",
            caption: leaf.leafable.caption,
            has_image: leaf.leafable.image.attached?,
            status: leaf.status,
            position: leaf.position_score,
            slug: leaf.slug,
            created_at: leaf.created_at,
            updated_at: leaf.updated_at
          }

          if leaf.leafable.image.attached?
            json[:image_url] = Rails.application.routes.url_helpers.rails_blob_path(leaf.leafable.image, only_path: true)
          end

          json
        end
    end
  end
end
